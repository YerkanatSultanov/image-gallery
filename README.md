# Image Gallery Service

## Overview

The Image Gallery Service is a microservices-based application designed to manage user authentication, image galleries, and user information. It provides a scalable and modular architecture that leverages technologies such as Go, gRPC, Kafka, and Docker to ensure flexibility and efficiency.

## Project Structure

The project is organized into several components:

- **cmd**: Contains the main entry points for different microservices (auth, gallery, user).
- **config**: Configuration files for each microservice, stored in YAML format.
- **database**: SQL migration files for initializing and updating the database schema.
- **docker**: Docker configuration, including Docker Compose for orchestrating services.
- **internal**: Core application logic and functionality, divided into microservices (auth, gallery, user).
- **pkg**: Shared packages and utilities, including protobuf files for gRPC communication.

## Microservices

### Auth Microservice

- Manages user authentication and token generation.
- Utilizes gRPC for communication.
- Configuration stored in `config/auth/config.yaml`.

### Gallery Microservice

- Handles image gallery-related operations.
- Implements gRPC for communication with other services.
- Configuration stored in `config/gallery/config.yaml`.

### User Microservice

- Manages user information and verification.
- Uses gRPC for inter-service communication.
- Configuration stored in `config/user/config.yaml`.

## Database

- SQL migration files for setting up and updating the database schema.
- Separated into `gallery_migrations` and `user_migrations` for clarity.

## Docker Configuration

- Docker Compose file (`docker/docker-compose.yaml`) for orchestrating Kafka and Zookeeper services.
- Data directories for persisting Kafka and Zookeeper data.

## Internal Packages

- **Auth**: User authentication logic.
- **Gallery**: Image gallery-related functionality.
- **Kafka**: Kafka consumer and producer implementations.
- **User**: User information and verification logic.
- **Util**: Shared utilities, including token handling and password utilities.

## Protobuf

- Protobuf files for defining gRPC service contracts and message structures.
- Located in the `pkg/protobuf` directory.

## Getting Started

1. Clone the repository: `git clone https://github.com/your-username/image-gallery.git`
2. Navigate to the project directory: `cd image-gallery`
3. Run `docker-compose up` to start Kafka and Zookeeper services.
4. Build and run individual microservices as needed (`cmd/auth`, `cmd/gallery`, `cmd/user`).

## Dependencies

- Go (version specified in `go.mod`)
- Docker and Docker Compose
- Kafka and Zookeeper

## Contributing

Feel free to contribute to the project by opening issues, submitting pull requests, or suggesting improvements. See the [CONTRIBUTING.md](CONTRIBUTING.md) file for more details.
