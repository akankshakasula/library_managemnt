Library Management System Backend
This is a backend API for a basic Library Management System, built with Go and PostgreSQL. It handles user authentication, book inventory, and borrowing/returning processes.

✨ Key Features
User Management: Sign up, log in, and role-based access (Librarian, Student, General).

Book Management: Add new books, record donations, and list all books.

Borrowing & Returns: Borrow books, return them, and calculate overdue fines.

Security: Uses JWT for API authentication and bcrypt for password hashing.

🚀 Technologies
Go (Golang)

Fiber v2 (Web Framework)

PostgreSQL (Database)

GORM (ORM)

Golang-JWT v5 (JWT)

Bcrypt (Password Hashing)

⚙️ How to Run Locally
Prerequisites: Install Go and have a running PostgreSQL database.

Database: Create a database (e.g., library) and a user (e.g., library_user) for it.

.env File: Create a .env file in the project root with your database connection string and a JWT secret:

DATABASE_URL="host=localhost port=5432 user=library_user password=your_db_password dbname=library sslmode=disable TimeZone=Asia/Kolkata"
JWT_SECRET="your_strong_jwt_secret_key"

Install Dependencies:

go mod tidy

Run Server:

go run ./cmd/server

The API will be available at http://127.0.0.1:3000.

🔌 API Endpoints (Examples)
POST /api/signup - Register a new user.

POST /api/signin - Login and get a JWT token.

GET /api/books - Get all books (requires JWT).

POST /api/books - Add a book (requires librarian JWT).

POST /api/books/borrow - Borrow a book (requires JWT).

POST /api/books/return/:id - Return a book (requires JWT).
