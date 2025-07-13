# Tuity - like twitter but tiny

## 🚀 Features

- **Tweet Management**: Create, read, delete tweets (280 character limit)
- **Follow System**: Follow/unfollow users with asymmetric relationships
- **Timeline Generation**: Personalized timeline with cached results
- **Rate Limiting**: Token bucket algorithm implementation
- **Hexagonal Architecture**
- **Docker Support**

## 🏗️ Architecture

This project follows hexagonal architecture with:

- **Core Domain**: Pure business logic (User, Tweet, Follow, Timeline)
- **Application Services**: Use cases and business workflows
- **Ports**: Interfaces for external interactions
- **Adapters**: External integrations (HTTP, Database, Cache)

## 🚦 Getting Started

### Prerequisites

- Go 1.23+
- Docker (optional)

### Quick Start

```bash
# Clone and run
git clone [<repository-url>](https://github.com/gcoria/tuity)
cd tuity
go mod download
make run
```

The API will be available at `http://localhost:8080`

### Environment Variables

```bash
PORT=8080
GIN_MODE=release
TIMELINE_DEFAULT_LIMIT=20
TIMELINE_MAX_LIMIT=100
RATE_LIMIT_TWEET_CREATE=10
RATE_LIMIT_FOLLOW_OPS=20
RATE_LIMIT_TIMELINE_REQUESTS=100
CACHE_TIMELINE_TTL=5m
```

## 🐳 Docker Deployment

```bash
# Build and run
make docker-build
make docker-run

# Development mode
make docker-dev
```

## 📡 API Documentation

### Authentication

Include `X-User-ID` header for protected endpoints.

### Key Endpoints

#### Users

```http
POST /api/v1/users                    # Create user
GET  /api/v1/users/{id}               # Get user by ID
GET  /api/v1/users/username/{username} # Get user by username
```

#### Tweets

```http
POST   /api/v1/tweets          # Create tweet (requires X-User-ID)
GET    /api/v1/tweets/{id}     # Get tweet
DELETE /api/v1/tweets/{id}     # Delete tweet (requires X-User-ID)
GET    /api/v1/users/{id}/tweets # Get user tweets
```

#### Follow System

```http
POST   /api/v1/users/{id}/follow        # Follow user (requires X-User-ID)
DELETE /api/v1/users/{id}/follow        # Unfollow user (requires X-User-ID)
GET    /api/v1/users/{id}/following     # Get following list
GET    /api/v1/users/{id}/followers     # Get followers list
```

#### Timeline

```http
GET  /api/v1/users/{id}/timeline?limit=20  # Get timeline
POST /api/v1/users/{id}/timeline/refresh   # Refresh timeline
```

### Example Usage

```bash
# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "display_name": "Alice Smith"}'

# Create a tweet
curl -X POST http://localhost:8080/api/v1/tweets \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-id-here" \
  -d '{"content": "Hello world!"}'

# Get timeline
curl http://localhost:8080/api/v1/users/:user_id/timeline?limit=10
```

## 🧪 Testing

```bash
make test              # Run tests
make test-coverage     # Run with coverage
make lint             # Run linter
make format           # Format code
```

## 🔧 Development

```bash
make setup    # Setup dev environment
make dev      # Start with hot reload
make help     # Show all commands
```

## 🚀 Performance Optimizations

### Timeline Generation


### Rate Limiting
- **Token Bucket**: Per-user token buckets with configurable limits

### Memory Repositories
- **Dual Indexing**: O(1) lookups for ID and username
- **Concurrent Safe**: Read-write mutexes

## 🔮 Future Scalability[WIP]

- **Database**: PostgreSQL with proper indexing, sharding
- **Cache**: Redis for distributed caching
- **Message Queue**: RabbitMQ for async processing
- **Load Balancing**: Multiple service instances
