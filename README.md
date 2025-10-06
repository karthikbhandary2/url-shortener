# URL Shortener

A high-performance URL shortener service built with Go, Fiber framework, and Redis. This service provides fast URL shortening with rate limiting and analytics capabilities.

## 🚀 Features

- **Fast URL Shortening**: Generate short URLs with custom or auto-generated IDs
- **Rate Limiting**: API quota system to prevent abuse
- **Redis Caching**: High-performance data storage and retrieval
- **URL Validation**: Comprehensive URL validation and sanitization
- **Analytics Ready**: Track URL usage and statistics
- **Docker Support**: Containerized deployment with Docker Compose
- **RESTful API**: Clean and intuitive API endpoints

## 📁 Project Structure

```
url-shortener/
├── api/                          # Main API application
│   ├── database/                 # Database connection and utilities
│   │   └── database.go          # Redis connection setup
│   ├── helpers/                  # Utility functions
│   │   └── helpers.go           # URL validation and helper functions
│   ├── routes/                   # API route handlers
│   │   ├── shorten.go           # URL shortening endpoint
│   │   └── resolve.go           # URL resolution endpoint
│   ├── .env                     # Environment variables
│   ├── Dockerfile               # API container configuration
│   ├── go.mod                   # Go module dependencies
│   ├── go.sum                   # Go module checksums
│   └── main.go                  # Application entry point
├── db/                          # Database configuration
│   └── Dockerfile               # Redis container configuration
├── .data/                       # Redis data persistence (gitignored)
├── .gitignore                   # Git ignore rules
├── docker-compose.yml           # Multi-container orchestration
└── README.md                    # Project documentation
```

## 🛠️ Technology Stack

- **Backend**: Go 1.24+
- **Web Framework**: Fiber v2 (Express-inspired web framework)
- **Database**: Redis (In-memory data structure store)
- **Containerization**: Docker & Docker Compose
- **Dependencies**:
  - `github.com/gofiber/fiber/v2` - Web framework
  - `github.com/go-redis/redis/v8` - Redis client
  - `github.com/google/uuid` - UUID generation
  - `github.com/asaskevich/govalidator` - URL validation
  - `github.com/joho/godotenv` - Environment variable loading

## 🚦 API Endpoints

### Shorten URL
```http
POST /api/v1
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "short": "custom-id",  // Optional: custom short ID
  "expiry": 24           // Optional: expiry in hours
}
```

**Response:**
```json
{
  "url": "https://example.com/very/long/url",
  "short": "localhost:3000/abc123",
  "expiry": 24,
  "rate_limit": 9,
  "rate_limit_reset": 1696608000
}
```

### Resolve URL
```http
GET /:shortId
```

**Response:** HTTP 301 redirect to original URL

## ⚙️ Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_ADD` | Redis server address | `db:6379` |
| `DB_PASS` | Redis password | `""` (empty) |
| `APP_PORT` | Application port | `:3000` |
| `DOMAIN` | Base domain for short URLs | `localhost:3000` |
| `API_QUOTA` | Rate limit per IP | `10` |

## 🐳 Quick Start with Docker

### Prerequisites
- Docker
- Docker Compose

### 1. Clone the Repository
```bash
git clone https://github.com/karthikbhandary2/url-shortener.git
cd url-shortener
```

### 2. Start the Services
```bash
docker-compose up -d
```

This will start:
- **API Service**: Available at `http://localhost:3000`
- **Redis Database**: Available at `localhost:6379`

### 3. Test the API
```bash
# Shorten a URL
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/karthikbhandary2/url-shortener"}'

# Access the shortened URL
curl -L http://localhost:3000/{short-id}
```

## 🔧 Local Development Setup

### Prerequisites
- Go 1.24+
- Redis server

### 1. Install Dependencies
```bash
cd api
go mod download
```

### 2. Start Redis
```bash
# Using Docker
docker run -d -p 6379:6379 redis:alpine

# Or install locally (Ubuntu/Debian)
sudo apt-get install redis-server
redis-server
```

### 3. Configure Environment
```bash
cp api/.env.example api/.env
# Edit api/.env with your configuration
```

### 4. Run the Application
```bash
cd api
go run main.go
```

## 📊 Rate Limiting

The service implements IP-based rate limiting:
- Default: 10 requests per IP
- Configurable via `API_QUOTA` environment variable
- Rate limit resets every hour
- Returns current limit and reset time in response headers

## 🔒 URL Validation

The service validates URLs using multiple checks:
- Protocol validation (http/https)
- Domain validation
- Malicious URL detection
- Custom domain restrictions (configurable)

## 📈 Monitoring and Logging

- **Request Logging**: All requests are logged with timestamps
- **Error Handling**: Comprehensive error responses
- **Health Checks**: Built-in health check endpoints
- **Metrics**: Redis connection and performance metrics

## 🚀 Production Deployment

### Docker Compose (Recommended)
```bash
# Production environment
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Deployment
1. Build the Go binary:
   ```bash
   cd api
   CGO_ENABLED=0 GOOS=linux go build -o url-shortener main.go
   ```

2. Set up Redis with persistence
3. Configure reverse proxy (Nginx/Apache)
4. Set up SSL certificates
5. Configure monitoring and logging

### Environment Configuration
```bash
# Production .env
DB_ADD=redis-server:6379
DB_PASS=your-secure-password
APP_PORT=:3000
DOMAIN=yourdomain.com
API_QUOTA=100
```

## 🧪 Testing

### Unit Tests
```bash
cd api
go test ./...
```

### Integration Tests
```bash
# Start services
docker-compose up -d

# Run tests
go test -tags=integration ./tests/
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 -H "Content-Type: application/json" \
   -p test-data.json http://localhost:3000/api/v1
```

## 🔧 Configuration Options

### Redis Configuration
- **Persistence**: Configured for data durability
- **Memory Policy**: LRU eviction for optimal performance
- **Connection Pooling**: Optimized for concurrent requests

### Application Configuration
- **CORS**: Configurable cross-origin resource sharing
- **Timeouts**: Request and database timeout settings
- **Logging**: Structured logging with different levels

## 📝 API Response Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | URL shortened successfully |
| 301 | Redirect to original URL |
| 400 | Bad request (invalid URL/parameters) |
| 429 | Rate limit exceeded |
| 500 | Internal server error |
| 404 | Short URL not found |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

### Common Issues

**Redis Connection Failed**
```bash
# Check Redis status
docker-compose logs db

# Restart Redis
docker-compose restart db
```

**Port Already in Use**
```bash
# Find process using port 3000
lsof -i :3000

# Kill the process
kill -9 <PID>
```

**Go Module Issues**
```bash
# Clean module cache
go clean -modcache
go mod download
```

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
go run main.go
```
