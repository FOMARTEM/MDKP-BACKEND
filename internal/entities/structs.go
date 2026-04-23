package entities

type User struct {
	ID          int    `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	IsActive    bool   `json:"is_active"`
	Role        string `json:"role"`
	RoleID      int    `json:"role_id"`
}

type Tasks struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateCreated string `json:"date_created"`
	DateDedline string `json:"date_deadline"`
	DateClosed  string `json:"date_closed"`
	Priority    int    `json:"priority"`
	IdCreator   int    `json:"id_creator"`
	IdRedactor  int    `json:"id_redactor"`
	IdAuthor    int    `json:"id_author"`
	IdStatus    int    `json:"id_status"`
}

type Material struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Extension   string `json:"extension"`
	Description string `json:"description"`
	CreatorID   int    `json:"creator_id"`
	TaskID      int    `json:"task_id"`
}

type Version struct {
	ID            int    `json:"id"`
	NumberVersion int    `json:"number_version"`
	DateCreated   string `json:"date_created"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CreatorID     int    `json:"creator_id"`
	MaterialID    int    `json:"material_id"`
	TaskID        int    `json:"task_id"`
}

type Status struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Role struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Revision struct {
	ID             int    `json:"id"`
	NumberRevision int    `json:"number_revision"`
	DateCreated    string `json:"date_created"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CreatorID      int    `json:"creator_id"`
	VersionID      int    `json:"version_id"`
	StatusID       int    `json:"status_id"`
}

type Log struct {
	ID          int    `json:"id"`
	DateCreated string `json:"date_created"`
	Action      string `json:"action"`
	UserID      int    `json:"user_id"`
}
