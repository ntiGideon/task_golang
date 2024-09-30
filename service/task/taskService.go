package task

import (
	"awesomeProject2/data/user"
	"awesomeProject2/helpers"
	"awesomeProject2/model/taskModel"
	"awesomeProject2/prisma/db"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type TaskService struct {
	Db *db.PrismaClient
}

func NewTaskService(db *db.PrismaClient) *TaskService {
	return &TaskService{Db: db}
}

func (p *TaskService) CreateTask(ctx context.Context, taskModel *taskModel.CreateTaskModel) *user.WebResponse {
	validateDto := helpers.RequestValidators(taskModel)
	if validateDto != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validateDto.Error(),
		}
	}

	existingTaskByTitle, _ := p.Db.Task.FindFirst(db.Task.Title.Equals(taskModel.Title)).Exec(ctx)
	if existingTaskByTitle != nil {
		return &user.WebResponse{
			Code:    http.StatusConflict,
			Message: "Task 'Title' already exists",
			Data:    nil,
		}
	}
	var dueDate time.Time
	if !taskModel.DueDate.IsZero() {
		dueDate = taskModel.DueDate
	}

	_, err := p.Db.Task.CreateOne(
		db.Task.Title.Set(taskModel.Title),
		db.Task.Priority.Set(taskModel.Priority),
		db.Task.Category.Set(taskModel.Category),
		db.Task.Status.Set(taskModel.Status),
		db.Task.DueDate.Set(dueDate),
		db.Task.User.Link(db.User.ID.Set(taskModel.UserId)),
		db.Task.Description.Set(taskModel.Description),
	).Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &user.WebResponse{
		Code:    http.StatusCreated,
		Message: "Task created",
		Data:    nil,
	}
}

func (p *TaskService) GetAllTask(ctx context.Context, userId int, page int, limit int) *user.WebResponsePagination {
	totalCount, err := p.Db.Task.FindMany(db.Task.UserID.Equals(userId)).Exec(ctx)

	if err != nil {
		return &user.WebResponsePagination{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		}
	}

	tasks, err := p.Db.Task.FindMany(
		db.Task.UserID.Equals(userId),
	).Omit(
		db.Task.UserID.Field(),
		db.Task.UpdatedAt.Field(),
	).Take(limit).Skip((page - 1) * limit).OrderBy(db.Task.CreatedAt.Order(db.SortOrderDesc)).Exec(ctx)

	var taskData []user.TaskPaginationData
	for _, task := range tasks {
		description, _ := task.Description()
		dueDate, _ := task.DueDate()
		data := user.TaskPaginationData{
			Id:          task.ID,
			Title:       task.Title,
			Description: description,
			Priority:    task.Priority,
			Category:    task.Category,
			DueDate:     dueDate,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
		}

		taskData = append(taskData, data)
	}

	if err != nil {
		return &user.WebResponsePagination{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		}
	}

	totalPages := (len(totalCount) + limit - 1) / limit

	route := fmt.Sprintf("%v/api/post", os.Getenv("FRONTEND_URL"))
	first := fmt.Sprintf("%v?limit=%v&page=1", route, limit)
	last := fmt.Sprintf("%v?limit=%v&page=%v", route, limit, totalPages)

	var prev, next string
	if page > 1 {
		prev = fmt.Sprintf("%v?limit=%v&page=%v", route, limit, page-1)
	}

	if page < totalPages {
		next = fmt.Sprintf("%v?limit=%v&page=%v", route, limit, page+1)
	}

	return &user.WebResponsePagination{
		Code:   http.StatusOK,
		Status: "OK!",
		Data:   taskData,
		Meta: &user.Meta{
			CurrentPage:  page,
			ItemsPerPage: limit,
			ItemCount:    len(taskData),
			TotalCount:   len(totalCount),
			TotalPages:   totalPages,
		},
		Links: &user.Links{
			First:    first,
			Previous: prev,
			Next:     next,
			Last:     last,
		},
	}
}

