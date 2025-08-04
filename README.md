


##  Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone and start
git clone <repository>
cd secret-manager
docker-compose up --build
```

### Manual Setup

```bash
# Install dependencies
go mod tidy

# Set environment variables
cp .env.example .env
# Edit .env with your configuration

# Start PostgreSQL (or use Docker)
# Update DATABASE_URL in .env

# Run the application
go run .
```

## 🔧 Configuration

Create a `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/secretdb?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
ENCRYPTION_KEY=your-32-byte-encryption-key-here12
API_KEYS=api-key-1,api-key-2,api-key-3
```

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login (username: admin, password: password)

### Secrets
- `POST /api/v1/secrets` - Create secret
- `GET /api/v1/secrets` - List secrets (with filtering)
- `GET /api/v1/secrets/{id}` - Get secret
- `PUT /api/v1/secrets/{id}` - Update secret
- `DELETE /api/v1/secrets/{id}` - Delete secret
- `GET /api/v1/secrets/{id}/versions` - Get secret versions
- `GET /api/v1/secrets/{id}/versions/{version}` - Get specific version

### Audit
- `GET /api/v1/audit/logs` - Get audit logs

##  Authentication

### JWT Token
```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Use token in requests
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/secrets
```

### API Key
```bash
curl -H "X-API-Key: your-api-key" \
  http://localhost:8080/api/v1/secrets
```

##  Testing

```bash
# Run all tests
make test

# Run with coverage
go test ./tests/... -cover -v

# Run specific test
go test ./tests/ -run TestCreateSecret -v
```

##  Deployment

### Docker
```bash
make docker-build
make docker-run
```

### Production
1. Set strong encryption keys and JWT secrets
2. Use secure database credentials
3. Enable HTTPS/TLS termination
4. Set up proper logging and monitoring
5. Configure backup strategies


##  Security Features

- **AES-256-GCM Encryption**: All secret values encrypted at rest
- **JWT Authentication**: Secure token-based auth with expiration
- **API Key Support**: Simple API access for services
- **Audit Logging**: Complete activity tracking
- **Input Validation**: Comprehensive request validation
- **CORS Support**: Configurable cross-origin policies

##  Example Usage

```bash
# Create a secret
curl -X POST http://localhost:8080/api/v1/secrets \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "database-password",
    "value": "super-secret-password",
    "description": "Production database password",
    "tags": [{"name": "production"}, {"name": "database"}]
  }'

# Get the secret
curl -H "X-API-Key: your-api-key" \
  http://localhost:8080/api/v1/secrets/{secret-id}

# Update secret (creates new version)
curl -X PUT http://localhost:8080/api/v1/secrets/{secret-id} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"value": "new-password"}'
```

The system automatically handles encryption, versioning, and audit logging for all operations.
