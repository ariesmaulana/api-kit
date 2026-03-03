package book_test

import (
	"context"
	"testing"

	"github.com/ariesmaulana/api-kit/src/app/book"
	"github.com/stretchr/testify/assert"
)

func TestBookCreate(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Book Create", func() {

			type input struct {
				title         string
				author        string
				isbn          string
				publisher     string
				publishedYear int
				pages         int
				description   string
			}
			type expected struct {
				success     bool
				message     string
				bookCreated bool
			}

			type testRow struct {
				name     string
				input    *input
				expected *expected
			}

			runtest := func(t *testing.T, app *BookApp, r *testRow) {
				ctx := context.Background()

				initialBooks := app.Helper.GetAllBooks(ctx, t)

				output := app.Service.CreateBook(ctx, &book.CreateBookInput{
					TraceId:       "trace-test",
					Title:         r.input.title,
					Author:        r.input.author,
					ISBN:          r.input.isbn,
					Publisher:     r.input.publisher,
					PublishedYear: r.input.publishedYear,
					Pages:         r.input.pages,
					Description:   r.input.description,
				})

				afterBooks := app.Helper.GetAllBooks(ctx, t)

				assert.Equal(t, r.expected.success, output.Success, r.name)
				assert.Equal(t, r.expected.message, output.Message, r.name)

				if r.expected.success == false {
					// Verify no books were created
					assert.Equal(t, initialBooks, afterBooks, r.name)
					return
				}

				// For successful creation
				assert.Equal(t, len(initialBooks)+1, len(afterBooks), r.name)
				assert.NotZero(t, output.Book.Id, r.name)
				assert.Equal(t, r.input.title, output.Book.Title, r.name)
				assert.Equal(t, r.input.author, output.Book.Author, r.name)
				assert.Equal(t, r.input.isbn, output.Book.ISBN, r.name)
				assert.Equal(t, r.input.publisher, output.Book.Publisher, r.name)
				assert.Equal(t, r.input.publishedYear, output.Book.PublishedYear, r.name)
				assert.Equal(t, r.input.pages, output.Book.Pages, r.name)
				assert.Equal(t, r.input.description, output.Book.Description, r.name)
				assert.NotZero(t, output.Book.CreatedAt, r.name)
				assert.NotZero(t, output.Book.UpdatedAt, r.name)
			}

			runRows := func(t *testing.T, app *BookApp, rows []*testRow) {
				for _, r := range rows {
					runtest(t, app, r)
				}
			}

			suite.Run(t, "Create Book scenarios", func(t *testing.T, ctx context.Context, app *BookApp) {
				runRows(t, app, []*testRow{
					// ===== Success Tests =====
					{
						name: "Should create book successfully with all fields",
						input: &input{
							title:         "The Go Programming Language",
							author:        "Alan A. A. Donovan",
							isbn:          "978-0134190440",
							publisher:     "Addison-Wesley",
							publishedYear: 2015,
							pages:         380,
							description:   "The authoritative resource to writing clear and idiomatic Go",
						},
						expected: &expected{
							success:     true,
							message:     "Book created successfully",
							bookCreated: true,
						},
					},
					{
						name: "Should create book with minimal required fields",
						input: &input{
							title:  "Minimal Book",
							author: "Minimal Author",
							isbn:   "978-1111111111",
						},
						expected: &expected{
							success:     true,
							message:     "Book created successfully",
							bookCreated: true,
						},
					},
					{
						name: "Should create book with long description",
						input: &input{
							title:       "Book with Long Description",
							author:      "Author Name",
							isbn:        "978-2222222222",
							description: "This is a very long description that contains many words and details about the book content. It should still be saved properly in the database without any issues.",
						},
						expected: &expected{
							success:     true,
							message:     "Book created successfully",
							bookCreated: true,
						},
					},

					// ===== Validation Tests: Title =====
					{
						name: "Should fail when title is empty",
						input: &input{
							title:  "",
							author: "Test Author",
							isbn:   "978-1234567890",
						},
						expected: &expected{
							success:     false,
							message:     "Title is mandatory",
							bookCreated: false,
						},
					},

					// ===== Validation Tests: Author =====
					{
						name: "Should fail when author is empty",
						input: &input{
							title:  "Test Book",
							author: "",
							isbn:   "978-1234567890",
						},
						expected: &expected{
							success:     false,
							message:     "Author is mandatory",
							bookCreated: false,
						},
					},

					// ===== Validation Tests: ISBN =====
					{
						name: "Should fail when ISBN is duplicate",
						input: &input{
							title:     "Duplicate ISBN Book",
							author:    "Another Author",
							isbn:      "978-0134190440", // Same as first book
							publisher: "Different Publisher",
						},
						expected: &expected{
							success:     false,
							message:     "Book with this ISBN already exists",
							bookCreated: false,
						},
					},
				})
			})

		})
	})
}

