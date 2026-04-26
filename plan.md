
# KeepSpace MVP Implementation Plan

## Project Goal:
Build a self-hosted minimal version of S3Drive, renamed to KeepSpace, that allows users to:
1. Sign up and log in.
2. Create secure "Spaces" (private storage areas).
3. Obtain an API key for each Space.
4. Use the API key to upload, list, and download files.
All files will be stored on an S3-compatible object store (MinIO locally).
A React developer portal will display Spaces and API keys.

## Naming Convention:
- User-facing concept: "Space" (private, isolated storage container).
- Internal S3 bucket name: "spaces" (only on the server-side, never exposed to users).
- Never use the word "bucket" in code, API docs, or UI.

## Tech Stack:
- Backend: Go (Gin framework)
- Database: PostgreSQL
- Object Storage: MinIO (S3-compatible)
- Frontend: React (Vite)
- Authentication: JWT (access/refresh tokens)
- API Key Authentication: SHA-256 hash for lookup, storing a 64-char hex string.

## Step-by-step Tasks:

### 1. Infrastructure (Docker Compose):
- Create `docker-compose.yml` with services:
    - `postgres:15-alpine` (credentials: `keepspace/devpassword`, db: `keepspace_dev`, port: `5432`)
    - `minio/minio` (root: `minioadmin`, password: `minioadmin`, ports: `9000` (S3 API), `9001` (console))
- Start with `docker compose up -d`.
- In MinIO console (http://localhost:9001), manually create an S3 bucket named `spaces`.

### 2. Database Schema:
- Migration file for PostgreSQL:
    - `users` table: `id` (UUID PK), `email` (VARCHAR UNIQUE), `password_hash` (VARCHAR), `created_at` (TIMESTAMPTZ).
    - `spaces` table: `id` (UUID PK), `owner_id` (UUID FK to users.id), `name` (VARCHAR), `api_key_hash` (VARCHAR, SHA-256 hash of API key), `created_at` (TIMESTAMPTZ).

### 3. Backend - Authentication (Go/Gin):
- Sign-up endpoint `POST /api/v1/auth/signup`:
    - Accepts `{email, password}`.
    - Hashes password (bcrypt), inserts into `users`, returns 201.
- Login endpoint `POST /api/v1/auth/login`:
    - Verifies email/password.
    - Returns JWT access token (15min) and refresh token (7 days).
- Refresh endpoint `POST /api/v1/auth/refresh`:
    - Handles token refresh.
- Middleware for JWT validation on protected routes.

### 4. Backend - Space Management (Go/Gin):
- All routes require valid JWT.
- Create Space `POST /api/v1/spaces`:
    - Request body: `{name}`.
    - Generate a random 64-char hex API key.
    - Store SHA-256 hash of the API key in `api_key_hash`.
    - Insert into `spaces` table.
    - Return plain API key `{id, name, api_key, created_at}` (user must copy it).
- List Spaces `GET /api/v1/spaces`:
    - Return all spaces for the logged-in user (id, name, created_at).
- Delete Space `DELETE /api/v1/spaces/:id`:
    - Remove related S3 objects with prefix `spaces/{space_id}/`.

### 5. Backend - File Operations (Go/Gin - Authenticated by API Key):
- Endpoints protected by `X-API-Key` header.
- S3 key pattern: `spaces/{space_id}/{file_path}`.
- Upload file `POST /api/v1/files`:
    - Header: `X-API-Key: <plain-text key>`.
    - Body: `multipart/form-data` with `file` and optional `path`.
    - Validate API key: Hash incoming key with SHA-256, find Space by hash, verify ownership.
    - Upload to MinIO: `PutObject` with key `spaces/{space_id}/{path}/{filename}`.
- List files `GET /api/v1/files?path=/`:
    - Use `X-API-Key` header.
    - List objects in MinIO with prefix `spaces/{space_id}/{path}`.
    - Return list of folders (prefixes) and files (objects) with name, size, last modified.
- Download file `GET /api/v1/files/:filename?path=/`:
    - Use `X-API-Key` header.
    - Generate a pre-signed URL from MinIO or stream directly.

### 6. Backend - CORS & Config:
- Enable CORS for `http://localhost:5173` (React dev server).
- Handle configuration via environment variables.

### 7. Frontend - Developer Portal (React/Vite):
- Routes: `/login`, `/signup`, `/spaces`, `/spaces/new`.
- Use React Router for navigation.
- Store JWT in localStorage (MVP).
- Login/Signup forms.
- Display list of user's Spaces.
- Form to create a new Space, display and copy API key.
- Fetch `/api/v1/spaces` using JWT in `Authorization` header.
- Display Space name, creation date.
- Button to delete Space.

### 8. Testing:
- Run Go backend (`go run .`).
- Run React frontend (`npm run dev`).
- Test full flow: Sign up, create Space, copy key.
- Use `curl` or Postman for file operations: upload, list, download using API key.

### 9. Code Refactoring:
- Ensure all code, API docs, and UI consistently use "Space" instead of "bucket".

