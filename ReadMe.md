# Ghostvox Backend RESTful API

This is the backend for the **Ghostvox app**. It is a RESTful API built using Go, hosted on Fly.io with a PostgreSQL database. The API handles storage and retrieval of polls and their associated data for the Ghostvox app.

## Endpoints

### User Endpoints

#### 🚀 Create User ✅
- **Route:** `POST /api/v1/users`
- **Request:**
  ```json
  {
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "password": "yourpassword",
    "provider": "",        // optional, if using external providers
    "provider_id": "",     // optional, if using external providers
    "role": "user"
  }
  ```
- **Response (201 Created):**
  - **Cookies Set:**
    - `access_token`: A short-lived JWT stored as an HTTP-only cookie and also returned in the Authorization header as Bearer <token>.
    - `refresh_token`: A long-lived refresh token stored as an HTTP-only cookie.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```
- **Notes:**
  - The refresh token is not stored on the user record but in its own dedicated table in the database. This facilitates better token management (e.g., rotation, revocation, and multiple sessions).
  - Passwords are hashed using bcrypt before being stored.
  - The API uses the request context (r.Context()) to handle graceful cancellation of database operations.

#### 🔍 Get All Users ✅
- **Route:** `GET /api/v1/users`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "email": "john@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user"
    },
    {
      "id": "2",
      "email": "jane@example.com",
      "first_name": "Jane",
      "last_name": "Doe",
      "role": "admin"
    }
  ]
  ```
- **Notes:** This endpoint returns user information without any tokens.

#### 🔍 Get Single User ✅
- **Route:** `GET /api/v1/users/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "role": "user"
  }
  ```

#### ✏️ Update User ✅
- **Route:** `PUT /api/v1/users/{id}`
- **Request:**
  ```json
  {
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "password": "newpassword",
    "provider": "",
    "provider_id": "",
    "role": "user"
  }
  ```
- **Response (200 OK):**
  - **Cookies Set:** New `access_token` and `refresh_token` are issued as HTTP-only cookies, and the access token is included in the Authorization header.
  - **Response Body:**
    ```json
    {
      "message": "User updated successfully"
    }
    ```

#### ❌ Delete User ✅
- **Route:** `DELETE /api/v1/users/{id}`
- **Response (204 No Content)**

### Poll Endpoints

#### 🚀 Create Poll ✅
- **Route:** `POST /api/v1/polls`
- **Request:**
  ```json
  {
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Active|Inactive|Archived" // Case sensitive
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "active"
  }
  ```

#### 🔍 Get Poll ✅
- **Route:** `GET /api/v1/polls/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "active"
  }
  ```

#### 🔍 Get Polls
- **Route:** `GET /api/v1/polls`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "userId": "user123",
      "title": "Sample Poll",
      "description": "This is a sample poll description",
      "expiresAt": "2024-12-31T23:59:59Z",
      "status": "active"
    },
    {
      "id": "2",
      "userId": "user123",
      "title": "Another Sample Poll",
      "description": "This is another sample poll description",
      "expiresAt": "2024-12-31T23:59:59Z",
      "status": "active"
    }
  ]
  ```

#### ✏️ Update Poll ✅
- **Route:** `PUT /api/v1/polls/{id}`
- **Request:**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "title": "Updated Poll",
    "description": "This is an updated poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Inactive"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "title": "Updated Poll",
    "description": "This is an updated poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "inactive"
  }
  ```

#### ❌ Delete Poll ✅
- **Route:** `DELETE /api/v1/polls/{id}`
- **Response (204 No Content)**

### Vote Endpoints

#### 🚀 Create Vote
- **Route:** `POST /api/v1/polls/{id}/votes`
- **Request:**
  ```json
  {
    "userId": "user123",
    "optionId": "option1"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "optionId": "option1"
  }
  ```

#### 🔍 Get Vote
- **Route:** `GET /api/v1/polls/{id}/votes/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "optionId": "option1"
  }
  ```

#### 🔍 Get Votes
- **Route:** `GET /api/v1/polls/{id}/votes`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "userId": "user123",
      "optionId": "option1"
    },
    {
      "id": "2",
      "userId": "user456",
      "optionId": "option2"
    }
  ]
  ```

#### ❌ Delete Vote
- **Route:** `DELETE /api/v1/polls/{id}/votes/{id}`
- **Response (204 No Content)**

### Options Endpoint

#### 🚀 Create Option
- **Route:** `POST /api/v1/polls/{id}/options`
- **Request:**
  ```json
  {
    "options": [
      {
        "name": "Option Name",
        "value": "Option 1"
      },
      {
        "name": "Option Name",
        "value": "Option 2"
      }
    ]
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "text": "Option 1"
  }
  ```

#### 🔍 Get Option
- **Route:** `GET /api/v1/polls/{id}/options/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "text": "Option 1"
  }
  ```

#### 🔍 Get Options
- **Route:** `GET /api/v1/polls/{id}/options`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "userId": "user123",
      "text": "Option 1"
    },
    {
      "id": "2",
      "userId": "user456",
      "text": "Option 2"
    }
  ]
  ```

#### ✏️ Update Option
- **Route:** `PUT /api/v1/polls/{id}/options/{id}`
- **Request:**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "text": "Updated Option"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "userId": "user123",
    "text": "Updated Option"
  }
  ```

#### ❌ Delete Option
- **Route:** `DELETE /api/v1/polls/{id}/options/{id}`
- **Response (204 No Content)**
