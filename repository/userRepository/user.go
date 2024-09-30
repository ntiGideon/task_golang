package userRepository

import (
	"awesomeProject2/prisma/db"
	"context"
)

func ExistingUserByEmail(ctx context.Context, dbClient *db.PrismaClient, email string, username string) (bool, error) {
	user, err := dbClient.User.FindFirst(
		db.User.Or(
			db.User.Email.Equals(email),
			db.User.Username.Equals(username),
		),
	).Exec(ctx)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}
