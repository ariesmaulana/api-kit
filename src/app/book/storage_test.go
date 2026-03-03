package book_test

import (
	"context"
	"testing"

	"github.com/ariesmaulana/api-kit/src/app/book"
	testsuite "github.com/ariesmaulana/api-kit/testing"
	"github.com/stretchr/testify/assert"
)

func TestStorageInsertBook(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage InsertBook", func() {
			suite.Runs(t, "Should insert book successfully", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				id, errType, err := tx.InsertBook(ctx, "Test Book", "Test Author", "978-1234567890", "Test Publisher", 2024, 300, "Test description")
				assert.Nil(t, err)
				assert.NotZero(t, id)
				assert.Equal(t, book.ErrTypeNone, errType)

				err = tx.Commit()
				assert.Nil(t, err)

				insertedBook := app.Helper.GetBookById(ctx, t, id)
				assert.Equal(t, "Test Book", insertedBook.Title)
				assert.Equal(t, "Test Author", insertedBook.Author)
				assert.Equal(t, "978-1234567890", insertedBook.ISBN)
				assert.Equal(t, "Test Publisher", insertedBook.Publisher)
				assert.Equal(t, 2024, insertedBook.PublishedYear)
				assert.Equal(t, 300, insertedBook.Pages)
				assert.Equal(t, "Test description", insertedBook.Description)
			})

			suite.Runs(t, "Should return error on duplicate ISBN", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				// Insert first book
				existingBook := app.Helper.InsertBook(ctx, t, "Existing Book", "Author", "978-0000000001", "Publisher", 2020, 200, "Description")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				// Try to insert book with same ISBN
				_, errType, err := tx.InsertBook(ctx, "New Book", "Different Author", existingBook.ISBN, "New Publisher", 2024, 300, "New description")
				assert.NotNil(t, err)
				assert.Equal(t, book.ErrTypeUniqueConstraint, errType)
			})
		})
	})
}

func TestStorageGetBookById(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage GetBookById", func() {
			suite.Runs(t, "Should get existing book", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				existingBook := app.Helper.InsertBook(ctx, t, "Existing Book", "Author", "978-1111111111", "Publisher", 2020, 250, "Description")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				retrievedBook, err := tx.GetBookById(ctx, int(existingBook.Id))
				assert.Nil(t, err)
				assert.Equal(t, existingBook.Id, retrievedBook.Id)
				assert.Equal(t, existingBook.Title, retrievedBook.Title)
				assert.Equal(t, existingBook.Author, retrievedBook.Author)
				assert.Equal(t, existingBook.ISBN, retrievedBook.ISBN)
				assert.Equal(t, existingBook.Publisher, retrievedBook.Publisher)
				assert.Equal(t, existingBook.PublishedYear, retrievedBook.PublishedYear)
				assert.Equal(t, existingBook.Pages, retrievedBook.Pages)
				assert.Equal(t, existingBook.Description, retrievedBook.Description)
			})

			suite.Runs(t, "Should return error for non-existent book", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				_, err = tx.GetBookById(ctx, 99999)
				assert.NotNil(t, err)
			})
		})
	})
}

func TestStorageGetBooks(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage GetBooks", func() {
			suite.Runs(t, "Should get all books with pagination", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				// Insert multiple books
				app.Helper.InsertBook(ctx, t, "Book 1", "Author A", "978-1111111111", "Publisher", 2020, 200, "Desc 1")
				app.Helper.InsertBook(ctx, t, "Book 2", "Author B", "978-2222222222", "Publisher", 2021, 250, "Desc 2")
				app.Helper.InsertBook(ctx, t, "Book 3", "Author C", "978-3333333333", "Publisher", 2022, 300, "Desc 3")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				// Get first 2 books
				books, err := tx.GetBooks(ctx, 2, 0, "")
				assert.Nil(t, err)
				assert.Equal(t, 2, len(books))
			})

			suite.Runs(t, "Should filter books by author", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				// Insert books with different authors
				app.Helper.InsertBook(ctx, t, "Go Book 1", "Alan Donovan", "978-4444444444", "Publisher", 2020, 200, "Desc")
				app.Helper.InsertBook(ctx, t, "Go Book 2", "Alan Donovan", "978-5555555555", "Publisher", 2021, 250, "Desc")
				app.Helper.InsertBook(ctx, t, "Java Book", "Robert Martin", "978-6666666666", "Publisher", 2022, 300, "Desc")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				// Search for books by "Alan"
				books, err := tx.GetBooks(ctx, 10, 0, "Alan")
				assert.Nil(t, err)
				assert.Equal(t, 2, len(books))
				for _, b := range books {
					assert.Contains(t, b.Author, "Alan")
				}
			})

			suite.Runs(t, "Should handle case-insensitive author search", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				app.Helper.InsertBook(ctx, t, "Book", "Martin Fowler", "978-7777777777", "Publisher", 2020, 200, "Desc")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				// Search with lowercase
				books, err := tx.GetBooks(ctx, 10, 0, "martin")
				assert.Nil(t, err)
				assert.Equal(t, 1, len(books))
			})

			suite.Runs(t, "Should return empty list when no books match", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				app.Helper.InsertBook(ctx, t, "Book", "Author", "978-8888888888", "Publisher", 2020, 200, "Desc")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				books, err := tx.GetBooks(ctx, 10, 0, "NonExistentAuthor")
				assert.Nil(t, err)
				assert.Equal(t, 0, len(books))
			})
		})
	})
}

