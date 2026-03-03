package book_test

import (
	"context"
	"testing"

	"github.com/ariesmaulana/api-kit/src/app/book"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// TestHelper provides utility methods for test fixtures
type TestHelper struct {
	pool *pgxpool.Pool
}

// NewTestHelper creates a new helper instance
func NewTestHelper(pool *pgxpool.Pool) *TestHelper {
	return &TestHelper{
		pool: pool,
	}
}

// DataBook represents a book fixture for testing
type DataBook struct {
	Idx           int // Index in the fixture array
	Id            int // Actual database ID (populated after insert)
	Title         string
	Author        string
	ISBN          string
	Publisher     string
	PublishedYear int
	Pages         int
	Description   string
}

// InsertBook inserts a single book and returns it
func (h *TestHelper) InsertBook(ctx context.Context, t *testing.T, title, author, isbn, publisher string, publishedYear, pages int, description string) *book.Book {
	query := `
		INSERT INTO books (title, author, isbn, publisher, published_year, pages, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at
	`

	var b book.Book
	err := h.pool.QueryRow(ctx, query, title, author, isbn, publisher, publishedYear, pages, description).Scan(
		&b.Id,
		&b.Title,
		&b.Author,
		&b.ISBN,
		&b.Publisher,
		&b.PublishedYear,
		&b.Pages,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	assert.Nil(t, err)

	return &b
}

// InsertBooks inserts multiple books and returns them
func (h *TestHelper) InsertBooks(ctx context.Context, t *testing.T, books []BookFixture) []*book.Book {
	result := make([]*book.Book, 0, len(books))

	for _, fixture := range books {
		b := h.InsertBook(ctx, t, fixture.Title, fixture.Author, fixture.ISBN, fixture.Publisher, fixture.PublishedYear, fixture.Pages, fixture.Description)
		result = append(result, b)
	}

	return result
}

// BookFixture represents a book fixture for testing
type BookFixture struct {
	Title         string
	Author        string
	ISBN          string
	Publisher     string
	PublishedYear int
	Pages         int
	Description   string
}

// ClearBooks removes all books from the database
func (h *TestHelper) ClearBooks(ctx context.Context, t *testing.T) {
	_, err := h.pool.Exec(ctx, "DELETE FROM books")
	assert.Nil(t, err)
}

// GetBookById retrieves a book by ID
func (h *TestHelper) GetBookById(ctx context.Context, t *testing.T, id int) *book.Book {
	query := `
		SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	var b book.Book
	err := h.pool.QueryRow(ctx, query, id).Scan(
		&b.Id,
		&b.Title,
		&b.Author,
		&b.ISBN,
		&b.Publisher,
		&b.PublishedYear,
		&b.Pages,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	assert.Nil(t, err)

	return &b
}

// GetBookByISBN retrieves a book by ISBN
func (h *TestHelper) GetBookByISBN(ctx context.Context, t *testing.T, isbn string) *book.Book {
	query := `
		SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at
		FROM books
		WHERE isbn = $1
	`

	var b book.Book
	err := h.pool.QueryRow(ctx, query, isbn).Scan(
		&b.Id,
		&b.Title,
		&b.Author,
		&b.ISBN,
		&b.Publisher,
		&b.PublishedYear,
		&b.Pages,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	assert.Nil(t, err)

	return &b
}

// CountBooks returns the total number of books
func (h *TestHelper) CountBooks(ctx context.Context, t *testing.T) int {
	var count int
	err := h.pool.QueryRow(ctx, "SELECT COUNT(*) FROM books").Scan(&count)
	assert.Nil(t, err)
	return count
}

// CountBooksByAuthor returns the number of books by a specific author
func (h *TestHelper) CountBooksByAuthor(ctx context.Context, t *testing.T, author string) int {
	var count int
	query := `SELECT COUNT(*) FROM books WHERE LOWER(author) LIKE LOWER($1)`
	err := h.pool.QueryRow(ctx, query, "%"+author+"%").Scan(&count)
	assert.Nil(t, err)
	return count
}

// GetAllBooks retrieves all books as a map indexed by book ID
func (h *TestHelper) GetAllBooks(ctx context.Context, t *testing.T) map[int]book.Book {
	query := `
		SELECT id, title, author, isbn, publisher, published_year, pages, description, created_at, updated_at
		FROM books
		ORDER BY id
	`

	rows, err := h.pool.Query(ctx, query)
	assert.Nil(t, err)
	defer rows.Close()

	books := make(map[int]book.Book)
	for rows.Next() {
		var b book.Book
		err := rows.Scan(
			&b.Id,
			&b.Title,
			&b.Author,
			&b.ISBN,
			&b.Publisher,
			&b.PublishedYear,
			&b.Pages,
			&b.Description,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		assert.Nil(t, err)
		books[b.Id] = b
	}

	assert.Nil(t, rows.Err())

	return books
}

// GetPool returns the connection pool
func (h *TestHelper) GetPool() *pgxpool.Pool {
	return h.pool
}

// BookExists checks if a book exists by ID
func (h *TestHelper) BookExists(ctx context.Context, t *testing.T, id int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)`
	err := h.pool.QueryRow(ctx, query, id).Scan(&exists)
	assert.Nil(t, err)
	return exists
}
