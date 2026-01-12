# Token Bucket Rate Limiter

A simple token bucket rate limiter implementation, with per-user tracking and burst handling, backed by Redis for distributed state management.

## Table of Contents

- [Key Features](#key-features)

- [Architecture Overview](#architecture-overview)

- [Tech Stack](#tech-stack)

- [Getting Started](#getting-started)

  - [Prerequisites](#prerequisites)

  - [Installation](#installation)

- [Configuration](#configuration)

- [Usage](#usage)

- [Project Structure](#project-structure)

- [Roadmap](#roadmap)

- [Contributing](#contributing)

- [Testing](#testing)

- [License](#license)

- [Acknowledgements](#acknowledgements)

## Key Features

*   **Token Bucket Algorithm**: Implements the classic token bucket algorithm for efficient rate limiting.

*   **Per-User Tracking**: Tracks rate limits for individual users or clients using a unique identifier (e.g., `X-Api-Key`).

*   **Burst Handling**: Allows for short bursts of requests up to a defined capacity.

*   **Redis-Backed State**: Leverages Redis for storing and managing token bucket states, enabling distributed rate limiting across multiple instances.

*   **Atomic Operations with Lua**: Utilizes Redis Lua scripting to ensure atomic updates to token buckets, preventing race conditions.

*   **Go HTTP Middleware**: Provides a simple HTTP middleware for easy integration into Go web applications.

## Architecture Overview

This project implements a distributed token bucket rate limiter. At its core, the `RateLimiter` struct manages the rate limiting logic, backed by a Redis instance. When a request comes in, the `rateLimitMiddleware` extracts a user identifier (e.g., from an `X-Api-Key` header). This identifier is then used as a key to interact with Redis.

The actual token bucket logic—checking if a request is allowed, refilling tokens, and decrementing tokens—is executed atomically within Redis using a Lua script. This ensures consistency and prevents race conditions in a concurrent environment. The Go application acts as a thin client, calling the Redis Lua script and interpreting its result to either allow the request to proceed or return a `429 Too Many Requests` status.

## Getting Started

Follow these instructions to set up and run the rate limiter locally.

### Prerequisites

Before you begin, ensure you have the following installed:

*   **Go**: Version 1.20 or higher.

    *   [Download Go](https://golang.org/dl/)

*   **Redis**: A running Redis server instance.

    *   [Install Redis](https://redis.io/docs/getting-started/installation/)

### Installation

1.  **Clone the repository**:

```bash
git clone https://github.com/sakshamg567/ratelimit.git

cd ratelimit

```
2.  **Download Go modules**:

```bash
go mod tidy

```
## Configuration

The rate limiter's behavior is configured via the `Config` struct in `limiter.go`. For a production environment, it's recommended to externalize these values, for example, using environment variables.

| ENV | Description | Example |
|---|---|---|
| `REDIS_ADDR` | Address of the Redis server. | `localhost:6379` |
| `CAPACITY` | Maximum number of tokens (burst capacity) for a bucket. | `5` |
| `REFILL_RATE` | Number of tokens added to the bucket per refill interval. | `1` |
| `REFILL_INTERVAL_MS` | Time in milliseconds between token refills. | `1000` (for 1 second) |

* Note: The current `main.go` hardcodes these values. For dynamic configuration, you would modify `main.go` to read from environment variables or a configuration file.*

## Usage

To run the example HTTP server with the rate limiter:

1.  **Start the Redis server** if it's not already running.

2.  **Run the Go application**:

```bash
go run main.go limiter.go bucket.go

```
The server will start on `http://localhost:8080`.

3.  **Test the endpoints**:

    *   **Unprotected endpoint**: This endpoint is not rate-limited.

```bash
curl http://localhost:8080/api/unprotected

# Expected: hello world

```
*   **Protected endpoint**: This endpoint is rate-limited. It requires an `X-Api-Key` header.

The default configuration allows 5 requests per second per API key.

        First 5 requests (should succeed):

```bash
curl -H "X-Api-Key: user123" http://localhost:8080/api/protected

# Expected: hello world

```
Subsequent requests within the same second (should be rate-limited):

```bash
curl -H "X-Api-Key: user123" http://localhost:8080/api/protected

# Expected: too many requests (HTTP 429)

```
*   **Quota endpoint**: Check the current token count for an API key.

```bash
curl -H "X-Api-Key: user123" http://localhost:8080/api/quota

# Expected: {"tokens":"X"} (where X is the remaining tokens)

```
## Project Structure

```
.

├── bucket.go
├── go.mod

├── go.sum
├── leaky_test.go

├── limit.lua
├── limiter.go

└── main.go

```
## Testing

To run the tests for the project:

```bash
go test ./...

```
The `leaky_test.go` file contains an example integration test that simulates multiple requests to the protected endpoint to verify the rate-limiting behavior, including burst handling and the `429 Too Many Requests` response.

