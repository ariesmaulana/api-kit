# Books API Testing Guide

Complete testing guide for the Books API with all scenarios.

## Prerequisites

1. Start the API server:
```bash
cd src
go run main.go
```

2. Create a test user (if not exists):
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "password": "password123"
  }'
```

3. Login and get token:
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }' | jq -r '.token'
```

Save the token:
```bash
export TOKEN="<your_token_here>"
```

## 1. CRUD: Create & Read

### Create Book (Public)
```bash
# Test 1: Create first book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan",
    "isbn": "978-0134190440",
    "publisher": "Addison-Wesley",
    "published_year": 2015,
    "pages": 380,
    "description": "The authoritative resource to writing clear and idiomatic Go"
  }' | jq

# Expected: 201 Created
# {
#   "success": true,
#   "message": "Book created successfully",
#   "data": { "id": 1, ... }
# }

# Test 2: Create second book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "isbn": "978-0132350884",
    "publisher": "Prentice Hall",
    "published_year": 2008,
    "pages": 464,
    "description": "A Handbook of Agile Software Craftsmanship"
  }' | jq

# Test 3: Create third book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Design Patterns",
    "author": "Gang of Four",
    "isbn": "978-0201633610",
    "publisher": "Addison-Wesley",
    "published_year": 1994,
    "pages": 395,
    "description": "Elements of Reusable Object-Oriented Software"
  }' | jq

# Test 4: Create book with same author as first book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Programming in Go",
    "author": "Alan A. A. Donovan",
    "isbn": "978-1234567890",
    "publisher": "Tech Press",
    "published_year": 2020,
    "pages": 300,
    "description": "Advanced Go programming"
  }' | jq
```

### Read Book by ID (Public)
```bash
# Test 5: Get book by ID
curl http://localhost:8080/books/1 | jq

# Expected: 200 OK with book data

# Test 6: Get another book
curl http://localhost:8080/books/2 | jq
```

### List Books (Protected - Requires Token)
```bash
# Test 7: List books (first page, default limit)
curl http://localhost:8080/books \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: 200 OK with array of books

# Test 8: List books with custom pagination
curl "http://localhost:8080/books?page=1&limit=2" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: First 2 books
```

## 2. CRUD: Update & Delete

### Update Book (Protected)
```bash
# Test 9: Update book
curl -X PUT http://localhost:8080/books/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language (2nd Edition)",
    "author": "Alan A. A. Donovan",
    "isbn": "978-0134190440",
    "publisher": "Addison-Wesley Professional",
    "published_year": 2016,
    "pages": 400,
    "description": "Updated edition with new content"
  }' | jq

# Expected: 200 OK with updated book data

# Verify update
curl http://localhost:8080/books/1 | jq
```

### Delete Book (Protected)
```bash
# Test 10: Delete book
curl -X DELETE http://localhost:8080/books/3 \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: 200 OK
# {
#   "success": true,
#   "message": "Book deleted successfully"
# }

# Verify deletion
curl http://localhost:8080/books/3 | jq

# Expected: 404 Not Found
```

## 3. Auth Guard

### Test Without Token (Should Fail)
```bash
# Test 11: Try to list books without token
curl http://localhost:8080/books | jq

# Expected: 401 Unauthorized
# {
#   "success": false,
#   "message": "Missing authentication token"
# }

# Test 12: Try to update without token
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","author":"Test"}' | jq

# Expected: 401 Unauthorized

# Test 13: Try to delete without token
curl -X DELETE http://localhost:8080/books/1 | jq

# Expected: 401 Unauthorized
```

### Test With Invalid Token
```bash
# Test 14: Invalid token
curl http://localhost:8080/books \
  -H "Authorization: Bearer invalid_token_here" | jq

# Expected: 401 Unauthorized
# {
#   "success": false,
#   "message": "Invalid authentication token"
# }
```

### Get Fresh Token
```bash
# Test 15: Get new token
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }' | jq

# Save the token
export TOKEN="<new_token>"
```

## 4. Search & Paginate

### Search by Author
```bash
# Test 16: Search for books by "Alan"
curl "http://localhost:8080/books?author=Alan" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Books by "Alan A. A. Donovan" (2 books if you created both)

# Test 17: Search for books by "Martin"
curl "http://localhost:8080/books?author=Martin" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Books by "Robert C. Martin"

# Test 18: Search with case insensitive
curl "http://localhost:8080/books?author=martin" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Same results as Test 17

# Test 19: Search with partial match
curl "http://localhost:8080/books?author=Dono" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Books by "Alan A. A. Donovan"
```

