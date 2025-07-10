# 🎟️ Tickitz Backend API

Tickitz is a backend RESTful API built for a movie ticket booking platform. It provides functionalities for users to explore movies, view showtimes, book tickets, and for admins to manage the system.

---

## 📌 Features

### 🔐 Authentication

- Register new users
- Secure login using JWT
- Role-based access (User / Admin)

### 🎬 Movies

- List of **popular** and **upcoming** movies
- Movie detail view
- Admin features: add, update, delete movies

### 🗓️ Schedules

- Filter movies by location, date, and cinema
- Showtimes management (Admin only)

### 🎟️ Ticket Booking

- Book tickets and select seats
- View booking history

### 📋 Extras

- Pagination and keyword search
- Admin dashboard endpoints

---

## 🧱 Tech Stack

- **Go (Golang)** using **Gin** framework
- **PostgreSQL** with **pgx** driver
- **JWT** for authentication and authorization
- **bcrypt** for password hashing
- **Redis** (for optional caching or session handling)
- **Swagger** for API documentation

---

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/ImA08/tickitz-backend.git
   cd tickitz-backend
   ```

Create a .env file and configure:

```# PostgreSQL Database Configuration

DBNAME=your_database_name
DBUSER=your_database_user
DBPASS=your_database_password
DBHOST=localhost
DBPORT=5432

# JWT Authentication

JWT_SECRET=your_jwt_secret
JWT_ISSUER=tickitz-api

# Google OAuth (if used)

GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Redis Configuration

RDS_HOST=localhost
RDS_PORT=6379```


2. Install dependencies
go mod tidy

3. Set up .env file

4. Run migrations (manually or with tool)
Example using psql:
psql -U your_user -d your_db -f migrations/init.sql

5. Run the app
go run cmd/main.go


🧪 API Endpoints (Summary)

| Method | Endpoint            | Description              |
| ------ | ------------------- | ------------------------ |
| POST   | `/auth/register`    | Register new user        |
| POST   | `/auth/login`       | User login               |
| GET    | `/movies`           | List movies with filters |
| GET    | `/movies/:id`       | Get movie details        |
| POST   | `/movies`           | Add new movie (Admin)    |
| PATCH  | `/movies/:id`       | Update movie (Admin)     |
| DELETE | `/movies/:id`       | Delete movie (Admin)     |
| GET    | `/schedules`        | Get movie schedules      |
| POST   | `/booking`          | Book tickets             |
| GET    | `/bookings/history` | View booking history     |


📂 Project Structure

tickitz-backend/
│
├── assets/ # Static or template assets
├── cmd/
│ └── main.go # Entry point
├── docs/ # Swagger documentation
│ ├── docs.go
│ ├── swagger.json
│ └── swagger.yaml
├── internal/
│ ├── handlers/ # HTTP handlers (controllers)
│ ├── middleware/ # Auth & role-check middleware
│ ├── models/ # Struct definitions for DB and domain models
│ ├── repositories/ # Database queries and logic
│ ├── routes/ # Router registration
│ └── utils/ # Helper utilities
├── migrations/ # SQL migration scripts
├── pkg/ # External packages (DB, JWT, Redis, etc.)
│ └── db.go, jwt.go, ...
├── public/
│ └── img/ # Uploaded movie images
├── tmp/ # Temp files
├── .dockerignore
├── .env # Environment variables
├── .fresh.yaml # For live reload (if using Fresh)
├── .gitignore
├── Dockerfile # Docker setup
├── go.mod
├── go.sum
└── README.md


👤 Author

Imanul Aufa
GitHub: @ImA08

📃 License

This project is licensed under the MIT License.

MIT License

Copyright (c) [2025] [Imanul Aufa]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell  
copies of the Software, and to permit persons to whom the Software is  
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in  
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR  
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,  
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE  
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER  
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN  
THE SOFTWARE.
