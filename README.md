# Redis Rate Limiter

This project is a high-performance rate limiter written in Go using Redis and Lua scripting, implementing the Fixed Window Counter algorithm.

---

## 📦 Dependencies & Requirements

### 🔧 Runtime Dependencies

| Dependency | Version | Description                                                         |
| ---------- | ------- | ------------------------------------------------------------------- |
| Go         | ≥ 1.24  | Required for local development and testing                          |
| Redis      | ≥ 6.0   | Required for storing token bucket state and executing Lua scripts   |
| Docker     | ≥ 20.10 | (Optional) Used for containerized setup                             |
| Make       | any     | (Optional) For running helper tasks like tests and load simulations |

> 📌 Redis must be running at the address defined in `.env` (`REDIS_ADDR`, default: `localhost:6379`).

---

## 🛠 Features

- ✅ Fixed Window Counter algorithm implemented in Redis Lua Script (atomic)
- ✅ Supports custom limits per API key or IP address
- ✅ Burst traffic handling
- ✅ Simple HTTP server demonstration
- ✅ Unit tested (80%+ coverage)
- ✅ Dockerized for quick setup

---

## 🚀 How to Run

### Option 1: Using Docker Compose (Recommended)

```bash
# setup env, using docker setup
cp example.env .env

# Build and start the application with Redis
docker-compose up

# Stop the application
docker-compose down
```

### Option 2: Manual Setup

#### 1. Run Redis locally

```bash
docker run -p 6379:6379 redis
```

#### 2. Setup Env (use local)

```bash
cp example.env .env
```

#### 3. Run the application

```
go run main.go rest
```

#### 4. Test endpoint

```bash
curl -H "X-API-KEY: test-key" http://localhost:8080/ping
```

#### 5. Run Unit Testing

```bash
make test-coverage
```

#### 6. Run Script Burst request

```bash
make test-burst
```

---

## 📌 Design Overview

### ⏳ Algorithm: Fixed Window Counter with Ban (Redis Lua Script)

- Each client is allowed a maximum of `limit` requests per `window` seconds.
- If the number of requests exceeds the limit, the client is **banned** for `ban_duration` seconds.
- While banned, all incoming requests are immediately rejected.
- The counter resets after each window, effectively "refilling" the quota.

#### Returned Values:
- `{1, remaining}`: Request accepted. `remaining` indicates how many more requests are allowed in the current window.
- `{0, ban_duration}`: Request rejected due to rate limit. Ban has just been activated.
- `{-1, ttl}`: Request rejected because the client is currently banned. `ttl` is the remaining ban time.

### 🧠 Decisions

- Redis + Lua for atomic, distributed-safe operations.
- Use `X-API-KEY` header or fallback to IP address.
- TTL used to clean up idle keys.

### 🧪 Unit Testing

- 100 requests allowed (limit 100).
- 101th request blocked.
- After waiting >5s, token refilled, request allowed again.

---

## 📋 Assumptions & Limitations

- Redis is required and must be available at `localhost:6379` (can be configured).
- No persistent configuration store for client rules — defined in code.
- Rate limits are static; no dynamic runtime reconfiguration implemented.
