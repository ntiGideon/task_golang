package user

import (
	"awesomeProject2/helpers"
	"awesomeProject2/model/userModel"
	userServ "awesomeProject2/service/user"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserController struct {
	UserService *userServ.UserService
}

func NewUserController(userService *userServ.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (controller *UserController) Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := userModel.RegisterUserModel{}
	helpers.ReadRequestBody(r, &userDto)

	webResponse := controller.UserService.Register(r.Context(), userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) VerifyAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := userModel.VerifyAccountModel{}
	helpers.ReadRequestBody(r, &userDto)

	webResponse := controller.UserService.VerifyAccount(r.Context(), userDto.EmailToken)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := userModel.LoginUserModel{}
	helpers.ReadRequestBody(r, &userDto)
	webResponse := controller.UserService.Login(r.Context(), userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *UserController) BeginPasswordReset(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userDto := userModel.BeginPasswordReset{}
	helpers.ReadRequestBody(r, &userDto)

	webResponse := controller.UserService.BeginPasswordReset(r.Context(), &userDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}