func TestBookGetById(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Book GetById", func() {

			var Books []DataBook

			type input struct {
				bookID int
			}
			type expected struct {
				success bool
				message string
				bookID  int
				title   string
				author  string
			}

			type testRow struct {
				name     string
				input    *input
				expected *expected
			}

			suite.Setup(func(ctx context.Context, app *BookApp) {
				Books = []DataBook{
					{
						Idx:           0,
						Title:         "Book One",
						Author:        "Author A",
						ISBN:          "978-1111111111",
						Publisher:     "Publisher A",
						PublishedYear: 2020,
						Pages:         200,
						Description:   "Description 1",
					},
					{
						Idx:           1,
						Title:         "Book Two",
						Author:        "Author B",
						ISBN:          "978-2222222222",
						Publisher:     "Publisher B",
						PublishedYear: 2021,
						Pages:         250,
						Description:   "Description 2",
					},
				}

				// Insert books and store actual database IDs
				for i, bookData := range Books {
					insertedBook := app.Helper.InsertBook(ctx, t, bookData.Title, bookData.Author, bookData.ISBN, bookData.Publisher, bookData.PublishedYear, bookData.Pages, bookData.Description)
					Books[i].Id = insertedBook.Id
				}
			})

			runtest := func(t *testing.T, app *BookApp, r *testRow) {
				ctx := context.Background()

				initialBooks := app.Helper.GetAllBooks(ctx, t)

				output := app.Service.GetBookById(ctx, &book.GetBookByIdInput{
					TraceId: "trace-test",
					Id:      r.input.bookID,
				})

				afterBooks := app.Helper.GetAllBooks(ctx, t)

				assert.Equal(t, r.expected.success, output.Success, r.name)
				assert.Equal(t, r.expected.message, output.Message, r.name)

				if r.expected.success == false {
					// Verify no books were modified (read operation)
					assert.Equal(t, initialBooks, afterBooks, r.name)
					return
				}

				// For successful retrieval
				assert.Equal(t, len(initialBooks), len(afterBooks), r.name)
				assert.NotZero(t, output.Book.Id, r.name)
				assert.Equal(t, r.expected.bookID, output.Book.Id, r.name)
				assert.Equal(t, r.expected.title, output.Book.Title, r.name)
				assert.Equal(t, r.expected.author, output.Book.Author, r.name)
				assert.NotZero(t, output.Book.CreatedAt, r.name)
				assert.NotZero(t, output.Book.UpdatedAt, r.name)

				// Verify no books were modified (read operation)
				assert.Equal(t, initialBooks, afterBooks, r.name)
			}

			runRows := func(t *testing.T, app *BookApp, rows []*testRow) {
				for _, r := range rows {
					runtest(t, app, r)
				}
			}

			suite.Run(t, "GetById scenarios", func(t *testing.T, ctx context.Context, app *BookApp) {
				runRows(t, app, []*testRow{
					// ===== Success Tests =====
					{
						name: "Should get book successfully for first book",
						input: &input{
							bookID: Books[0].Id,
						},
						expected: &expected{
							success: true,
							message: "Book retrieved successfully",
							bookID:  Books[0].Id,
							title:   "Book One",
							author:  "Author A",
						},
					},
					{
						name: "Should get book successfully for second book",
						input: &input{
							bookID: Books[1].Id,
						},
						expected: &expected{
							success: true,
							message: "Book retrieved successfully",
							bookID:  Books[1].Id,
							title:   "Book Two",
							author:  "Author B",
						},
					},

					// ===== Error Tests =====
					{
						name: "Should fail when book does not exist",
						input: &input{
							bookID: 99999,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
					{
						name: "Should fail when book ID is zero",
						input: &input{
							bookID: 0,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
					{
						name: "Should fail when book ID is negative",
						input: &input{
							bookID: -1,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
				})
			})

		})
	})
}

func TestBookGetBooks(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Book GetBooks", func() {

			var Books []DataBook

			type input struct {
				page   int
				limit  int
				author string
			}
			type expected struct {
				success   bool
				message   string
				minBooks  int
				maxBooks  int
				hasAuthor string
			}

			type testRow struct {
				name     string
				input    *input
				expected *expected
			}

			suite.Setup(func(ctx context.Context, app *BookApp) {
				Books = []DataBook{
					{
						Idx:           0,
						Title:         "Go Book 1",
						Author:        "Alan Donovan",
						ISBN:          "978-3333333331",
						Publisher:     "Publisher",
						PublishedYear: 2020,
						Pages:         200,
						Description:   "Go programming",
					},
					{
						Idx:           1,
						Title:         "Go Book 2",
						Author:        "Alan Donovan",
						ISBN:          "978-3333333332",
						Publisher:     "Publisher",
						PublishedYear: 2021,
						Pages:         250,
						Description:   "Advanced Go",
					},
					{
						Idx:           2,
						Title:         "Java Book",
						Author:        "Robert Martin",
						ISBN:          "978-3333333333",
						Publisher:     "Publisher",
						PublishedYear: 2022,
						Pages:         300,
						Description:   "Java programming",
					},
				}

				// Insert books
				for i, bookData := range Books {
					insertedBook := app.Helper.InsertBook(ctx, t, bookData.Title, bookData.Author, bookData.ISBN, bookData.Publisher, bookData.PublishedYear, bookData.Pages, bookData.Description)
					Books[i].Id = insertedBook.Id
				}
			})

			runtest := func(t *testing.T, app *BookApp, r *testRow) {
				ctx := context.Background()

				initialBooks := app.Helper.GetAllBooks(ctx, t)

				output := app.Service.GetBooks(ctx, &book.GetBooksInput{
					TraceId: "trace-test",
					Page:    r.input.page,
					Limit:   r.input.limit,
					Author:  r.input.author,
				})

				afterBooks := app.Helper.GetAllBooks(ctx, t)

				assert.Equal(t, r.expected.success, output.Success, r.name)
				assert.Equal(t, r.expected.message, output.Message, r.name)

				if r.expected.success == false {
					// Verify no books were modified (read operation)
					assert.Equal(t, initialBooks, afterBooks, r.name)
					return
				}

				// For successful retrieval
				assert.GreaterOrEqual(t, len(output.Books), r.expected.minBooks, r.name)
				assert.LessOrEqual(t, len(output.Books), r.expected.maxBooks, r.name)

				// If author filter is specified, verify all books match
				if r.expected.hasAuthor != "" {
					for _, b := range output.Books {
						assert.Contains(t, b.Author, r.expected.hasAuthor, r.name)
					}
				}

				// Verify no books were modified (read operation)
				assert.Equal(t, initialBooks, afterBooks, r.name)
			}

			runRows := func(t *testing.T, app *BookApp, rows []*testRow) {
				for _, r := range rows {
					runtest(t, app, r)
				}
			}

			suite.Run(t, "GetBooks scenarios", func(t *testing.T, ctx context.Context, app *BookApp) {
				runRows(t, app, []*testRow{
					// ===== Success Tests: Pagination =====
					{
						name: "Should get all books with default pagination",
						input: &input{
							page:  1,
							limit: 10,
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 3,
							maxBooks: 3,
						},
					},
					{
						name: "Should get first 2 books",
						input: &input{
							page:  1,
							limit: 2,
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 2,
							maxBooks: 2,
						},
					},
					{
						name: "Should get books from second page",
						input: &input{
							page:  2,
							limit: 2,
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 1,
							maxBooks: 1,
						},
					},
					{
						name: "Should use default limit when limit is 0",
						input: &input{
							page:  1,
							limit: 0,
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 3,
							maxBooks: 10,
						},
					},
					{
						name: "Should cap limit at maximum",
						input: &input{
							page:  1,
							limit: 200,
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 3,
							maxBooks: 3,
						},
					},

					// ===== Success Tests: Author Filter =====
					{
						name: "Should filter books by author 'Alan'",
						input: &input{
							page:   1,
							limit:  10,
							author: "Alan",
						},
						expected: &expected{
							success:   true,
							message:   "Books retrieved successfully",
							minBooks:  2,
							maxBooks:  2,
							hasAuthor: "Alan",
						},
					},
					{
						name: "Should filter books by author 'Martin'",
						input: &input{
							page:   1,
							limit:  10,
							author: "Martin",
						},
						expected: &expected{
							success:   true,
							message:   "Books retrieved successfully",
							minBooks:  1,
							maxBooks:  1,
							hasAuthor: "Martin",
						},
					},
					{
						name: "Should handle case-insensitive author search",
						input: &input{
							page:   1,
							limit:  10,
							author: "alan",
						},
						expected: &expected{
							success:   true,
							message:   "Books retrieved successfully",
							minBooks:  2,
							maxBooks:  2,
							hasAuthor: "Alan",
						},
					},
					{
						name: "Should return empty list for non-existent author",
						input: &input{
							page:   1,
							limit:  10,
							author: "NonExistent",
						},
						expected: &expected{
							success:  true,
							message:  "Books retrieved successfully",
							minBooks: 0,
							maxBooks: 0,
						},
					},

					// ===== Combined Tests =====
					{
						name: "Should filter by author with pagination",
						input: &input{
							page:   1,
							limit:  1,
							author: "Alan",
						},
						expected: &expected{
							success:   true,
							message:   "Books retrieved successfully",
							minBooks:  1,
							maxBooks:  1,
							hasAuthor: "Alan",
						},
					},
				})
			})

		})
	})
}

