package main

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AvatarURL    string `json:"avatar_url"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	AvatarURL   string `json:"avatar_url"`
	IsActive    bool   `json:"is_active"`
}

type DeleteUserRequest struct {
	Password string `json:"password"`
}

type Project struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IsArchived  bool   `json:"is_archived"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateProjectRequest struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IsArchived  bool   `json:"is_archived"`
}

type List struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Task struct {
	ID          int    `json:"id"`
	ListID      int    `json:"list_id"`
	CreatorID   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Position    int    `json:"position"`
	DueDate     string `json:"due_date"`
	CompletedAt string `json:"completed_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Label struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Comment struct {
	ID        int    `json:"id"`
	TaskID    int    `json:"task_id"`
	AuthorID  int    `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
