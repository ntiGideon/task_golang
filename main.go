package main

import (
	"awesomeProject2/config"
	taskCont "awesomeProject2/controller/task"
	userCont "awesomeProject2/controller/user"
	"awesomeProject2/helpers"
	"awesomeProject2/router"
	taskServ "awesomeProject2/service/task"
	userServ "awesomeProject2/service/user"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic("Sorry, we can't proceed!")
	}

	fmt.Printf("Server starting on PORT %s \n", os.Getenv("PORT"))

	db, err := config.ConnectDB()
	if err != nil {
		helpers.PanicAllErrors(err)
	}

	redisClient, err := config.ConnectRedis()

	defer db.Prisma.Disconnect()

	userService := userServ.NewUserService(db, redisClient)
	userController := userCont.NewUserController(userService)
	taskService := taskServ.NewTaskService(db)
	taskController := taskCont.NewTaskController(taskService)

	routes := router.NewRouter(userController, taskController)

	server := http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverError := server.ListenAndServe()
	if serverError != nil {
		helpers.PanicAllErrors(serverError)
	}
}
