package models

type Message struct {
	ID        string `json:"id"`
	To        []uint `json:"to"`
	ProjectID *uint  `json:"project_id,omitempty"`
	TaskID    *uint  `json:"task_id,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
