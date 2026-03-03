package book

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time check to ensure storage implements Storage interface
var _ Storage = (*storage)(nil)

// storage implements the Storage interface
type storage struct {
	pool *pgxpool.Pool
}

// NewStorage creates a new storage instance
func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{
		pool: pool,
	}
}

// BeginTx starts a new database transaction
func (s *storage) BeginTx(ctx context.Context) (StorageTx, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &storageTx{
		tx: tx,
	}, nil
}

// Compile-time check to ensure storageTx implements StorageTx interface
var _ StorageTx = (*storageTx)(nil)

// storageTx implements the StorageTx interface
type storageTx struct {
	tx pgx.Tx
}

func (st *storageTx) InsertBook(ctx context.Context, title, author, isbn, publisher string, publishedYear, pages int, description string) (int, StorageErrorType, error) {
	query := `INSERT INTO books (title, author, isbn, publisher, published_year, pages, description) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int
	err := st.tx.QueryRow(ctx, query, title, author, isbn, publisher, publishedYear, pages, description).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// 23505 is the PostgreSQL error code for unique_violation
			if pgErr.Code == "23505" || strings.Contains(pgErr.Message, "duplicate key") {
				return 0, ErrTypeUniqueConstraint, fmt.Errorf("failed to insert book: %w", err)
			}
		}
		return 0, ErrTypeCommon, fmt.Errorf("failed to insert book: %w", err)
	}
	return id, ErrTypeNone, nil
}

func (st *storageTx) GetBookById(ctx context.Context, id int) (Book, error) {
	query := `SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at FROM books WHERE id = $1`
	row := st.tx.QueryRow(ctx, query, id)
	book, err := convertBookRow(row)
	if err != nil {
		return Book{}, fmt.Errorf("failed to get book by id: %w", err)
	}
	return book, nil
}

func (st *storageTx) GetBooks(ctx context.Context, limit, offset int, author string) ([]Book, error) {
	var query string
	var args []interface{}
	
	if author != "" {
		// Search by author (case-insensitive)
		query = `SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at 
		         FROM books 
		         WHERE LOWER(author) LIKE LOWER($1)
		         ORDER BY created_at DESC 
		         LIMIT $2 OFFSET $3`
		args = []interface{}{"%" + author + "%", limit, offset}
	} else {
		// Get all books
		query = `SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at 
		         FROM books 
		         ORDER BY created_at DESC 
		         LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}
	
	rows, err := st.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.Publisher, &book.PublishedYear, &book.Pages, &book.Description, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book row: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating book rows: %w", err)
	}

	return books, nil
}

func (st *storageTx) UpdateBook(ctx context.Context, id int, title, author, isbn, publisher string, publishedYear, pages int, description string) error {
	query := `UPDATE books SET title = $1, author = $2, isbn = $3, publisher = $4, published_year = $5, pages = $6, description = $7, updated_at = NOW() WHERE id = $8`
	_, err := st.tx.Exec(ctx, query, title, author, isbn, publisher, publishedYear, pages, description, id)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}
	return nil
}

func (st *storageTx) DeleteBook(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := st.tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// LockBookById locks a book row for update and returns the book
// This implements pessimistic locking to prevent concurrent modifications
func (st *storageTx) LockBookById(ctx context.Context, id int) (Book, StorageErrorType, error) {
	query := `SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at FROM books WHERE id = $1 FOR UPDATE`
	row := st.tx.QueryRow(ctx, query, id)
	book, err := convertBookRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Book{}, ErrTypeNotFound, fmt.Errorf("failed to lock book by id: %w", err)
		}
		return Book{}, ErrTypeCommon, fmt.Errorf("failed to lock book by id: %w", err)
	}
	return book, ErrTypeNone, nil
}

func convertBookRow(row pgx.Row) (Book, error) {
	var book Book
	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.Publisher, &book.PublishedYear, &book.Pages, &book.Description, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

// Commit commits the transaction
func (st *storageTx) Commit() error {
	err := st.tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the transaction
func (st *storageTx) Rollback() error {
	err := st.tx.Rollback(context.Background())
	if err != nil && err != pgx.ErrTxClosed {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
