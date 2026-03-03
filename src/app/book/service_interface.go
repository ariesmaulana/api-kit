package book

import (
	"context"
)

// Service defines the interface for book business logic
type Service interface {
	// CreateBook creates a new book
	CreateBook(ctx context.Context, input *CreateBookInput) *CreateBookOutput

	// GetBookById retrieves a book by ID
	GetBookById(ctx context.Context, input *GetBookByIdInput) *GetBookByIdOutput

	// GetBooks retrieves books with pagination
	GetBooks(ctx context.Context, input *GetBooksInput) *GetBooksOutput

	// UpdateBook updates a book's information
	UpdateBook(ctx context.Context, input *UpdateBookInput) *UpdateBookOutput

	// DeleteBook deletes a book by ID
	DeleteBook(ctx context.Context, input *DeleteBookInput) *DeleteBookOutput
}

// CreateBookInput represents input for creating a book
type CreateBookInput struct {
	TraceId       string
	Title         string
	Author        string
	ISBN          string
	Publisher     string
	PublishedYear int
	Pages         int
	Description   string
}

// CreateBookOutput represents output after creating a book
type CreateBookOutput struct {
	Success bool
	Message string
	TraceId string
	Book    Book
}

// GetBookByIdInput represents input for getting a book by ID
type GetBookByIdInput struct {
	TraceId string
	Id      int
}

// GetBookByIdOutput represents output after getting a book by ID
type GetBookByIdOutput struct {
	Success bool
	Message string
	TraceId string
	Book    Book
}

// GetBooksInput represents input for getting books with pagination
type GetBooksInput struct {
	TraceId string
	Page    int    // Page number (1-based)
	Limit   int    // Items per page
	Author  string // Filter by author (optional)
}

// GetBooksOutput represents output after getting books
type GetBooksOutput struct {
	Success bool
	Message string
	TraceId string
	Books   []Book
	Total   int
}

// UpdateBookInput represents input for updating a book
type UpdateBookInput struct {
	TraceId       string
	Id            int
	Title         string
	Author        string
	ISBN          string
	Publisher     string
	PublishedYear int
	Pages         int
	Description   string
}

// UpdateBookOutput represents output after updating a book
type UpdateBookOutput struct {
	Success bool
	Message string
	TraceId string
	Book    Book
}

// DeleteBookInput represents input for deleting a book
type DeleteBookInput struct {
	TraceId string
	Id      int
}

// DeleteBookOutput represents output after deleting a book
type DeleteBookOutput struct {
	Success bool
	Message string
	TraceId string
}
