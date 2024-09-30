package router

import (
	"awesomeProject2/controller/task"
	"awesomeProject2/controller/user"
	"awesomeProject2/middleware"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(userController *user.UserController, taskController *task.TaskController) *httprouter.Router {
	router := httprouter.New()

	// user routes
	router.POST("/api/user/register", userController.Register)
	router.PUT("/api/user/verify", userController.VerifyAccount)
	router.POST("/api/user/begin-password-reset", userController.BeginPasswordReset)
	router.POST("/api/user/login", userController.Login)

	// task routes
	router.POST("/api/task/create", middleware.AuthMiddleware(taskController.CreateTask))
	router.GET("/api/task", middleware.AuthMiddleware(taskController.GetAllTask))
	router.PUT("/api/task/update/:taskId", middleware.AuthMiddleware(taskController.UpdateTask))
	router.DELETE("/api/task/delete/:taskId", middleware.AuthMiddleware(taskController.DeleteTask))
	router.PUT("/api/task/mark-complete/:taskId", middleware.AuthMiddleware(taskController.MarkTaskCompleted))
	router.PUT("/api/task/set-due-date/:taskId", middleware.AuthMiddleware(taskController.SetDueDate))
	router.GET("/api/task/:taskId", middleware.AuthMiddleware(taskController.GetTask))

	return router
}
