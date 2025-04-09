# AdMetric - Video Advertisement Tracking System

A fault-tolerant Go backend for managing and tracking video advertisements. The system provides real-time analytics, handles high traffic efficiently, and remains resilient under partial failures.

## Features

- **Ad Management**: Store and retrieve ad metadata
- **Click Tracking**: Record and analyze user clicks with playback time
- **Real-time Analytics**: Get performance metrics for ads over various timeframes
- **Fault Tolerance**: Circuit breaker pattern to handle partial failures
- **Data Resilience**: Local backup for click events during system issues
- **Scalability**: Efficient batch processing and in-memory caching

## Architecture

The system uses a multi-layered architecture:

1. **HTTP Layer**: Fiber-based REST API
2. **Service Layer**: Business logic with circuit breakers
3. **Repository Layer**: Data access with GORM
4. **Storage Layer**: MySQL for persistent storage, Redis for caching

## API Endpoints

### GET /ads

Returns a list of ads with basic metadata.

### POST /ads/click

Records a click event with details (Ad ID, timestamp, IP, video playback time).

### GET /ads/analytics

Returns analytics data for ads within a specified timeframe.

## Setup and Installation

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (for local development)

### Running with Docker Compose

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/admetric.git
   cd admetric
   ```

2. Start the services:

   ```
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

### Local Development

1. Install dependencies:

   ```
   go mod download
   ```

2. Set up environment variables (copy from .env.example):

   ```
   cp .env.example .env
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

## Configuration

The application can be configured using environment variables:

- `DB_HOST`: MySQL host
- `DB_PORT`: MySQL port
- `DB_USER`: MySQL user
- `DB_PASSWORD`: MySQL password
- `DB_NAME`: MySQL database name
- `REDIS_URL`: Redis connection URL

## Fault Tolerance Features

1. **Circuit Breaker**: Prevents cascading failures when downstream services are unavailable
2. **Local Backup**: Saves click events to disk when the database is unavailable
3. **Batch Processing**: Efficiently processes click events in batches
4. **In-memory Caching**: Keeps frequently accessed data in memory
5. **Periodic Pruning**: Removes inactive ads from memory to prevent memory leaks

## Analytics Implementation

The system provides analytics at different timeframes:

- 1 hour
- 6 hours
- 24 hours
- 7 days
- 30 days

Analytics include:

- Total clicks
- Unique visitors
- Average playback time
- Click-through rate (CTR)

## License

MIT