func (p *TaskService) UpdateTask(ctx context.Context, updateTaskM *taskModel.UpdateTaskModel) *user.WebResponse {
	validateDto := helpers.RequestValidators(updateTaskM)
	if validateDto != nil {
		return &user.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validateDto.Error(),
		}
	}

	existingTask, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.ID.Equals(updateTaskM.TaskId),
			db.Task.UserID.Equals(updateTaskM.UserId),
		),
	).Exec(ctx)
	if existingTask == nil {
		return &user.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Task not found",
			Data:    nil,
		}
	}

	existingTaskByTitle, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.Title.Equals(updateTaskM.Title),
			db.Task.UserID.Equals(updateTaskM.UserId),
		),
	).Exec(ctx)
	if existingTaskByTitle != nil {
		return &user.WebResponse{
			Code:    http.StatusConflict,
			Message: "Task 'Title' already exists",
			Data:    nil,
		}
	}

	_, err := p.Db.Task.FindUnique(
		db.Task.ID.Equals(updateTaskM.TaskId),
	).Update(
		db.Task.Title.Set(updateTaskM.Title),
		db.Task.Description.Set(updateTaskM.Description),
		db.Task.Priority.Set(updateTaskM.Priority),
		db.Task.Status.Set(updateTaskM.Status),
		db.Task.Status.Set(updateTaskM.Status),
		db.Task.Category.Set(updateTaskM.Category),
		db.Task.User.Link(db.User.ID.Set(updateTaskM.UserId)),
	).Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "Task updated",
		Data:    nil,
	}

}

func (p *TaskService) DeleteTask(ctx context.Context, userId int, taskId int) *user.WebResponse {
	existingTask, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.UserID.Equals(userId),
			db.Task.ID.Equals(taskId)),
	).Exec(ctx)
	if existingTask == nil {
		return &user.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Task not found",
			Data:    nil,
		}
	}
	_, err := p.Db.Task.FindUnique(
		db.Task.ID.Equals(taskId),
	).Delete().Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "Task deleted",
		Data:    nil,
	}
}

func (p *TaskService) MarkTaskCompleted(ctx context.Context, userId int, taskId int) *user.WebResponse {
	existingTask, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.UserID.Equals(userId),
			db.Task.ID.Equals(taskId)),
	).Exec(ctx)
	if existingTask == nil {
		return &user.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Task not found",
			Data:    nil,
		}
	}

	_, err := p.Db.Task.FindUnique(
		db.Task.ID.Equals(taskId),
	).Update(db.Task.Status.Set(db.TaskStatusCompleted)).Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "Task marked as completed!",
		Data:    nil,
	}
}

func (p *TaskService) SetDueDate(ctx context.Context, taskDto *taskModel.SetDueDateTaskModel) *user.WebResponse {
	existingTask, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.UserID.Equals(taskDto.UserId),
			db.Task.ID.Equals(taskDto.TaskId)),
	).Exec(ctx)
	if existingTask == nil {
		return &user.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Task not found",
			Data:    nil,
		}
	}

	_, err := p.Db.Task.FindUnique(
		db.Task.ID.Equals(taskDto.TaskId),
	).Update(db.Task.DueDate.Set(taskDto.DueDate)).Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "Task due date set!",
		Data:    nil,
	}
}

func (p *TaskService) GetTask(ctx context.Context, userId int, taskId int) *user.WebResponse {
	existingTask, _ := p.Db.Task.FindFirst(
		db.Task.And(
			db.Task.UserID.Equals(userId),
			db.Task.ID.Equals(taskId)),
	).Exec(ctx)
	if existingTask == nil {
		return &user.WebResponse{
			Code:    http.StatusNotFound,
			Message: "Task not found",
			Data:    nil,
		}
	}
	data, err := p.Db.Task.FindUnique(db.Task.ID.Equals(taskId)).Exec(ctx)
	if err != nil {
		return &user.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	description, _ := data.Description()
	dueDate, _ := data.DueDate()

	taskData := &user.TaskPaginationData{
		Id:          data.ID,
		Title:       data.Title,
		Description: description,
		Priority:    data.Priority,
		Category:    data.Category,
		DueDate:     dueDate,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
	}

	return &user.WebResponse{
		Code:    http.StatusOK,
		Message: "OK!",
		Data:    taskData,
	}
}
