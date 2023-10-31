package forms

type SignupUserCommand struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Agency   *int   `json:"agency"`
	Daily    *int   `json:"daily"`
}

type UpdateUserCommand struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Agency *int   `json:"agency"`
	Daily  *int   `json:"daily"`
}

type LoginUserCommand struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PasswordResetCommand struct {
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
}

type ChangePasswordCommand struct {
	CurrentPassword string `json:"password_current" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Confirm         string `json:"password_confirmation" binding:"required"`
}

// ResendCommand defines resend email payload
type ResendCommand struct {
	Email string `json:"email" binding:"required"`
}
