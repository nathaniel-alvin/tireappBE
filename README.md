Here's a comprehensive README for the Tire Eye application GitHub repository:

---

# Tire Eye

## Development of Tire Eye: An AI-Powered Vehicle Tire Management Application

### Backend and Database Development

---

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

---

## Introduction

Tire Eye is a mobile application designed to assist users in identifying, managing, and locating tires and tire shops. The application aims to enhance vehicle safety and maintenance by providing detailed tire information, tracking tire usage, and simplifying the process of finding replacement tires and nearby tire shops. 

## Features

- **Automatic Tire Identification:** Upload images of tire walls to identify tire details such as size, brand, type, and DOT.
- **Tire Management:** Save and track tire information and maintenance history.
- **Tire Marketplace Search:** Search for replacement tires in online marketplaces.
- **Nearby Tire Shops:** Locate tire shops using Google Maps integration.

## Architecture

Tire Eye's backend system is developed using:
- **Go:** For server-side logic.
- **PostgreSQL:** For database management.
- **Amazon S3:** For image storage.

The application follows the MVC (Model-View-Controller) architecture to ensure a clear separation of concerns and maintainability.

## Installation

### Prerequisites

- Go 1.18 or later
- PostgreSQL 13 or later
- AWS account for S3

### Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/nathaniel-alvin/tireappBE.git
    cd tireappBE
    ```

2. Set up environment variables:
    ```bash
    cp .env.example .env
    # Edit .env with your configuration
    ```

3. Install dependencies:
    ```bash
    go mod tidy
    ```

4. Set up PostgreSQL database:
    ```sql
    CREATE DATABASE tire_eye;
    ```

5. Run migrations:
    ```bash
    make migrate-up
    ```

6. Start the server:
    ```bash
    make run
    ```

## Usage

To run the application, ensure your environment variables are correctly configured in the `.env` file. Start the server using the `go run main.go` command. The application will be accessible at `http://localhost:8080`.

## API Endpoints

### Authentication
- **POST** `/api/v1/auth/register` - Register a new user
- **POST** `/api/v1/auth/login` - Log in a user

### Tires
- **POST** `/api/v1/inventories` - Add a new tire
- **GET** `/api/v1/inventories` - List all tires
- **GET** `/api/v1/inventories/{id}` - Get tire details
- **PUT** `/api/v1/inventories/{id}` - Update tire information
- **DELETE** `/api/v1/inventories/{id}` - Delete a tire

### Cars
- **POST** `/api/v1/cars` - Add a new car
- **GET** `/api/v1/cars` - List all cars
- **GET** `/api/v1/cars/{id}` - Get car details
- **PUT** `/api/v1/cars/{id}` - Update car information
- **DELETE** `/api/v1/cars/{id}` - Delete a car

### Shops
- **GET** `/api/v1/shops` - Search for nearby tire shops

## Testing

Currently, all testing of the Tire Eye backend is conducted using Postman.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes. Ensure your code adheres to the project's coding standards and passes all tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.

## Acknowledgements

- **Supervisor:** S. Pradono Suryodiningrat
- **Team Members:** Bernard Wijaya & Jason Mikael
- **Special Thanks:** N. Nurul Qomariyah as co-supervisor

---

Feel free to reach out with any questions or feedback. Thank you for contributing to Tire Eye!

---
