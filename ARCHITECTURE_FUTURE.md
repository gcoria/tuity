## Scalable Architecture

## DB Selection

### Primary Database: PostgreSQL

- **ACID Compliance**: Strong consistency for follow operations and user data
- **JSON Support**: Flexible schema for tweets and user profiles
- **Performance**: Excellent read/write performance with proper indexing
- **Scalability**: Master-slave replication and sharding capabilities
- **JSON Support**: Native JSONB for tweet metadata
- **UUID Support**: Better for distributed systems
- **Proven at Scale**: Big tech companies used, tooling and community support
- **Sharding**: Can shard PostgreSQL when needed

### Secondary Database: Redis

- **Timeline Cache**: Users timelines
- **Session Storage**: User sessions and authentication tokens
- **Rate Limiting**: Token bucket counters
- **Clusters**: Horizontal scaling when needed

#### Scenario 1 Single Instance (0-100K users)

PostgreSQL Master + Redis Single Node

#### Scenario 2: Read Replicas (100K-1M users)

PostgreSQL Master → Read Replica 1
→ Read Replica 2
→ Read Replica 3
Redis Master → Redis Replica

#### Scenario 3: Sharding (1M-10M+ users)

Shard 1: Users 0-333333
Shard 2: Users 333334-666666  
Shard 3: Users 666667-999999
Each shard has its own master + replicas

### [WIP]More issues to think at scale

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
