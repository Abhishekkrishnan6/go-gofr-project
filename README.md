# Car Service Center Application

The Car Service Center Application is a web-based system that allows users to manage and track cars in a service center. It provides functionalities such as adding new cars, updating their status, searching by license plate, and deleting entries.

## Features

- **Add New Car:** Users can add new cars to the service center by entering the license plate information.

- **Update Car Status:** The application allows users to update the status of a car, indicating whether it's washed, done, or completed.

- **Search by License Plate:** Users can search for cars by their license plate, making it easy to find specific entries.

- **Delete Car Entry:** Users can delete car entries from the system, removing them from the service center records.

## Technologies Used

- **Frontend:** React.js
- **Backend:** Gofr (a lightweight HTTP framework for Go)
- **Database:** MongoDB
- **Styling:** CSS

## Setup and Installation

### Prerequisites

- Node.js and npm installed (for the frontend)
- Go installed (for the backend)
- MongoDB installed and running locally

### Frontend

1. Navigate to the `frontend` directory.

```bash
cd frontend

2. Install dependencies.

```bash
npm install

3. Run the frontend application.

```bash
npm start

The frontend will be accessible at http://localhost:3000.

4. Backend

Navigate to the backend directory.

```bash
cd backend

5. Install Gofr using go get.

```bash
go get -u gofr.dev/pkg/gofr

6. Run the backend server.

```bash
go run main.go
The backend will be accessible at http://localhost:8000.
