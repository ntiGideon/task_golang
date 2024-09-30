package user

import (
	"awesomeProject2/data/user"
	mailDto "awesomeProject2/data/user"
	"awesomeProject2/helpers"
	mailServ "awesomeProject2/mail/user"
	"awesomeProject2/model/userModel"
	"awesomeProject2/prisma/db"
	"awesomeProject2/repository/userRepository"
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"time"
)

type UserService struct {
	Db          *db.PrismaClient
	RedisClient *redis.Client
}

func NewUserService(db *db.PrismaClient, redisClient *redis.Client) *UserService {
	return &UserService{
		Db:          db,
		RedisClient: redisClient,
	}
}

func (p *UserService) Register(ctx context.Context, registerDto userModel.RegisterUserModel) *user.WebResponse {
	validateDto := helpers.RequestValidators(registerDto)
	if validateDto != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validateDto.Error(),
		}
	}

	existingUserByEmail, _ := userRepository.ExistingUserByEmail(ctx, p.Db, registerDto.Email, registerDto.Username)

	if existingUserByEmail {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
			Data:    nil,
		}
	}
	hashedPassword := helpers.HashPassword(registerDto.Password)
	var emailToken = uuid.New().String()
	createdUser, err := p.Db.User.CreateOne(
		db.User.Username.Set(registerDto.Username),
		db.User.Email.Set(registerDto.Email),
		db.User.Name.Set(registerDto.Name),
		db.User.Password.Set(hashedPassword),
		db.User.State.Set(db.StateEnumFresh),
	).Exec(ctx)
	if err != nil {
		helpers.PanicAllErrors(err)
	}

	// redis store
	redisContext := context.Background()
	_ = p.RedisClient.Set(redisContext, emailToken, createdUser.ID, time.Hour*3)

	mail := &mailDto.MailInputs{
		Email:    registerDto.Email,
		Code:     emailToken,
		Username: registerDto.Username,
	}

	err = mailServ.VerifyEmail(mail)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "could not send email",
			Data:    err.Error(),
		}
	}

	return &user.WebResponse{
		Code:    http.StatusCreated,
		Message: "User registered!",
		Data:    nil,
	}
}

func (p *UserService) VerifyAccount(ctx context.Context, emailToken string) *user.WebResponse {
	userId := p.RedisClient.Get(ctx, emailToken).Val()
	if userId == "" {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification token!",
		}
	}
	id, _ := strconv.ParseInt(userId, 10, 0)
	userExist, _ := p.Db.User.FindUnique(db.User.ID.Equals(int(id))).Exec(ctx)
	if userExist == nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist!",
			Data:    nil,
		}
	}

	p.Db.User.FindUnique(db.User.ID.Equals(int(id))).Update(
		db.User.State.Set(db.StateEnumVerified),
	).Exec(ctx)
	p.RedisClient.Del(ctx, emailToken)
	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "User verified!",
		Data:    nil,
	}
}

func (p *UserService) Login(ctx context.Context, loginDto userModel.LoginUserModel) *user.WebResponse {
	validateDto := helpers.RequestValidators(loginDto)
	if validateDto != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validateDto.Error(),
		}
	}

	existingUser, _ := p.Db.User.FindFirst(
		db.User.Or(
			db.User.Email.Equals(loginDto.EmailOrUsername),
			db.User.Username.Equals(loginDto.EmailOrUsername),
		)).Exec(ctx)
	if existingUser == nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist!",
			Data:    nil,
		}
	}

	if existingUser.State == db.StateEnumFresh {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Account not verified!",
			Data:    nil,
		}
	}

	if existingUser.State == db.StateEnumDeleted {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Account deleted!",
			Data:    nil,
		}
	}

	correctPassword := helpers.CheckPasswordHash(loginDto.Password, existingUser.Password)
	if !correctPassword {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid password!",
			Data:    nil,
		}
	}
	jwtPayload := &userModel.JWTPayload{
		Username: existingUser.Username,
		Email:    existingUser.Email,
		Id:       existingUser.ID,
	}
	token, err := helpers.GenerateJwt(jwtPayload)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Token error",
			Data:    nil,
		}
	}
	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "User logged in!",
		Data: struct {
			Id          int    `json:"id"`
			Email       string `json:"email"`
			Username    string `json:"username"`
			AccessToken string `json:"access_token"`
		}{
			Id:          existingUser.ID,
			Email:       existingUser.Email,
			Username:    existingUser.Username,
			AccessToken: token,
		},
	}
}

func (p *UserService) BeginPasswordReset(ctx context.Context, emailDto *userModel.BeginPasswordReset) *user.WebResponse {
	validateDto := helpers.RequestValidators(emailDto)
	if validateDto != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validateDto.Error(),
		}
	}

	userExist, _ := p.Db.User.FindUnique(db.User.Email.Equals(emailDto.Email)).Exec(ctx)
	if userExist != nil && userExist.State == db.StateEnumVerified {
		resetPasswordToken := uuid.New().String()
		p.RedisClient.Set(ctx, resetPasswordToken, userExist.ID, time.Hour*3)

		mailInput := &mailDto.MailInputs{
			Email:    userExist.Email,
			Code:     resetPasswordToken,
			Username: userExist.Username,
		}
		err := mailServ.ResetPassword(mailInput)
		if err != nil {
			return &user.WebResponse{
				Code:    http.StatusBadRequest,
				Message: "Could not send email!",
				Data:    err.Error(),
			}
		}

		return &user.WebResponse{
			Code:    http.StatusOK,
			Message: "Reset-password email sent!",
			Data:    nil,
		}
	}

	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "Reset-password email sent!",
		Data:    nil,
	}

}
