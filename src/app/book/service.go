package book

import (
	"context"

	"github.com/rs/zerolog/log"
)

// Compile-time check to ensure service implements Service interface
var _ Service = (*service)(nil)

// service implements the Service interface
type service struct {
	storage Storage
}

// NewService creates a new book service instance
func NewService(storage Storage) Service {
	return &service{
		storage: storage,
	}
}

// CreateBook creates a new book
func (s *service) CreateBook(ctx context.Context, input *CreateBookInput) *CreateBookOutput {
	resp := &CreateBookOutput{TraceId: input.TraceId}

	// Validate input
	if input.Title == "" {
		log.Warn().Msg("Title empty")
		resp.Message = "Title is mandatory"
		return resp
	}

	if input.Author == "" {
		log.Warn().Msg("Author empty")
		resp.Message = "Author is mandatory"
		return resp
	}

	// Begin transaction
	db, err := s.storage.BeginTx(ctx)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to begin transaction")
		return resp
	}
	defer db.Rollback()

	// Insert book
	insertedId, errType, err := db.InsertBook(ctx, input.Title, input.Author, input.ISBN, input.Publisher, input.PublishedYear, input.Pages, input.Description)
	if err != nil {
		if errType == ErrTypeUniqueConstraint {
			log.Err(err).Str("traceId", input.TraceId).Msg("failed to insert book")
			resp.Message = "Book with this ISBN already exists"
			return resp
		}
		log.Err(err).Str("traceId", input.TraceId).Msg("failed to insert book")
		return resp
	}

	// Get created book
	data, err := db.GetBookById(ctx, insertedId)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("failed to get book")
		return resp
	}

	// Commit transaction
	err = db.Commit()
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("failed to commit")
		return resp
	}

	log.Info().
		Str("traceId", input.TraceId).
		Str("title", input.Title).
		Str("author", input.Author).
		Msg("Book created successfully")

	resp.Success = true
	resp.Message = "Book created successfully"
	resp.Book = data
	return resp
}

// GetBookById retrieves a book by ID
func (s *service) GetBookById(ctx context.Context, input *GetBookByIdInput) *GetBookByIdOutput {
	resp := &GetBookByIdOutput{TraceId: input.TraceId}

	// Begin transaction
	db, err := s.storage.BeginTx(ctx)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to begin transaction")
		resp.Message = "Failed to fetch book"
		return resp
	}
	defer db.Rollback()

	// Get book by ID
	book, err := db.GetBookById(ctx, input.Id)
	if err != nil {
		log.Err(err).
			Str("traceId", input.TraceId).
			Int("id", input.Id).
			Msg("Book not found")
		resp.Message = "Book not found"
		return resp
	}

	log.Info().
		Str("traceId", input.TraceId).
		Int("id", input.Id).
		Msg("Book retrieved successfully")

	resp.Success = true
	resp.Message = "Book retrieved successfully"
	resp.Book = book

	return resp
}

// GetBooks retrieves books with pagination
func (s *service) GetBooks(ctx context.Context, input *GetBooksInput) *GetBooksOutput {
	resp := &GetBooksOutput{TraceId: input.TraceId}

	// Set default pagination values
	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	page := input.Page
	if page <= 0 {
		page = 1 // Default to first page
	}

	// Convert page to offset
	offset := (page - 1) * limit

	// Begin transaction
	db, err := s.storage.BeginTx(ctx)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to begin transaction")
		resp.Message = "Failed to fetch books"
		return resp
	}
	defer db.Rollback()

	// Get books with pagination and optional author filter
	books, err := db.GetBooks(ctx, limit, offset, input.Author)
	if err != nil {
		log.Err(err).
			Str("traceId", input.TraceId).
			Msg("Failed to get books")
		resp.Message = "Failed to fetch books"
		return resp
	}

	log.Info().
		Str("traceId", input.TraceId).
		Int("count", len(books)).
		Int("page", page).
		Int("limit", limit).
		Str("author", input.Author).
		Msg("Books retrieved successfully")

	resp.Success = true
	resp.Message = "Books retrieved successfully"
	resp.Books = books
	resp.Total = len(books)

	return resp
}

// UpdateBook updates a book's information
func (s *service) UpdateBook(ctx context.Context, input *UpdateBookInput) *UpdateBookOutput {
	resp := &UpdateBookOutput{TraceId: input.TraceId}

	// Validate input
	if input.Title == "" {
		log.Warn().Msg("Title empty")
		resp.Message = "Title is mandatory"
		return resp
	}

	if input.Author == "" {
		log.Warn().Msg("Author empty")
		resp.Message = "Author is mandatory"
		return resp
	}

	// Begin transaction
	db, err := s.storage.BeginTx(ctx)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to begin transaction")
		return resp
	}
	defer db.Rollback()

	// Lock book row for update (pessimistic lock)
	_, errType, err := db.LockBookById(ctx, input.Id)
	if err != nil {
		if errType == ErrTypeNotFound {
			log.Err(err).Str("traceId", input.TraceId).Msg("Book not found")
			resp.Message = "Book not found"
			return resp
		}
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to lock book")
		resp.Message = "Failed to update book"
		return resp
	}

	// Update book
	err = db.UpdateBook(ctx, input.Id, input.Title, input.Author, input.ISBN, input.Publisher, input.PublishedYear, input.Pages, input.Description)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to update book")
		resp.Message = "Failed to update book"
		return resp
	}

	// Get updated book
	data, err := db.GetBookById(ctx, input.Id)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to get book")
		resp.Message = "Book not found"
		return resp
	}

	// Commit transaction
	err = db.Commit()
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to commit")
		resp.Message = "Failed to update book"
		return resp
	}

	log.Info().
		Str("traceId", input.TraceId).
		Int("id", input.Id).
		Str("title", input.Title).
		Msg("Book updated successfully")

	resp.Success = true
	resp.Message = "Book updated successfully"
	resp.Book = data

	return resp
}

// DeleteBook deletes a book by ID
func (s *service) DeleteBook(ctx context.Context, input *DeleteBookInput) *DeleteBookOutput {
	resp := &DeleteBookOutput{TraceId: input.TraceId}

	// Begin transaction
	db, err := s.storage.BeginTx(ctx)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to begin transaction")
		return resp
	}
	defer db.Rollback()

	// Lock book row for update (pessimistic lock)
	_, errType, err := db.LockBookById(ctx, input.Id)
	if err != nil {
		if errType == ErrTypeNotFound {
			log.Err(err).Str("traceId", input.TraceId).Msg("Book not found")
			resp.Message = "Book not found"
			return resp
		}
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to lock book")
		resp.Message = "Failed to delete book"
		return resp
	}

	// Delete book
	err = db.DeleteBook(ctx, input.Id)
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to delete book")
		resp.Message = "Failed to delete book"
		return resp
	}

	// Commit transaction
	err = db.Commit()
	if err != nil {
		log.Err(err).Str("traceId", input.TraceId).Msg("Failed to commit")
		resp.Message = "Failed to delete book"
		return resp
	}

	log.Info().
		Str("traceId", input.TraceId).
		Int("id", input.Id).
		Msg("Book deleted successfully")

	resp.Success = true
	resp.Message = "Book deleted successfully"

	return resp
}
