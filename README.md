# Todo List API with User Authentication

This project is a RESTful API for managing a todo list, allowing users to register, log in, and create, update, delete, and retrieve their todo items. The API uses JWT-based authentication for secure access, and data is persisted in MongoDB database.

## Features
- **User Authentication**: Users can register and log in using JWT for secure authentication.
- **Task Management**: Users can create, update, delete, and retrieve their own todo items.
- **JWT-based Authentication**: Each user’s session is managed with JWT, ensuring secure access.
- **Database Persistence**: Todo items and user data are stored in a NoSQL (MongoDB) database.
  
## Technologies Used
- **Go**: Programming language used for building the API.
- **Gorilla Mux**: Router for handling routes and endpoints.
- **MongoDB**: Database for persisting user and todo data.
- **JWT**: For secure token-based authentication.
  
## Project Structure
![image](https://github.com/user-attachments/assets/d5852879-7b05-4d84-a1a2-93864be668d0)


## API Endpoints
### Authentication
- **Register**: `POST /register`
    - Request body: 
      ```json
      {
        "username": "test1",
        "password": "test1"
      }
      ```
    - Response: 
      ```json
       {
        "token": "jwt_token"
      }
      ```

- **Login**: `POST /login`
    - Request body: 
      ```json
      {
        "username": "test1",
        "password": "test1"
      }
      ```
    - Response: 
      ```json
      {
        "token": "jwt_token"
      }
      ```
### Todo Management (Requires Authentication)
- **Get All Todos with Pagination**: `GET /todos?page=1&limit=2`
    - Request headers:
      - `Authorization: Bearer <jwt_token>`
    - Query parameters:
      - `page`: Specifies the page number (default: 1)
      - `limit`: Specifies the number of todos to retrieve per page (default: 10)
    - Response:
      ```json
      {
        "data": [
          {
            "ID": "66edd306d77b65a2d931234d",
            "UserID": "test1",
            "Title": "walking",
            "Done": false
          },
          {
            "ID": "66edd30bd77b65a2d931234e",
            "UserID": "test1",
            "Title": "Dancing",
            "Done": false
          }
        ],
        "limit": 2,
        "page": 1,
        "success": true
      }
      ```
- **Create Todo**: `POST /todos`
    - Request body:
      ```json
      {
        "title": "Buy groceries",
      }
      ```
    - Response:
      ```json
      {
       "data": {
		      "ID": "66edd31dd77b65a2d9312350",
		      "UserID": "test1",
		      "Title":  "Buy groceries",
		      "Done": false
	     },
	     "success": true
      }
      ```
  - **Update Todo**: `PUT /todos?id=66edc75d5746feae033bdd61`
    - Request body:
      ```json
      {
        "title": "Update Title",
      }
      ```
    - Response:
      ```json
      {
        "message": "Todo updated successfully",
        "success": true
      }
      ```
  - **Delete Todo**: `DELETE /todos?id=66edc75d5746feae033bdd61`
    - Response:
      ```json
      {
        "message": "Todo deleted successfully",
        "success": true
      }
      ```
- **Update Todo Status**: `PUT /todos/status?id=66edc75d5746feae033bdd61`
    - Response:
      ```json
      {
        "message": "Todo status updated successfully.",
        "success": true
      }
      ```
