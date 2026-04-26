# KeepSpace - Self-Hosted S3 Storage Gateway

A minimal, secure, self-hosted S3-compatible storage gateway. Manage your files with private Spaces and programmatic API keys.

![KeepSpace](https://storage.googleapis.com/dala-prod-public-storage/generated-images/7147aa4f-7147-437c-b5ae-29ea02e8cbd1/dashboard-preview-6e4c76fd-1777228188793.webp)

## 🚀 Features

- **Private Spaces**: Isolate your data in secure containers (no "buckets" terminology)
- **API Key Access**: Generate programmatic keys for each Space
- **S3 Compatible**: Built on MinIO for high performance
- **User Authentication**: JWT-based auth with access and refresh tokens
- **Modern UI**: Beautiful React dashboard with shadcn/ui
- **Self-Hosted**: Full control over your data

## 🏗️ Architecture

- **Frontend**: React + TypeScript + Vite + shadcn/ui
- **Backend**: Go + Gin framework
- **Database**: PostgreSQL 15
- **Storage**: MinIO (S3-compatible)
- **Auth**: JWT tokens + bcrypt password hashing
- **API Keys**: SHA-256 hashed for secure storage

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.26+
- Node.js 18+ (or Bun)
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd KeepSpace
```

### 2. Start Infrastructure

```bash
docker compose up -d
```

This starts:
- PostgreSQL on port 5432
- MinIO on ports 9000 (API) and 9001 (Console)

### 3. Start Backend

```bash
cd backend
go run .
```

The backend will:
- Connect to PostgreSQL and run migrations
- Connect to MinIO and create the "spaces" bucket
- Start API server on port 8080

### 4. Start Frontend

```bash
# In a new terminal, from project root
npm install
npm run dev
```

Frontend will be available at: http://localhost:3001

### 5. Access the Application

1. Open http://localhost:3001
2. Click "Sign Up" and create an account
3. Login with your credentials
4. Create your first Space
5. **Copy the API key** (shown only once!)
6. Use the API key to upload/download files

## 📚 Usage

### Web Dashboard

1. **Sign Up**: Create a new account
2. **Login**: Access your dashboard
3. **Create Space**: Click "Create Space" and give it a name
4. **Copy API Key**: Save the API key securely (shown only once)
5. **Manage Spaces**: View, delete spaces from the dashboard

### API Usage

#### Upload a File

```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "X-API-Key: ks_live_your_api_key_here" \
  -F "file=@myfile.pdf" \
  -F "path=/"
```

#### List Files

```bash
curl -X GET "http://localhost:8080/api/v1/files?path=/" \
  -H "X-API-Key: ks_live_your_api_key_here"
```

#### Download a File

```bash
curl -X GET "http://localhost:8080/api/v1/files/download?path=myfile.pdf" \
  -H "X-API-Key: ks_live_your_api_key_here" \
  -o myfile.pdf
```

#### Delete a File

```bash
curl -X DELETE "http://localhost:8080/api/v1/files?path=myfile.pdf" \
  -H "X-API-Key: ks_live_your_api_key_here"
```

## 🔧 Configuration

### Backend Environment Variables

Edit `backend/.env`:

```env
DATABASE_URL=postgres://keepspace:devpassword@localhost:5432/keepspace_dev?sslmode=disable
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### Frontend Environment Variables

Create `.env` in project root (optional):

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 📖 API Documentation

See [backend/README.md](backend/README.md) for complete API documentation.

### Quick Reference

**Authentication:**
- `POST /api/v1/auth/signup` - Create account
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token

**Spaces (JWT required):**
- `GET /api/v1/spaces` - List spaces
- `POST /api/v1/spaces` - Create space
- `DELETE /api/v1/spaces/:id` - Delete space

**Files (API Key required):**
- `POST /api/v1/files` - Upload file
- `GET /api/v1/files` - List files
- `GET /api/v1/files/download` - Download file
- `DELETE /api/v1/files` - Delete file
- `GET /api/v1/files/presigned-url` - Get presigned URL

## 🧪 Testing

### Complete Flow Test

```bash
# 1. Sign up
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 2. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

# 3. Create space
API_KEY=$(curl -s -X POST http://localhost:8080/api/v1/spaces \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test Space"}' \
  | grep -o '"api_key":"[^"]*' | cut -d'"' -f4)

echo "API Key: $API_KEY"

# 4. Upload file
echo "Hello KeepSpace!" > test.txt
curl -X POST http://localhost:8080/api/v1/files \
  -H "X-API-Key: $API_KEY" \
  -F "file=@test.txt"

# 5. List files
curl -X GET "http://localhost:8080/api/v1/files?path=/" \
  -H "X-API-Key: $API_KEY"

# 6. Download file
curl -X GET "http://localhost:8080/api/v1/files/download?path=test.txt" \
  -H "X-API-Key: $API_KEY"
```

## 🗂️ Project Structure

```
KeepSpace/
├── backend/
│   ├── auth/           # JWT and password hashing
│   ├── db/             # Database connection
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Auth middleware
│   ├── models/         # Database models
│   ├── storage/        # MinIO operations
│   ├── main.go         # Entry point
│   └── .env            # Configuration
├── src/
│   ├── components/     # React components
│   ├── pages/          # Page components
│   ├── lib/            # API client & utilities
│   └── main.tsx        # React entry point
├── docker-compose.yml  # Infrastructure
└── README.md           # This file
```

## 🔒 Security

- ✅ Passwords hashed with bcrypt
- ✅ JWT tokens with expiration
- ✅ API keys hashed with SHA-256
- ✅ CORS configured
- ✅ Input validation
- ✅ SQL injection protection (GORM)

### Production Checklist

- [ ] Change JWT_SECRET to a strong random value
- [ ] Enable HTTPS/TLS
- [ ] Use strong database passwords
- [ ] Add rate limiting
- [ ] Set up monitoring and logging
- [ ] Configure backups
- [ ] Review CORS settings
- [ ] Enable MinIO SSL

## 🐛 Troubleshooting

### Backend won't start
- Check if PostgreSQL is running: `docker compose ps`
- Check if port 8080 is available: `lsof -i :8080`
- Review backend logs for errors

### Frontend won't connect
- Ensure backend is running on port 8080
- Check browser console for CORS errors
- Verify API_URL in frontend configuration

### MinIO connection failed
- Access MinIO console: http://localhost:9001
- Login: minioadmin / minioadmin
- Verify bucket "spaces" exists

### Database errors
- Check PostgreSQL logs: `docker compose logs postgres`
- Verify DATABASE_URL in backend/.env
- Try restarting: `docker compose restart postgres`

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- UI components from [shadcn/ui](https://ui.shadcn.com/)
- Storage powered by [MinIO](https://min.io/)
- Database with [GORM](https://gorm.io/)

---

