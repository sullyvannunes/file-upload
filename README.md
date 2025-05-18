# File Upload/Download Service

A simple file upload and download service built with Go, demonstrating best practices for handling file operations in a web application.

## Features

- File upload with automatic MIME type detection
- File download with proper content type and disposition headers
- PostgreSQL database storage for file metadata and content
- Simple web interface for managing files
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or later
- PostgreSQL 15.4 or later
- Docker and Docker Compose (optional)

## Project Structure

```
.
├── app.go           # Main application logic and HTTP handlers
├── main.go          # Application entry point
├── db/             # Database related files
│   ├── schema.sql  # Database schema
│   └── queries/    # SQL queries
├── pg/             # Generated database code
├── web/            # Web interface components
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## Database Schema

The application uses a simple database schema with a single `files` table:

```sql
CREATE TABLE files (
    id bigint NOT NULL,
    file bytea NOT NULL,
    name character varying NOT NULL,
    mimetype character varying NOT NULL
);
```

## Getting Started

1. Start the PostgreSQL database:
   ```bash
   docker-compose up -d
   ```

2. Run the application:
   ```bash
   make run
   ```

The application will be available at `http://localhost:3030`

## API Endpoints

- `GET /cvs` - List all uploaded files
- `GET /cvs/new` - Show file upload form
- `POST /cvs` - Upload new file(s)
- `GET /cvs/{id}` - Download a specific file

## File Upload Features

- Supports multiple file uploads
- Automatic MIME type detection
- Files are stored in PostgreSQL as bytea

## Development

The project uses:
- `pgx` for PostgreSQL database access
- `sqlc` for type-safe SQL queries
- Standard library `net/http` for HTTP server
- `mime/multipart` for file upload handling

## License

MIT