func TestBookUpdate(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Book Update", func() {

			var Books []DataBook

			type input struct {
				bookID        int
				title         string
				author        string
				isbn          string
				publisher     string
				publishedYear int
				pages         int
				description   string
			}
			type expected struct {
				success bool
				message string
			}

			type testRow struct {
				name     string
				input    *input
				expected *expected
			}

			suite.Setup(func(ctx context.Context, app *BookApp) {
				Books = []DataBook{
					{
						Idx:           0,
						Title:         "Original Book",
						Author:        "Original Author",
						ISBN:          "978-4444444441",
						Publisher:     "Original Publisher",
						PublishedYear: 2020,
						Pages:         200,
						Description:   "Original description",
					},
				}

				// Insert books
				for i, bookData := range Books {
					insertedBook := app.Helper.InsertBook(ctx, t, bookData.Title, bookData.Author, bookData.ISBN, bookData.Publisher, bookData.PublishedYear, bookData.Pages, bookData.Description)
					Books[i].Id = insertedBook.Id
				}
			})

			runtest := func(t *testing.T, app *BookApp, r *testRow) {
				ctx := context.Background()

				initialBooks := app.Helper.GetAllBooks(ctx, t)

				output := app.Service.UpdateBook(ctx, &book.UpdateBookInput{
					TraceId:       "trace-test",
					Id:            r.input.bookID,
					Title:         r.input.title,
					Author:        r.input.author,
					ISBN:          r.input.isbn,
					Publisher:     r.input.publisher,
					PublishedYear: r.input.publishedYear,
					Pages:         r.input.pages,
					Description:   r.input.description,
				})

				afterBooks := app.Helper.GetAllBooks(ctx, t)

				assert.Equal(t, r.expected.success, output.Success, r.name)
				assert.Equal(t, r.expected.message, output.Message, r.name)

				if r.expected.success == false {
					// Verify no books were modified
					assert.Equal(t, initialBooks, afterBooks, r.name)
					return
				}

				// For successful update
				assert.Equal(t, len(initialBooks), len(afterBooks), r.name)
				assert.NotZero(t, output.Book.Id, r.name)
				assert.Equal(t, r.input.title, output.Book.Title, r.name)
				// Other fields should remain unchanged from the original book
				originalBook := Books[0]
				assert.Equal(t, originalBook.Author, output.Book.Author, r.name)
				assert.Equal(t, originalBook.ISBN, output.Book.ISBN, r.name)
				assert.Equal(t, originalBook.Publisher, output.Book.Publisher, r.name)
				assert.Equal(t, originalBook.PublishedYear, output.Book.PublishedYear, r.name)
				assert.Equal(t, originalBook.Pages, output.Book.Pages, r.name)
				assert.Equal(t, originalBook.Description, output.Book.Description, r.name)
				assert.NotZero(t, output.Book.CreatedAt, r.name)
				assert.NotZero(t, output.Book.UpdatedAt, r.name)

				// Verify the book was actually updated
				updatedBook, exists := afterBooks[r.input.bookID]
				assert.True(t, exists, r.name+" - book should exist in afterBooks")
				assert.Equal(t, r.input.title, updatedBook.Title, r.name)
			}

			runRows := func(t *testing.T, app *BookApp, rows []*testRow) {
				for _, r := range rows {
					runtest(t, app, r)
				}
			}

			suite.Run(t, "Update scenarios", func(t *testing.T, ctx context.Context, app *BookApp) {
				runRows(t, app, []*testRow{
					// ===== Success Tests =====
					{
						name: "Should update book successfully",
						input: &input{
							bookID: Books[0].Id,
							title:  "Updated Title",
						},
						expected: &expected{
							success: true,
							message: "Book updated successfully",
						},
					},
					{
						name: "Should update only title",
						input: &input{
							bookID: Books[0].Id,
							title:  "New Title",
						},
						expected: &expected{
							success: true,
							message: "Book updated successfully",
						},
					},

					// ===== Validation Tests: Title =====
					{
						name: "Should fail when title is empty",
						input: &input{
							bookID: Books[0].Id,
							title:  "",
							author: "Author",
						},
						expected: &expected{
							success: false,
							message: "Title is mandatory",
						},
					},

					// ===== Validation Tests: Book ID =====
					{
						name: "Should fail when book does not exist",
						input: &input{
							bookID: 99999,
							title:  "Title",
							author: "Author",
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
				})
			})

		})
	})
}

