package task

import (
	"awesomeProject2/helpers"
	"awesomeProject2/model/taskModel"
	"awesomeProject2/service/task"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type TaskController struct {
	TaskService *task.TaskService
}

func NewTaskController(taskService *task.TaskService) *TaskController {
	return &TaskController{TaskService: taskService}
}

func (controller *TaskController) CreateTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	taskDto := taskModel.CreateTaskModel{}
	helpers.ReadRequestBody(r, &taskDto)
	userId := r.Context().Value("userId").(int)
	taskDto.UserId = userId

	webResponse := controller.TaskService.CreateTask(r.Context(), &taskDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) GetAllTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := r.Context().Value("userId").(int)
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "6"
	}
	pageNumber, _ := strconv.Atoi(page)
	limitNumber, _ := strconv.Atoi(limit)
	webResponse := controller.TaskService.GetAllTask(r.Context(), userId, pageNumber, limitNumber)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := r.Context().Value("userId").(int)
	taskId := params.ByName("taskId")
	taskIdConv, _ := strconv.Atoi(taskId)
	taskDto := taskModel.UpdateTaskModel{}
	helpers.ReadRequestBody(r, &taskDto)
	taskDto.UserId = userId
	taskDto.TaskId = taskIdConv
	webResponse := controller.TaskService.UpdateTask(r.Context(), &taskDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := r.Context().Value("userId").(int)
	taskId := params.ByName("taskId")
	taskIdConv, _ := strconv.Atoi(taskId)
	webResponse := controller.TaskService.DeleteTask(r.Context(), userId, taskIdConv)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) MarkTaskCompleted(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := r.Context().Value("userId").(int)
	taskId := params.ByName("taskId")
	taskIdConv, _ := strconv.Atoi(taskId)
	webResponse := controller.TaskService.MarkTaskCompleted(r.Context(), userId, taskIdConv)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) SetDueDate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	taskDto := taskModel.SetDueDateTaskModel{}
	helpers.ReadRequestBody(r, &taskDto)

	userId := r.Context().Value("userId").(int)
	taskId := params.ByName("taskId")
	taskIdConv, _ := strconv.Atoi(taskId)
	taskDto.UserId = userId
	taskDto.TaskId = taskIdConv
	webResponse := controller.TaskService.SetDueDate(r.Context(), &taskDto)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}

func (controller *TaskController) GetTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := r.Context().Value("userId").(int)
	taskId := params.ByName("taskId")
	taskIdConv, _ := strconv.Atoi(taskId)
	webResponse := controller.TaskService.GetTask(r.Context(), userId, taskIdConv)
	helpers.WriteResponseBody(w, webResponse, webResponse.Code)
}
