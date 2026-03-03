package book

import (
	"context"
)

type StorageErrorType string

const (
	ErrTypeNone             StorageErrorType = ""
	ErrTypeUniqueConstraint StorageErrorType = "unique_constraint"
	ErrTypeNotFound         StorageErrorType = "not_found"
	ErrTypeCommon           StorageErrorType = "common"
)

// Storage defines the interface for book data access layer
type Storage interface {
	// BeginTx starts a new database transaction
	BeginTx(ctx context.Context) (StorageTx, error)
}

// StorageTx defines the interface for transactional book operations
type StorageTx interface {
	// InsertBook inserts a new book and returns the book ID
	InsertBook(ctx context.Context, title, author, isbn, publisher string, publishedYear, pages int, description string) (int, StorageErrorType, error)

	// GetBookById retrieves a book by ID
	GetBookById(ctx context.Context, id int) (Book, error)

	// GetBooks retrieves books with pagination and optional author filter
	GetBooks(ctx context.Context, limit, offset int, author string) ([]Book, error)

	// UpdateBook updates a book's title
	UpdateBook(ctx context.Context, id int, title string) error

	// DeleteBook deletes a book by ID
	DeleteBook(ctx context.Context, id int) error

	// Pessimistic Locking Methods (FOR UPDATE)
	// These methods acquire row-level locks and return the locked entity

	// LockBookById locks a book row for update and returns the book
	LockBookById(ctx context.Context, id int) (Book, StorageErrorType, error)

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}
