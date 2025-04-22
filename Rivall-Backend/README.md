# Rivall Backend
This backend was constructed for the Rivall mobile app.  This project is part of a undergrad capstone project for Nathaniel Reeves at Utah Tech University in Jan 2025.

## Available Resources

The Rivall Backend API provides the following resources:

### Public Routes
- **POST /api/v1/auth/register**: Register a new user.
- **POST /api/v1/auth/login**: Log in an existing user.
- **POST /api/v1/auth/recovery/send-code**: Send an account recovery email.
- **POST /api/v1/auth/recovery/validate-code**: Validate an account recovery code.
- **GET /api/v1/contacts/{user_id}**: Retrieve a user's contacts.

### Private Routes (Require Authentication)
- **GET /api/v1/users/{user_id}**: Retrieve user details.
- **PUT /api/v1/auth/recovery/{user_id}/reset-password**: Reset a user's password.
- **POST /api/v1/auth/{user_id}/refresh**: Renew an access token.
- **DELETE /api/v1/auth/{user_id}/logout**: Log out a user.
- **POST /api/v1/users/{user_id}/contacts**: Add a new contact for a user.
- **GET /api/v1/users/{user_id}/contacts/{chat_id}/chat**: Retrieve a chat for a specific contact.
- **GET /api/v1/ws/connect/{user_id}**: Establish a WebSocket connection.

## Starting the Service

To start the Rivall Backend service, follow these steps:

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/learning-cloud-native-go/Rivall-Backend.git
    cd Rivall-Backend
    ```

2. **Set Up Environment Variables**:
    Create a `.env` file in the root directory with the following variables:
    ```env
    MONGO_URI=<your-mongodb-atlas-uri>
    JWT_SECRET=<your-jwt-secret>
    ```

3. **Build and Run the Service**:
    Using Docker:
    ```bash
    docker build -t rivall-backend .
    docker run -p 8080:8080 --env-file .env rivall-backend
    ```
    Or, using Go directly:
    ```bash
    go mod download
    go run main.go
    ```

4. **Access the API**:
    The API will be available at `http://localhost:8080`.

## Database Interaction

The Rivall Backend uses MongoDB Atlas for data storage. The database is connected using the `ConnectMongoDB` function in `main.go`. Key details include:

- **Connection**: The MongoDB URI is retrieved from the environment variables. The connection is established using the official MongoDB Go driver.
- **Ping**: A ping command is sent to ensure the connection is successful.
- **Collections**: Data is stored in collections such as `users`, `contacts`, and `chats`.
- **CRUD Operations**: The backend performs Create, Read, Update, and Delete operations on the database to manage user data, authentication, and chat functionality.

Make sure your MongoDB Atlas cluster is properly configured and accessible from your environment.
