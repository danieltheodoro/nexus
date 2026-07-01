package main

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Project struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type List struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
}

type Task struct {
	ID          int    `json:"id"`
	ListID      int    `json:"list_id"`
	CreatorID   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Position    int    `json:"position"`
}

type Label struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
}

type Comment struct {
	ID       int    `json:"id"`
	TaskID   int    `json:"task_id"`
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}
