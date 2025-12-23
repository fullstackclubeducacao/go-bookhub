-- name: CreateLoan :one
INSERT INTO loans (id, user_id, book_id, borrowed_at, due_date, returned_at, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetLoanByID :one
SELECT * FROM loans WHERE id = $1;

-- name: GetLoanByIDWithDetails :one
SELECT
    l.*,
    u.name as user_name,
    b.title as book_title
FROM loans l
JOIN users u ON l.user_id = u.id
JOIN books b ON l.book_id = b.id
WHERE l.id = $1;

-- name: GetActiveByUserAndBook :one
SELECT * FROM loans
WHERE user_id = $1 AND book_id = $2 AND status = 'active';

-- name: ListLoans :many
SELECT * FROM loans
ORDER BY borrowed_at DESC
LIMIT $1 OFFSET $2;

-- name: ListLoansWithDetails :many
SELECT
    l.*,
    u.name as user_name,
    b.title as book_title
FROM loans l
JOIN users u ON l.user_id = u.id
JOIN books b ON l.book_id = b.id
ORDER BY l.borrowed_at DESC
LIMIT $1 OFFSET $2;

-- name: ListLoansByUser :many
SELECT * FROM loans
WHERE user_id = $1
ORDER BY borrowed_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLoansByUserWithDetails :many
SELECT
    l.*,
    u.name as user_name,
    b.title as book_title
FROM loans l
JOIN users u ON l.user_id = u.id
JOIN books b ON l.book_id = b.id
WHERE l.user_id = $1
ORDER BY l.borrowed_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLoansByStatus :many
SELECT * FROM loans
WHERE status = $1
ORDER BY borrowed_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLoansByStatusWithDetails :many
SELECT
    l.*,
    u.name as user_name,
    b.title as book_title
FROM loans l
JOIN users u ON l.user_id = u.id
JOIN books b ON l.book_id = b.id
WHERE l.status = $1
ORDER BY l.borrowed_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLoansByUserAndStatus :many
SELECT * FROM loans
WHERE user_id = $1 AND status = $2
ORDER BY borrowed_at DESC
LIMIT $3 OFFSET $4;

-- name: ListLoansByUserAndStatusWithDetails :many
SELECT
    l.*,
    u.name as user_name,
    b.title as book_title
FROM loans l
JOIN users u ON l.user_id = u.id
JOIN books b ON l.book_id = b.id
WHERE l.user_id = $1 AND l.status = $2
ORDER BY l.borrowed_at DESC
LIMIT $3 OFFSET $4;

-- name: CountLoans :one
SELECT COUNT(*) FROM loans;

-- name: CountLoansByUser :one
SELECT COUNT(*) FROM loans WHERE user_id = $1;

-- name: CountLoansByStatus :one
SELECT COUNT(*) FROM loans WHERE status = $1;

-- name: CountLoansByUserAndStatus :one
SELECT COUNT(*) FROM loans WHERE user_id = $1 AND status = $2;

-- name: UpdateLoan :one
UPDATE loans
SET returned_at = $2, status = $3
WHERE id = $1
RETURNING *;
