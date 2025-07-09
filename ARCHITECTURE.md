# Tuity - Like twitter but tiny

## Core Requirements Analysis

- **Tweet Management**: Users can post messages (≤280 characters)
- **Follow System**: Users can follow/unfollow other users
- **Timeline Generation**: Users view chronological feeds of followed users tweets
- **Scalability**: Handle millions of users
- **Performance**: Optimized for read operations
- **Consistency**: Eventual consistency for reads

## Hexagonal Architecture

- **Technology Independence**: Core business logic isolated from external concerns
- **Testability**: Easy unit testing of business logic
- **Flexibility**: Easy to swap adapters without affecting core
- **Maintainability**: Clear separation of concerns

###

## MVP Architecture

### Core Layer

#### Domain

- **User**: Represents a platform user with profile information
- **Tweet**: Represents a message with content, timestamp, and metadata
- **Follow**: Represents follow relationships between users
- **Timeline**: Represents user's feed of tweets

#### Application Layer

- **TweetService**: Create, retrieve, and manage tweets
- **FollowService**: Handle follow/unfollow operations
- **TimelineService**: Generate and manage user timelines
- **UserService**: Basic user operations

#### Domain Services

- **FanoutService**: Timeline generation and distribution logic
- **ContentValidator**: Tweet content validation
- **NotificationService**: User notifications

### Ports Layer

#### Input Ports

```go
type TweetPort interface {}

type FollowPort interface {}

type TimelinePort interface {}
```

#### Output Ports

```go
type TweetRepositoryPort interface {}

type FollowRepositoryPort interface {}

type CachePort interface {}

type EventPort interface {}
```

## Adapters Layer

- **HTTP Adapter**: Gin framework for REST API
- **Database Adapter**: GORM for PostgreSQL
- **Cache Adapter**: KVS in memory
- **Event Adapter**: event bus in memory

## 5. Scalable Architecture (Wip Miro)

### Requirements / Issues to think

- **Timeline Generation**: For users with many followed, Message broker
- **Tweets**: 100:1 Read/Write ratio
- **Cache**: Needs a distributed cache, handle consistency
- **Timeline Updates**: Fan-out to followers, Redis(cluster?),
- **DB**: Sharding + eventual consistency, celebrity accounts
- **CDN** required ?
- **Gateway**: rate limit, Authentication
- **Container Orchestration**: Kubernetes
- **Load Balancing**: NGINX/HAProxy

- **Logging**:
- **Monitoring**:
  - **Business Metrics**: DAU, tweets/day, engagement rates
  - **Technical Metrics**: Response times, error rates, throughput
  - **Infrastructure Metrics**: CPU, memory, disk, network
  - **SLA Violations**: Response time > 500ms
  - **Error Rates**: > 1% error rate
  - **Resource Usage**: > 80% CPU/memory
