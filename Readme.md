# Receipt Processor

A simple web service that processes receipts and awards points based on defined rules.
The service exposes a RESTful API that allows you to submit receipts and later retrieve the points awarded.

## Table of Contents

- [Overview](#overview)
- [API Specification](#api-specification)
- [Building and Running](#building-and-running)
- [Running Tests with Docker](#running-tests-with-docker)

## Overview

The Receipt Processor service accepts a JSON receipt via the **POST /receipts/process** endpoint and returns an ID. 
You can then retrieve the points awarded for that receipt by calling the **GET /receipts/{id}/points** endpoint.
The points calculation is based on several rules (such as counting alphanumeric characters in the retailer's name, checking total amounts, bonus for certain item descriptions, etc).

## API Specification

Key endpoints include:

- **POST /receipts/process**
  - **Description:** Submits a receipt for processing.
  - **Request:** JSON receipt conforming to the `Receipt` schema.
  - **Response:** JSON with an `id` field.
  
- **GET /receipts/{id}/points**
  - **Description:** Returns the points awarded for the receipt identified by the provided ID.
  - **Response:** JSON with a `points` field.

## Building and Running

1. Build the Docker image:

```bash
    docker-compose build app
```

2. Run the application container:

```bash
    docker-compose up app
```

The service will be available at http://localhost:8080.

## Running Tests with Docker

A dedicated service in `docker-compose.yml` buils the image from the builder stage(which has go installed) and runs the tests.

1. Run tests:

```bash
    docker-compose up test --build
```

To run it locally:

```bash
    go test -v ./...
```
Built with ❤️ by Sirius829[https://github.com/sirius829]