# status-board-demo

## Backend
### Status Board API

The backend-server directory contains a RESTful API for a status board built using Go, Echo framework, and PostgreSQL.
The API allows you to create, update, delete, and retrieve status entries. 
It also provides Server-Sent Events (SSE) for real-time notifications of status changes.

#### Setting Up the Database

First, you need to set up the PostgreSQL database for the application. Follow these steps:

1. Create a new PostgreSQL database.
2. Set the `POSTGRESQL_ADDON_URI` environment variable to your PostgreSQL connection string.
3. Run the `sql/init.sql` SQL script with your pgsql client to create the `status` table and necessary triggers for notifications

### Building the Application

To build and run the application, follow these steps:

1. Install Go (version 1.16 or later) from [the official website](https://golang.org/dl/).
2. Clone this repository: `git clone https://github.com/fredericalix/status-board-demo.git`.
3. Change to the repository directory: `cd status-board-demo/backend-server`.
4. Install dependencies: `go mod download`.
5. Build the application: `go build -o status-board-api`.
6. Run the application: `./status-board-api`.

The API will listen on port `8080` or the port specified in the `PORT` environment variable.

### Using the API

Here are some example `curl` commands to interact with the API:

- Get API status:

  ```sh
  curl -X GET http://localhost:8080/
  ```

- Create a new status:

  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"designation": "Sample Status", "state": "green"}' http://localhost:8080/status
  ```

- Update an existing status (replace `<id>` with the status ID):

  ```sh
  curl -X PUT -H "Content-Type: application/json" -d '{"designation": "Updated Status", "state": "red"}' http://localhost:8080/status/<id>
  ```

- Delete a status (replace `<id>` with the status ID):

  ```sh
  curl -X DELETE http://localhost:8080/status/<id>
  ```

- Get a status by ID (replace `<id>` with the status ID):

  ```sh
  curl -X GET http://localhost:8080/status/<id>
  ```

- Get all statuses:

  ```sh
  curl -X GET http://localhost:8080/status
  ```

### Unit Tests

Set the `POSTGRESQL_ADDON_URI` environment variable to your PostgreSQL connection string.

```
cd backend-server
go test -v ./api
```

### Swagger Documentation

To view the API documentation using Swagger, visit the following URL:

```
http://localhost:8080/swagger/index.html
```

Replace `localhost` with your server's hostname if you're not running the application locally.

## Frontend

The frontend-server directory contains a Vue.js frontend application for displaying notifications fetched from the backend-server API.
The application uses the Vue.js framework, and it includes features like real-time updates through SSE.

### Features

- Fetches notifications from a backend API
- Displays notifications as colored squares
- Real-time updates with Server-Sent Events (SSE)
- Responsive and clean design
### Installation and Setup

1. Change the working directory to the frontend's root folder:

```
cd frontend-server
```

1. Install the required dependencies:

```
npm install
```

4. Define the `REACT_APP_API_URL` environment variable if you want to point to a different backend API URL. If not set, the default value will be `http://localhost:8080`.

```
export REACT_APP_API_URL="http://yourbackendurl.com"
```

5. Start the development server:

```
npm run dev
```

The React frontend application should now be running on `http://localhost:3000` (or another available port).

### Usage

Open the application in your browser and observe the notifications displayed as colored squares. The squares' colors correspond to the `state` attribute of each notification.
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.