func TestStorageUpdateBook(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage UpdateBook", func() {
			suite.Runs(t, "Should update book successfully", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				existingBook := app.Helper.InsertBook(ctx, t, "Old Title", "Old Author", "978-9999999999", "Old Publisher", 2020, 200, "Old description")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				err = tx.UpdateBook(ctx, existingBook.Id, "New Title")
				assert.Nil(t, err)

				err = tx.Commit()
				assert.Nil(t, err)

				updatedBook := app.Helper.GetBookById(ctx, t, existingBook.Id)
				assert.Equal(t, "New Title", updatedBook.Title)
			})

			// Note: PostgreSQL UPDATE doesn't error on 0 rows affected
			// Would need to check affected rows count to validate
		})
	})
}

func TestStorageDeleteBook(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage DeleteBook", func() {
			suite.Runs(t, "Should delete book successfully", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				existingBook := app.Helper.InsertBook(ctx, t, "Book to Delete", "Author", "978-0000000000", "Publisher", 2020, 200, "Description")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				err = tx.DeleteBook(ctx, existingBook.Id)
				assert.Nil(t, err)

				err = tx.Commit()
				assert.Nil(t, err)

				bookExists := app.Helper.BookExists(ctx, t, existingBook.Id)
				assert.False(t, bookExists)
			})

			suite.Runs(t, "Should return error when deleting non-existent book", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				err = tx.DeleteBook(ctx, 99999)
				assert.NotNil(t, err)
			})
		})
	})
}

func TestStorageLockBookById(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage LockBookById", func() {
			suite.Runs(t, "Should lock existing book", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				existingBook := app.Helper.InsertBook(ctx, t, "Book to Lock", "Author", "978-1010101010", "Publisher", 2020, 200, "Description")

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				lockedBook, errType, err := tx.LockBookById(ctx, existingBook.Id)
				assert.Nil(t, err)
				assert.Equal(t, book.ErrTypeNone, errType)
				assert.Equal(t, existingBook.Id, lockedBook.Id)
				assert.Equal(t, existingBook.Title, lockedBook.Title)
			})

			suite.Runs(t, "Should return error when locking non-existent book", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				_, errType, err := tx.LockBookById(ctx, 99999)
				assert.NotNil(t, err)
				assert.Equal(t, book.ErrTypeNotFound, errType)
			})
		})
	})
}

func TestStorageTransactionRollback(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage Transaction Rollback", func() {
			suite.Runs(t, "Should rollback changes when transaction is not committed", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				initialCount := app.Helper.CountBooks(ctx, t)

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)

				_, _, err = tx.InsertBook(ctx, "Rollback Test", "Author", "978-1212121212", "Publisher", 2024, 300, "Description")
				assert.Nil(t, err)

				tx.Rollback()

				finalCount := app.Helper.CountBooks(ctx, t)
				assert.Equal(t, initialCount, finalCount)
			})
		})
	})
}

func TestStorageTransactionCommit(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Storage Transaction Commit", func() {
			suite.Runs(t, "Should persist changes when transaction is committed", func(t *testing.T, appCtx *testsuite.AppContext) {
				app := initBookApp(appCtx)
				ctx := context.Background()

				initialCount := app.Helper.CountBooks(ctx, t)

				tx, err := app.Storage.BeginTx(ctx)
				assert.Nil(t, err)
				defer tx.Rollback()

				_, _, err = tx.InsertBook(ctx, "Commit Test", "Author", "978-1313131313", "Publisher", 2024, 300, "Description")
				assert.Nil(t, err)

				err = tx.Commit()
				assert.Nil(t, err)

				finalCount := app.Helper.CountBooks(ctx, t)
				assert.Equal(t, initialCount+1, finalCount)
			})
		})
	})
}
