package taskModel

import (
	"awesomeProject2/prisma/db"
	"time"
)

type CreateTaskModel struct {
	Title       string          `json:"title" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Priority    db.TaskPriority `json:"priority" validate:"required"`
	Category    db.TaskCategory `json:"category" validate:"required"`
	Status      db.TaskStatus   `json:"status" validate:"required"`
	DueDate     time.Time       `json:"dueDate" validate:"required"`
	UserId      int             `json:"user_id"`
}

type UpdateTaskModel struct {
	TaskId      int             `json:"task_id"`
	Title       string          `json:"title" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Priority    db.TaskPriority `json:"priority" validate:"required"`
	Category    db.TaskCategory `json:"category" validate:"required"`
	Status      db.TaskStatus   `json:"status" validate:"required"`
	DueDate     time.Time       `json:"dueDate" validate:"required"`
	UserId      int             `json:"user_id"`
}

type SetDueDateTaskModel struct {
	TaskId  int
	UserId  int
	DueDate time.Time `json:"dueDate" validate:"required"`
}
