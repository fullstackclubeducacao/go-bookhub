-- name: CreateBook :one
INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetBookByID :one
SELECT * FROM books WHERE id = $1;

-- name: GetBookByISBN :one
SELECT * FROM books WHERE isbn = $1;

-- name: ListBooks :many
SELECT * FROM books
ORDER BY title ASC
LIMIT $1 OFFSET $2;

-- name: ListAvailableBooks :many
SELECT * FROM books
WHERE available_copies > 0
ORDER BY title ASC
LIMIT $1 OFFSET $2;

-- name: CountBooks :one
SELECT COUNT(*) FROM books;

-- name: CountAvailableBooks :one
SELECT COUNT(*) FROM books WHERE available_copies > 0;

-- name: UpdateBook :one
UPDATE books
SET title = $2, author = $3, isbn = $4, published_year = $5,
    total_copies = $6, available_copies = $7, updated_at = $8
WHERE id = $1
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1;
