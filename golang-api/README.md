# Golang Volley Live Score

This repository is part of a larger project that includes Golang, React, and Podman/Docker. The goal of the project is to provide live scores and updates for a volleyball team, allowing parents to stay informed about their children's performance in matches. This specific repository contains the backend code for the project.

## Project Overview

The **Golang Volley Live Score** project aims to:

- Provide real-time updates on volleyball matches.
- Allow one person to update the score while others receive updates on their devices (cell phones or computers).
- Focus on keeping parents informed about their children's performance.

## Features

- **Matches Management**: Create, delete, update, and get matches.
- **Sets Management**: Create, delete, update, and get sets within matches.
- **Live Updates**: The backend sends live updates to the frontend when sets are updated.

## Technologies Used

- **Golang**: The primary language for the backend.
- **React**: Used for the frontend (not included in this repository).
- **Podman**: For containerization and deployment (preferred for security reasons).
- **Docker**: An alternative option for containerization.
- **Nix**: Used for creating reproducible development environments.

## Why Podman?

This project uses **Podman** for containerization due to its fully open-source nature and enhanced security compared to Docker. However, you can still use **Docker** if that's your preference.

## Getting Started

### Prerequisites

- Golang installed on your machine.
- Podman or Docker installed on your machine.
- Nix installed on your machine (optional, for development).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/JorgeSaicoski/volley-live-score.git
   cd golang-api
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run cmd/matches/main.go
   ```

### Running with Podman

1. Build the Podman image:

   ```bash
   podman build -t golang-volley-live-score .
   ```

2. Run the Podman container:

   ```bash
   podman run -p 8080:8080 golang-volley-live-score
   ```

Alternatively, you can use Podman Compose for orchestrating multiple services, such as the matches and sets services:

1. Start the application using Podman Compose:

   ```bash
   podman-compose up --build
   ```

2. Access the application at `http://localhost:8080` or `http://localhost:8081` depending on the service.

### Running with Docker

The project is compatible with Docker as well. You can follow these steps to run it using Docker:

1. Build the Docker image:

   ```bash
   docker build -t golang-volley-live-score .
   ```

2. Run the Docker container:

   ```bash
   docker run -p 8080:8080 golang-volley-live-score
   ```

3. Alternatively, use Docker Compose:

   ```bash
   docker-compose up --build
   ```

### Using Nix for Development

You can use Nix to ensure a consistent development environment. Nix handles dependencies like Podman/Docker and Golang, making the development setup reproducible.

1. Install Nix: [Nix Installation Guide](https://nixos.org/manual/nix/stable/#chap-installation).

2. Enter the Nix Shell:

   ```bash
   cd nix
   nix-shell
   ```

   This will open a shell with Podman/Docker and Golang installed and ready to use.

### API Endpoints

#### Matches

- `GET /matches`: Get all matches (with pagination).
- `GET /matches/:id`: Get a match by ID.
- `POST /matches`: Create a new match.
- `PUT /matches/:id`: Update a match.
- `DELETE /matches/:id`: Delete a match.

#### Sets

- `GET /sets`: Get all sets.
- `GET /sets/:id`: Get a set by ID.
- `POST /sets`: Create a new set.
- `PUT /sets/:id`: Update a set.
- `DELETE /sets/:id`: Delete a set.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

### Recommendations

Some code is repeated in the project. Refactoring and updating this code will help improve maintainability.

## License

This project is licensed under the MIT License.


