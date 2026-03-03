package user

import (
	"context"
)

// Service defines the interface for user business logic
type Service interface {
	// Register creates a new user account
	Register(ctx context.Context, input *RegisterInput) *RegisterOutput

	// Login authenticates a user
	Login(ctx context.Context, input *LoginInput) *LoginOutput

	// UpdateUsername updates a user's username
	UpdateUsername(ctx context.Context, input *UpdateUsernameInput) *UpdateUsernameOutput

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, input *UpdatePasswordInput) *UpdatePasswordOutput

	// GetProfileById retrieves a user profile by ID
	GetProfileById(ctx context.Context, input *GetProfileByIdInput) *GetProfileByIdOutput
}

// RegisterInput represents input for user registration
type RegisterInput struct {
	TraceId  string
	Username string
	Email    string
	FullName string
	Password string
}

// RegisterOutput represents output after user registration
type RegisterOutput struct {
	Success bool
	Message string
	TraceId string
	User    User
}

// LoginInput represents input for user login
type LoginInput struct {
	TraceId  string
	Username string
	Password string
}

// LoginOutput represents output after user login
type LoginOutput struct {
	Success bool
	Message string
	TraceId string
	User    User
}

// UpdateUsernameInput represents input for updating username
type UpdateUsernameInput struct {
	TraceId     string
	Id          int
	NewUsername string
}

// UpdateUsernameOutput represents output after updating username
type UpdateUsernameOutput struct {
	Success bool
	Message string
	TraceId string
	User    User
}

// UpdatePasswordInput represents input for updating password
type UpdatePasswordInput struct {
	TraceId     string
	Id          int
	OldPassword string
	NewPassword string
}

// UpdatePasswordOutput represents output after updating password
type UpdatePasswordOutput struct {
	Success bool
	Message string
	TraceId string
}

// GetProfileByIdInput represents input for getting user profile
type GetProfileByIdInput struct {
	TraceId string
	Id      int
}

// GetProfileByIdOutput represents output after getting user profile
type GetProfileByIdOutput struct {
	Success bool
	Message string
	TraceId string
	User    User
}