### Pagination
```bash
# Test 20: First page with 2 items
curl "http://localhost:8080/books?page=1&limit=2" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: First 2 books, total: 2

# Test 21: Second page with 2 items
curl "http://localhost:8080/books?page=2&limit=2" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Next 2 books (or less if fewer books exist)

# Test 22: Page 1 with limit 1
curl "http://localhost:8080/books?page=1&limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Just 1 book

# Test 23: Large limit (should cap at max)
curl "http://localhost:8080/books?page=1&limit=200" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: All books (limit capped at 100)
```

### Combined Search and Pagination
```bash
# Test 24: Search by author with pagination
curl "http://localhost:8080/books?author=Alan&page=1&limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: First book by Alan (if multiple exist)

# Test 25: Second page of search results
curl "http://localhost:8080/books?author=Alan&page=2&limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: Second book by Alan (if exists)
```

## 5. Error Handling

### Invalid Input - Create Book
```bash
# Test 26: Missing title
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Test Author",
    "isbn": "123456"
  }' | jq

# Expected: 400 Bad Request
# {
#   "success": false,
#   "message": "Title is required"
# }

# Test 27: Missing author
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Book",
    "isbn": "123456"
  }' | jq

# Expected: 400 Bad Request
# {
#   "success": false,
#   "message": "Author is required"
# }

# Test 28: Invalid JSON
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{invalid json}' | jq

# Expected: 400 Bad Request
# {
#   "success": false,
#   "message": "Invalid request body"
# }

# Test 29: Duplicate ISBN
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Duplicate Book",
    "author": "Test Author",
    "isbn": "978-0134190440"
  }' | jq

# Expected: 400 Bad Request
# {
#   "success": false,
#   "message": "Book with this ISBN already exists"
# }
```

### Not Found Errors
```bash
# Test 30: Get non-existent book
curl http://localhost:8080/books/99999 | jq

# Expected: 404 Not Found (or 500 if not properly handled)
# {
#   "success": false,
#   "message": "Book not found"
# }

# Test 31: Update non-existent book
curl -X PUT http://localhost:8080/books/99999 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test",
    "author": "Test"
  }' | jq

# Expected: 404 Not Found
# {
#   "success": false,
#   "message": "Book not found"
# }

# Test 32: Delete non-existent book
curl -X DELETE http://localhost:8080/books/99999 \
  -H "Authorization: Bearer $TOKEN" | jq

# Expected: 404 Not Found
```

### Invalid ID Format
```bash
# Test 33: Invalid book ID (not a number)
curl http://localhost:8080/books/abc | jq

# Expected: 400 Bad Request
# {
#   "success": false,
#   "message": "Invalid book ID"
# }
```

## Summary Test Script

Run all tests at once:

```bash
#!/bin/bash

echo "=== Books API Test Suite ==="
echo ""

# Setup
echo "1. Creating test user..."
curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","full_name":"Test User","password":"password123"}' > /dev/null

echo "2. Getting token..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' | jq -r '.token')

echo "Token: ${TOKEN:0:20}..."
echo ""

# Create books
echo "3. Creating books..."
BOOK1=$(curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Programming Language","author":"Alan A. A. Donovan","isbn":"978-0134190440"}' | jq -r '.data.id')

BOOK2=$(curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code","author":"Robert C. Martin","isbn":"978-0132350884"}' | jq -r '.data.id')

echo "Created books: $BOOK1, $BOOK2"
echo ""

# Test read
echo "4. Testing GET /books/:id (public)..."
curl -s http://localhost:8080/books/$BOOK1 | jq '.success'

# Test list (protected)
echo "5. Testing GET /books (protected)..."
curl -s http://localhost:8080/books -H "Authorization: Bearer $TOKEN" | jq '.success'

# Test search
echo "6. Testing search..."
curl -s "http://localhost:8080/books?author=Alan" -H "Authorization: Bearer $TOKEN" | jq '.total'

# Test pagination
echo "7. Testing pagination..."
curl -s "http://localhost:8080/books?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | jq '.total'

# Test update
echo "8. Testing PUT /books/:id (protected)..."
curl -s -X PUT http://localhost:8080/books/$BOOK1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","author":"Alan A. A. Donovan","isbn":"978-0134190440"}' | jq '.success'

# Test delete
echo "9. Testing DELETE /books/:id (protected)..."
curl -s -X DELETE http://localhost:8080/books/$BOOK2 \
  -H "Authorization: Bearer $TOKEN" | jq '.success'

# Test error cases
echo "10. Testing error: missing title..."
curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"author":"Test"}' | jq '.message'

echo "11. Testing error: not found..."
curl -s http://localhost:8080/books/99999 | jq '.message'

echo "12. Testing error: no auth..."
curl -s http://localhost:8080/books | jq '.message'

echo ""
echo "=== Tests Complete ==="
```

Save this as `test_books_api.sh`, make it executable, and run:
```bash
chmod +x test_books_api.sh
./test_books_api.sh
```
