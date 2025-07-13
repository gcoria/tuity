# Tuity - Like twitter but tiny

## 5. [WIP Miro] Scalable Architecture

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
