CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    publisher VARCHAR(255),
    published_year INTEGER,
    pages INTEGER,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for fast lookups of books by title
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);

-- Index for fast lookups of books by author
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);

-- Index for fast lookups of books by ISBN
CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