func TestBookDelete(t *testing.T) {
	RunTest(t, func(t *testing.T, suite *TestSuite) {
		suite.Describe(t, "Book Delete", func() {

			var Books []DataBook

			type input struct {
				bookID int
			}
			type expected struct {
				success bool
				message string
			}

			type testRow struct {
				name     string
				input    *input
				expected *expected
			}

			suite.Setup(func(ctx context.Context, app *BookApp) {
				Books = []DataBook{
					{
						Idx:           0,
						Title:         "Book to Delete",
						Author:        "Author",
						ISBN:          "978-5555555551",
						Publisher:     "Publisher",
						PublishedYear: 2020,
						Pages:         200,
						Description:   "Description",
					},
					{
						Idx:           1,
						Title:         "Another Book",
						Author:        "Author",
						ISBN:          "978-5555555552",
						Publisher:     "Publisher",
						PublishedYear: 2021,
						Pages:         250,
						Description:   "Description",
					},
				}

				// Insert books
				for i, bookData := range Books {
					insertedBook := app.Helper.InsertBook(ctx, t, bookData.Title, bookData.Author, bookData.ISBN, bookData.Publisher, bookData.PublishedYear, bookData.Pages, bookData.Description)
					Books[i].Id = insertedBook.Id
				}
			})

			runtest := func(t *testing.T, app *BookApp, r *testRow) {
				ctx := context.Background()

				initialBooks := app.Helper.GetAllBooks(ctx, t)

				output := app.Service.DeleteBook(ctx, &book.DeleteBookInput{
					TraceId: "trace-test",
					Id:      r.input.bookID,
				})

				afterBooks := app.Helper.GetAllBooks(ctx, t)

				assert.Equal(t, r.expected.success, output.Success, r.name)
				assert.Equal(t, r.expected.message, output.Message, r.name)

				if r.expected.success == false {
					// Verify no books were modified
					assert.Equal(t, initialBooks, afterBooks, r.name)
					return
				}

				// For successful deletion
				assert.Equal(t, len(initialBooks)-1, len(afterBooks), r.name)

				// Verify the book was actually deleted
				_, exists := afterBooks[r.input.bookID]
				assert.False(t, exists, r.name+" - book should not exist in afterBooks")
			}

			runRows := func(t *testing.T, app *BookApp, rows []*testRow) {
				for _, r := range rows {
					runtest(t, app, r)
				}
			}

			suite.Run(t, "Delete scenarios", func(t *testing.T, ctx context.Context, app *BookApp) {
				runRows(t, app, []*testRow{
					// ===== Success Tests =====
					{
						name: "Should delete book successfully",
						input: &input{
							bookID: Books[0].Id,
						},
						expected: &expected{
							success: true,
							message: "Book deleted successfully",
						},
					},
					{
						name: "Should delete another book successfully",
						input: &input{
							bookID: Books[1].Id,
						},
						expected: &expected{
							success: true,
							message: "Book deleted successfully",
						},
					},

					// ===== Error Tests =====
					{
						name: "Should fail when book does not exist",
						input: &input{
							bookID: 99999,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
					{
						name: "Should fail when book ID is zero",
						input: &input{
							bookID: 0,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
					{
						name: "Should fail when book ID is negative",
						input: &input{
							bookID: -1,
						},
						expected: &expected{
							success: false,
							message: "Book not found",
						},
					},
				})
			})

		})
	})
}
