package user

import (
	"awesomeProject2/prisma/db"
	"time"
)

type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MailInputs struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Link     string `json:"link"`
}

type Meta struct {
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
	ItemCount    int `json:"item_count"`
	TotalCount   int `json:"total_count"`
	TotalPages   int `json:"total_pages"`
}

type Links struct {
	First    string `json:"first"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Last     string `json:"last"`
}

type WebResponsePagination struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Links  interface{} `json:"links,omitempty"`
}

type TaskPaginationData struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    db.TaskPriority `json:"priority"`
	Category    db.TaskCategory `json:"category"`
	Status      db.TaskStatus   `json:"status"`
	DueDate     time.Time       `json:"dueDate"`
	CreatedAt   time.Time       `json:"created_at"`
}
