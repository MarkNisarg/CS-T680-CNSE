# Voter API

The Voter API is a RESTful API that allows users to manage voters and their voting history. It utilizes the Gin framework for handling HTTP requests and Redis for caching voter data. It provides endpoints to add, retrieve, update, and delete voter records, as well as manage their voting history. Additionally, the API offers a health check endpoint to monitor the API's status.

## Prerequisites

Make sure you have Go Programming Language installed on your machine. If not, you can download and install it from the official website:
- Go: https://golang.org/
- Docker: https://www.docker.com/get-started

## Installation

To clone the repository, use the following command:

```bash
git clone https://github.com/MarkNisarg/CS-T680-CNSE.git
cd voter-container
```

## Usage

## Running with Docker

### Method 1. Build and Run Using Single Docker Command:

Build and run the docker image:

```bash
docker compose up --build
```

-----

### Method 2. Build and Run Using Single Make Command:

Build and run the docker image:

```bash
make docker-start
```

-----

### Method 3. Build and Run Using Docker Commands:

Build the docker image:

```bash
docker compose build
```

Run the docker container:

```bash
docker compose up
```

-----

### Method 4. Build and Run Using Make Commands:
Build the Docker image:

```bash
make docker-build
```

Run the Docker container:

```bash
make docker-run
```

-----

### Method 5. Using Published Docker Image

If you prefer not to build the Docker image or have access to the source code, you can use the published Docker image directly. The image is available at [`nisargrajendrakumar/voter-api`](https://hub.docker.com/r/nisargrajendrakumar/voter-api).

The API will be accessible at http://localhost:1080  
The Redis stack will be accessible at http://localhost:8001

## Configuring Redis

The Voter API is designed to work seamlessly with Redis, which is set up to automatically start in its own Docker container before the Voter API starts up (as defined in `docker-compose.yml`). The two containers share the same network, making communication easy.

By default, the Voter API uses the URL `redis:6379` to establish a connection with the Redis container. This URL points to the Redis service within the Docker network. If you wish to use a different Redis instance or have a specific Redis server you'd like to connect to, you can configure this by setting the `REDIS_URL` environment variable for the Voter API container

## API Endpoints

The following are the available API endpoints for the Voter API:

1. **Welcome Message**: `GET /`  
  Returns a welcome message for the Voter API.

2. **List All Voters**: `GET /voters`  
  Retrieves a list of all voters along with their voting history.

3. **Get Voter by ID**: `GET /voters/:id`  
  Retrieves a single voter record based on the provided `:id`.

4. **Add Voter**: `POST /voters/:id`  
   Adds a new voter record with the provided `:id`, `firstName`, and `lastName`.

5. **Update Voter**: `PUT /voters/:id`  
   Updates an existing voter record with the provided `:id`, `firstName`, and `lastName`.

6. **Delete All Voters**: `DELETE /voters`  
   Deletes all voters and their voting history.

7. **Delete Voter by ID**: `DELETE /voters/:id`  
   Deletes a single voter record based on the provided `:id`.

8. **Get Voter History**: `GET /voters/:id/polls`  
   Retrieves the voting history of a voter based on the provided `:id`.

9. **Get Voter Poll**: `GET /voters/:id/polls/:pollid`  
   Retrieves a specific poll from a voter's voting history based on the provided `:id` and `:pollid`.

10. **Add Voter Poll**: `POST /voters/:id/polls/:pollid`  
    Adds a new poll record to a voter's voting history based on the provided `:id` and `:pollid`.

11. **Update Voter Poll**: `PUT /voters/:id/polls/:pollid`  
    Updates an existing poll record in a voter's voting history based on the provided `:id` and `:pollid`.

12. **Delete Voter Poll**: `DELETE /voters/:id/polls/:pollid`  
    Deletes a specific poll record from a voter's voting history based on the provided `:id` and `:pollid`.

13. **Health Check**: `GET /voters/health`  
    Provides health metadata for the Voter API, including status, uptime, total API calls, total API calls with errors, total request time, average request time, and boot time.

## API Usage Examples

### Build the Voter API executable
```bash
make build
```

### Run the Voter API from the code
```bash
make run
```

### Run the Voter API executable
```bash
make run-bin
```

### Get All Voters
```bash
make get-all
```

### Get a Voter by ID
```bash
make get-voter id=1
```

### Add a Voter
```bash
make add-voter id=2 firstName="John" lastName="Doe"
```

### Update a Voter
```bash
make update-voter id=2 firstName="John" lastName="Smith"
```

### Delete All Voters
```bash
make delete-all
```

### Delete a Voter by ID
```bash
make delete-voter id=2
```

### Get Voter History
```bash
make get-voter-history id=1
```

### Get Voter Poll
```bash
make get-voter-poll id=1 pollid=1
```

### Add Voter Poll
```bash
make add-voter-poll id=1 pollid=2
```

### Update Voter Poll
```bash
make update-voter-poll id=1 pollid=2
```

### Delete Voter Poll
```bash
make delete-voter-poll id=1 pollid=2
```

### Health Check
```bash
make health-check
```

## Middleware - Health Metadata

The Voter API includes a middleware that provides health metadata for the API. The following metadata is available:

- **Status**: Indicates the health status of the API.
- **Uptime**: The duration since the API was started.
- **TotalAPICalls**: The total number of API calls made to the Voter API.
- **TotalAPICallsError**: The total number of API calls that resulted in an error.
- **BootTime**: The timestamp when the API was started.
- **TotalRequestTime**: The total duration of all requests made to the API.
- **AverageRequestTime**: The average duration of a request to the API.
