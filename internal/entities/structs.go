package entities

type User struct {
	ID           int    `json:"id,omitempty"`
	LastName     string `json:"last_name"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	Phone        string `json:"phone"`
	DateOfBirth  string `json:"date_of_birth"`
	Email        string `json:"email" validate:"email"`
	Password     string `json:"password,omitempty" validate:"min=8,max=16"`
	PasswordHash string `json:"password_hash,omitempty"`
	Token        string `json:"token,omitempty"`
	IsActive     bool   `json:"is_active,omitempty"`
	Role         string `json:"role,omitempty"`
	RoleID       int    `json:"role_id,omitempty"`
}

type Tasks struct {
	ID          int    `json:"id"`
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"true"`
	DateCreated string `json:"date_created"`
	DateDedline string `json:"date_deadline"`
	DateClosed  string `json:"date_closed" `
	Priority    int    `json:"priority" required:"true"`
	IdCreator   int    `json:"id_creator"`
	IdRedactor  int    `json:"id_redactor"`
	IdAuthor    int    `json:"id_author"`
	IdStatus    int    `json:"id_status"`
}

type Material struct {
	ID          int    `json:"id"`
	Title       string `json:"title" form:"title" required:"true"`
	Extension   string `json:"extension"`
	Description string `json:"description" form:"description" required:"true"`
	CreatorID   int    `json:"creator_id" form:"creator_id"`
	TaskID      int    `json:"task_id" form:"task_id"`
}

type Version struct {
	ID            int    `json:"id"`
	NumberVersion int    `json:"number_version"`
	DateCreated   string `json:"date_created" validate:"date"`
	Title         string `json:"title" required:"true"`
	Description   string `json:"description" required:"true"`
	CreatorID     int    `json:"creator_id"`
	MaterialID    int    `json:"material_id"`
	TaskID        int    `json:"task_id"`
}

type Status struct {
	ID          int    `json:"id"`
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"true"`
}

type Role struct {
	ID          int    `json:"id"`
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"true"`
}

type Revision struct {
	ID             int    `json:"id"`
	NumberRevision int    `json:"number_revision"`
	DateCreated    string `json:"date_created" validate:"date"`
	Title          string `json:"title" required:"true"`
	Description    string `json:"description" required:"true"`
	CreatorID      int    `json:"creator_id"`
	VersionID      int    `json:"version_id"`
	StatusID       int    `json:"status_id"`
}

type Log struct {
	ID          int    `json:"id"`
	DateCreated string `json:"date_created" validate:"date"`
	Action      string `json:"action"`
	UserID      int    `json:"user_id"`
}
