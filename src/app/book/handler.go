package book

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ariesmaulana/api-kit/src/app/user"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// Handler handles HTTP requests for book operations
type Handler struct {
	service    Service
	jwtService *user.JWTService
}

// NewHandler creates a new book handler
func NewHandler(service Service, jwtService *user.JWTService) *Handler {
	return &Handler{
		service:    service,
		jwtService: jwtService,
	}
}

// RegisterRoutes registers all book routes
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/books")

	// Public routes
	g.POST("", h.CreateBook)
	g.GET("/:id", h.GetBookById)

	// Protected routes (requires authentication)
	protected := g.Group("")
	protected.Use(h.jwtService.JWTMiddleware())
	protected.GET("", h.GetBooks)
	protected.PUT("/:id", h.UpdateBook)
	protected.DELETE("/:id", h.DeleteBook)
}

// HTTP Request/Response structs

// CreateBookRequest represents the HTTP request body for creating a book
type CreateBookRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	ISBN          string `json:"isbn"`
	Publisher     string `json:"publisher"`
	PublishedYear int    `json:"published_year"`
	Pages         int    `json:"pages"`
	Description   string `json:"description"`
}

// UpdateBookRequest represents the HTTP request body for updating a book
type UpdateBookRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	ISBN          string `json:"isbn"`
	Publisher     string `json:"publisher"`
	PublishedYear int    `json:"published_year"`
	Pages         int    `json:"pages"`
	Description   string `json:"description"`
}

// BookDTO represents a book in HTTP responses
type BookDTO struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	Publisher     string    `json:"publisher"`
	PublishedYear int       `json:"published_year"`
	Pages         int       `json:"pages"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BookResponse represents the HTTP response for single book operations
type BookResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    *BookDTO `json:"data,omitempty"`
}

// BooksResponse represents the HTTP response for multiple books operations
type BooksResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    []BookDTO  `json:"data,omitempty"`
	Total   int        `json:"total,omitempty"`
}

func toBookDTO(book Book) BookDTO {
	return BookDTO{
		Id:            book.Id,
		Title:         book.Title,
		Author:        book.Author,
		ISBN:          book.ISBN,
		Publisher:     book.Publisher,
		PublishedYear: book.PublishedYear,
		Pages:         book.Pages,
		Description:   book.Description,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
}

func toBooksDTO(books []Book) []BookDTO {
	dtos := make([]BookDTO, len(books))
	for i, book := range books {
		dtos[i] = toBookDTO(book)
	}
	return dtos
}

// CreateBook handles POST /books
// @Summary Create a new book
// @Description Create a new book entry
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "Book creation data"
// @Success 201 {object} BookResponse
// @Failure 400 {object} BookResponse
// @Router /books [post]
func (h *Handler) CreateBook(c echo.Context) error {
	traceID := xid.New().String()

	var req CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Title is required",
		})
	}

	if req.Author == "" {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Author is required",
		})
	}

	output := h.service.CreateBook(c.Request().Context(), &CreateBookInput{
		TraceId:       traceID,
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		Publisher:     req.Publisher,
		PublishedYear: req.PublishedYear,
		Pages:         req.Pages,
		Description:   req.Description,
	})

	if !output.Success {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: output.Message,
		})
	}

	dto := toBookDTO(output.Book)
	return c.JSON(http.StatusCreated, BookResponse{
		Success: true,
		Message: "Book created successfully",
		Data:    &dto,
	})
}

// GetBookById handles GET /books/:id
// @Summary Get book by ID
// @Description Get a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} BookResponse
// @Failure 400 {object} BookResponse
// @Failure 404 {object} BookResponse
// @Router /books/{id} [get]
func (h *Handler) GetBookById(c echo.Context) error {
	traceID := xid.New().String()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Invalid book ID",
		})
	}

	output := h.service.GetBookById(c.Request().Context(), &GetBookByIdInput{
		TraceId: traceID,
		Id:      id,
	})

	if !output.Success {
		status := http.StatusInternalServerError
		if output.Message == "Book not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, BookResponse{
			Success: false,
			Message: output.Message,
		})
	}

	dto := toBookDTO(output.Book)
	return c.JSON(http.StatusOK, BookResponse{
		Success: true,
		Message: output.Message,
		Data:    &dto,
	})
}

// GetBooks handles GET /books
// @Summary Get books with pagination and search
// @Description Get books with pagination support and optional author search
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (1-based)" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Param author query string false "Filter by author name"
// @Success 200 {object} BooksResponse
// @Failure 401 {object} BooksResponse
// @Failure 500 {object} BooksResponse
// @Router /books [get]
func (h *Handler) GetBooks(c echo.Context) error {
	traceID := xid.New().String()

	// Parse query parameters
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")
	authorParam := c.QueryParam("author")

	page := 1 // default to first page
	if pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10 // default
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	output := h.service.GetBooks(c.Request().Context(), &GetBooksInput{
		TraceId: traceID,
		Page:    page,
		Limit:   limit,
		Author:  authorParam,
	})

	if !output.Success {
		return c.JSON(http.StatusInternalServerError, BooksResponse{
			Success: false,
			Message: output.Message,
		})
	}

	dtos := toBooksDTO(output.Books)
	return c.JSON(http.StatusOK, BooksResponse{
		Success: true,
		Message: output.Message,
		Data:    dtos,
		Total:   output.Total,
	})
}

// UpdateBook handles PUT /books/:id
// @Summary Update book
// @Description Update a book's information
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param book body UpdateBookRequest true "Book update data"
// @Success 200 {object} BookResponse
// @Failure 400 {object} BookResponse
// @Failure 401 {object} BookResponse
// @Failure 404 {object} BookResponse
// @Router /books/{id} [put]
func (h *Handler) UpdateBook(c echo.Context) error {
	traceID := xid.New().String()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Invalid book ID",
		})
	}

	var req UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	output := h.service.UpdateBook(c.Request().Context(), &UpdateBookInput{
		TraceId:       traceID,
		Id:            id,
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		Publisher:     req.Publisher,
		PublishedYear: req.PublishedYear,
		Pages:         req.Pages,
		Description:   req.Description,
	})

	if !output.Success {
		status := http.StatusBadRequest
		if output.Message == "Book not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, BookResponse{
			Success: false,
			Message: output.Message,
		})
	}

	dto := toBookDTO(output.Book)
	return c.JSON(http.StatusOK, BookResponse{
		Success: true,
		Message: "Book updated successfully",
		Data:    &dto,
	})
}

// DeleteBook handles DELETE /books/:id
// @Summary Delete book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} BookResponse
// @Failure 400 {object} BookResponse
// @Failure 401 {object} BookResponse
// @Failure 404 {object} BookResponse
// @Router /books/{id} [delete]
func (h *Handler) DeleteBook(c echo.Context) error {
	traceID := xid.New().String()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BookResponse{
			Success: false,
			Message: "Invalid book ID",
		})
	}

	output := h.service.DeleteBook(c.Request().Context(), &DeleteBookInput{
		TraceId: traceID,
		Id:      id,
	})

	if !output.Success {
		status := http.StatusBadRequest
		if output.Message == "Book not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, BookResponse{
			Success: false,
			Message: output.Message,
		})
	}

	return c.JSON(http.StatusOK, BookResponse{
		Success: true,
		Message: "Book deleted successfully",
	})
}
