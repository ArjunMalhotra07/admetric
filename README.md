# AdMetric - Ad Click Analytics Service

AdMetric is a service for tracking and analyzing ad clicks. It provides APIs for recording ad clicks and retrieving analytics data.

## Features

- Record ad clicks with playback time
- Track total clicks per ad
- Get analytics data for different time frames (minutes, hours, days)
- Kafka integration for reliable message processing
- Circuit breaker pattern for fault tolerance
- Batch processing for efficient database operations

## API Endpoints

### 1. Get All Ads

- **URL**: `localhost:8888/ads`
- **Method**: `GET`
- **Description**: Retrieves all ads from the database
- **Response**: Array of Ad objects
  ```json
  [
        {
        "id": "1",
        "image_url": "https://example.com/images/ad1.jpg",
        "target_url": "https://example.com/landing/ad1",
        "created_at": "2025-04-10T16:25:10.223Z",
        "updated_at": "2025-04-10T16:25:10.223Z",
        "deleted_at": null,
        "Clicks": [
            {
                "id": "000b25f2-ca19-4d7f-ae0a-96ec24356635",
                "ad_id": "1",
                "Ad": {
                    "id": "",
                    "image_url": "",
                    "target_url": "",
                    "created_at": "0001-01-01T00:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "deleted_at": null,
                    "Clicks": null,
                    "total_clicks": 0
                },
                "ip": "127.0.0.1",
                "playback_time": 258,
                "timestamp": "2025-04-10T16:59:26.239Z"
            },{...}
        ],
        "total_clicks": 121
      }
  ]
  ```

### 2. Record Click

- **URL**: `localhost:8888/ads/click`
- **Method**: `POST`
- **Description**: Records a click event for an ad
- **Request Body**:
  ```json
   {
   "ad_id": "2",
   "playback_time": 120,
   "timestamp": "2025-04-10T14:30:00Z"
   }
  ```
- **Response**:
  ```json
  {
    "message": "Click recorded"
  }
  ```
- **Status Code**: 202 Accepted

### 3. Get Click Count

- **URL**: `localhost:8888/ads/:id/clicks`
- **Method**: `GET`
- **Description**: Gets the total number of clicks for a specific ad
- **URL Parameters**:
  - `id`: Ad ID
- **Response**:
  ```json
   {
      "ad_id": "2",
      "total_clicks": 1
   }
  ```

### 4. Get Click Analytics

- **URL**: `localhost:8888/ads/:id/analytics`
- **Method**: `GET`
- **Description**: Gets click analytics for a specific ad within a time frame
- **URL Parameters**:
  - `id`: Ad ID
- **Query Parameters**:
  - `timeframe`: Time frame for analytics (default: "1h")
    - Options: "1m", "5m", "15m", "30m", "1h", "6h", "12h", "1d", "7d", "30d"
- **Response**:
  ```json
   {
      "ad_id": "5",
      "clicks": 62,
      "timeframe": "7d"
   }
  ```

## Running the Application

You have two options to run the application:

### Option 1: Using Docker Compose with App Image

1. Make sure Docker and Docker Compose are installed
2. Run the following command:
   ```
   make compose
   ```
3. The application will be available at http://localhost:8888

Note: This option might take 5-10 seconds to completely start the application since mysql container might take some time to come to healthy state.

### Option 2: Running Locally without App image (Stack.yml file)

1. Run the following command:
   ```
   make stack
   ```
2. Set up the environment variables:
   ```
   source export.sh
   ```
3. Run the application:
   ```
   make run
   ```
4. The application will be available at http://localhost:8888

## Click Simulator Client

To test the analytics API which fetches counts for a specific timeframe, you need at least 100 clicks recorded. A client programn has been added that simulates multiple clicks:

1. Navigate to the client directory:
   ```bash
   cd client
   ```
2. Run the client:
   ```bash
   go run click_simulator.go
   ```

The client will:

- Send 450 clicks in total
- Send them in batches of 50
- Pause for 1 minute between batches
- Randomly select ad IDs between 1-10
- Provide progress updates in the console

This will quickly populate your database with enough data to test the analytics API.

## Architecture

The application uses:

- Fiber for the HTTP server
- GORM for database operations
- Kafka for message processing
- Circuit breaker pattern for fault tolerance
- Batch processing for efficient database operations

## Environment Variables

- `BASE_URL`: Base URL for the application
- `HTTP_HOST`: HTTP host
- `HTTP_PORT`: HTTP port
- `LOG_FILE`: Log file path
- `KAFKA_BROKER`: Kafka broker address
- `MYSQL_USER`: MySQL username
- `MYSQL_PASSWORD`: MySQL password
- `MYSQL_HOST`: MySQL host
- `MYSQL_PORT`: MySQL port
- `MYSQL_DB`: MySQL database name
- `MYSQL_ROOT_PASSWORD`: MySQL root password
- `MYSQL_DATA`: MySQL data directory

## Workflow

- Additionally a github workflow has been added in the repository which builds the latest image of the application and pushes to the Docker Hub whenever the main branch is pushed to.

## License

MIT
