// Switch to bookhub database
db = db.getSiblingDB('bookhub');

// Create users collection with schema validation
// Field names match Go entity struct fields (lowercase): id, name, email, passwordhash, active, createdat, updatedat
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'email', 'passwordhash', 'active', 'createdat', 'updatedat'],
      properties: {
        id: {
          bsonType: 'binData',
          description: 'UUID stored as binary'
        },
        name: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        email: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        passwordhash: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        active: {
          bsonType: 'bool',
          description: 'must be a boolean and is required'
        },
        createdat: {
          bsonType: 'date',
          description: 'must be a date and is required'
        },
        updatedat: {
          bsonType: 'date',
          description: 'must be a date and is required'
        }
      }
    }
  }
});

// Create unique index on email
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ active: 1 });

// Insert default admin user (password: admin123) if not exists
const adminExists = db.users.findOne({ email: 'admin@bookhub.com' });
if (!adminExists) {
  db.users.insertOne({
    id: UUID(),
    name: 'Admin',
    email: 'admin@bookhub.com',
    passwordhash: '$2a$10$otJUHlZifNL133mJxahlJuDq7w5xv1S3RDeYVCHgikFZ0FOtov2f6',
    active: true,
    createdat: new Date(),
    updatedat: new Date()
  });
  print('Admin user created successfully');
} else {
  print('Admin user already exists');
}

// Create books collection with schema validation
// Field names match Go entity struct fields (lowercase): id, title, author, isbn, publishedyear, totalcopies, availablecopies, createdat, updatedat
db.createCollection('books', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['title', 'author', 'isbn', 'totalcopies', 'availablecopies', 'createdat', 'updatedat'],
      properties: {
        id: {
          bsonType: 'binData',
          description: 'UUID stored as binary'
        },
        title: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        author: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        isbn: {
          bsonType: 'string',
          description: 'must be a string and is required'
        },
        publishedyear: {
          bsonType: 'int',
          description: 'must be an integer'
        },
        totalcopies: {
          bsonType: 'int',
          description: 'must be an integer and is required'
        },
        availablecopies: {
          bsonType: 'int',
          description: 'must be an integer and is required'
        },
        createdat: {
          bsonType: 'date',
          description: 'must be a date and is required'
        },
        updatedat: {
          bsonType: 'date',
          description: 'must be a date and is required'
        }
      }
    }
  }
});

// Create indexes for books
db.books.createIndex({ isbn: 1 }, { unique: true });
db.books.createIndex({ title: 1 });
db.books.createIndex({ author: 1 });
db.books.createIndex({ availablecopies: 1 });

// Insert sample books (only if they don't exist)
const sampleBooks = [
  {
    title: 'Clean Code',
    author: 'Robert C. Martin',
    isbn: '9780132350884',
    publishedyear: 2008,
    totalcopies: 3,
    availablecopies: 3
  },
  {
    title: 'The Pragmatic Programmer',
    author: 'David Thomas, Andrew Hunt',
    isbn: '9780135957059',
    publishedyear: 2019,
    totalcopies: 2,
    availablecopies: 2
  },
  {
    title: 'Design Patterns',
    author: 'Gang of Four',
    isbn: '9780201633610',
    publishedyear: 1994,
    totalcopies: 2,
    availablecopies: 2
  },
  {
    title: 'Domain-Driven Design',
    author: 'Eric Evans',
    isbn: '9780321125217',
    publishedyear: 2003,
    totalcopies: 1,
    availablecopies: 1
  },
  {
    title: 'Clean Architecture',
    author: 'Robert C. Martin',
    isbn: '9780134494166',
    publishedyear: 2017,
    totalcopies: 2,
    availablecopies: 2
  }
];

sampleBooks.forEach(function(book) {
  const exists = db.books.findOne({ isbn: book.isbn });
  if (!exists) {
    db.books.insertOne({
      id: UUID(),
      title: book.title,
      author: book.author,
      isbn: book.isbn,
      publishedyear: NumberInt(book.publishedyear),
      totalcopies: NumberInt(book.totalcopies),
      availablecopies: NumberInt(book.availablecopies),
      createdat: new Date(),
      updatedat: new Date()
    });
    print('Book "' + book.title + '" created successfully');
  } else {
    print('Book "' + book.title + '" already exists');
  }
});

// Create loans collection with schema validation
// Field names match Go entity struct fields (lowercase): id, userid, bookid, borrowedat, duedate, returnedat, status
db.createCollection('loans', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['userid', 'bookid', 'borrowedat', 'duedate', 'status'],
      properties: {
        id: {
          bsonType: 'binData',
          description: 'UUID stored as binary'
        },
        userid: {
          bsonType: 'binData',
          description: 'UUID stored as binary and is required'
        },
        bookid: {
          bsonType: 'binData',
          description: 'UUID stored as binary and is required'
        },
        borrowedat: {
          bsonType: 'date',
          description: 'must be a date and is required'
        },
        duedate: {
          bsonType: 'date',
          description: 'must be a date and is required'
        },
        returnedat: {
          bsonType: ['date', 'null'],
          description: 'must be a date or null'
        },
        status: {
          enum: ['active', 'returned'],
          description: 'must be either active or returned'
        }
      }
    }
  }
});

// Create indexes for loans
db.loans.createIndex({ userid: 1 });
db.loans.createIndex({ bookid: 1 });
db.loans.createIndex({ status: 1 });
db.loans.createIndex({ duedate: 1 });
db.loans.createIndex({ userid: 1, bookid: 1, status: 1 });

print('Loans collection created successfully');
print('MongoDB initialization completed');
