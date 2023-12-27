# GoFiber Hexagonal Architecture with PostgreSQL

This repository contains a Golang application that implements a robust architecture using GoFiber as the HTTP framework and PostgreSQL for database management. The application follows the hexagonal architecture pattern, emphasizing modularity, testability, and maintainability.

## Overview

The project is focused on creating and managing various endpoints, including user authentication, product, order, category, and handling image files uploaded to Google Cloud Platform Storage.

## Features

- **User Management:** Handles user operations such as registration, login, and user profile management.
- **Authentication:** Implements secure authentication mechanisms such as access token and refresh token for accessing endpoints.
- **Product and Category Management:** Manages products and their categories, supporting CRUD operations.
- **Order Handling:** Implements functionalities for order creation, modification, and tracking.
- **Google Cloud Platform Integration:** Enables storage and retrieval of image files in Google Cloud Platform Storage and use Cloud SQL to Deploy databases and Deploy Application on Cloud Run

## Prerequisites

- Go (v1.16 or higher)
- PostgreSQL
- Google Cloud Platform account

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/NatthawutSK/ri-shop.git
2. **Create .env file:**
   ```bash
   APP_HOST=
   APP_PORT=
   APP_NAME=
   APP_VERSION=
   APP_BODY_LIMIT=
   APP_READ_TIMEOUT=
   APP_WRITE_TIMEOUT=
   APP_FILE_LIMIT=
   APP_GCP_BUCKET=
   
   JWT_SECRET_KEY=
   JWT_API_KEY=
   JWT_ADMIN_KEY=
   JWT_ACCESS_EXPIRES=
   JWT_REFRESH_EXPIRES=
   
   DB_HOST=
   DB_PORT=
   DB_PROTOCOL=
   DB_USERNAME=
   DB_PASSWORD=
   DB_DATABASE=
   DB_SSL_MODE=
   DB_MAX_CONNECTIONS=
3. **Run Command:**
   ```bash
   go run main.go .
