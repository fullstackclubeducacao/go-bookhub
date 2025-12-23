CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL,
    isbn VARCHAR(13) NOT NULL UNIQUE,
    published_year INTEGER,
    total_copies INTEGER NOT NULL DEFAULT 1,
    available_copies INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_copies CHECK (available_copies >= 0 AND available_copies <= total_copies)
);

CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
CREATE INDEX IF NOT EXISTS idx_books_available ON books(available_copies);
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);

-- Insert sample books (only if they don't exist)
INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies)
SELECT gen_random_uuid(), 'Clean Code', 'Robert C. Martin', '9780132350884', 2008, 3, 3
WHERE NOT EXISTS (SELECT 1 FROM books WHERE isbn = '9780132350884');

INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies)
SELECT gen_random_uuid(), 'The Pragmatic Programmer', 'David Thomas, Andrew Hunt', '9780135957059', 2019, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM books WHERE isbn = '9780135957059');

INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies)
SELECT gen_random_uuid(), 'Design Patterns', 'Gang of Four', '9780201633610', 1994, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM books WHERE isbn = '9780201633610');

INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies)
SELECT gen_random_uuid(), 'Domain-Driven Design', 'Eric Evans', '9780321125217', 2003, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM books WHERE isbn = '9780321125217');

INSERT INTO books (id, title, author, isbn, published_year, total_copies, available_copies)
SELECT gen_random_uuid(), 'Clean Architecture', 'Robert C. Martin', '9780134494166', 2017, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM books WHERE isbn = '9780134494166');
