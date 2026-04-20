package models

type Exec struct {
	ID                   int     `json:"id" db:"id"`
	FirstName            string  `json:"first_name" db:"first_name"`
	LastName             string  `json:"last_name" db:"last_name"`
	Email                string  `json:"email" db:"email"`
	Username             string  `json:"username" db:"username"`
	Password             string  `json:"password,omitempty" db:"password"`
	PasswordUpdatedAt    *string `json:"password_updated_at,omitempty" db:"password_updated_at"`
	UserCreatedAt        *string `json:"user_created_at,omitempty" db:"user_created_at"`
	PasswordResetCode    *string `json:"password_reset_code,omitempty" db:"password_reset_code"`
	PasswordCodeExpires  *string `json:"password_code_expires_at,omitempty" db:"password_code_expires_at"`
	PasswordTokenExpires *string `json:"password_token_expires,omitempty" db:"password_token_expires"`
	IsActive             bool    `json:"is_active" db:"is_active"`
	Role                 string  `json:"role" db:"role"`
}
