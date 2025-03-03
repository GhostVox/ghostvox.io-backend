# Ghostvox Backend RESTful API

This is the backend for the **Ghostvox app**. It is a RESTful API built using Go, hosted on Fly.io with a PostgreSQL database. The API handles storage and retrieval of polls and their associated data for the Ghostvox app.

## API Endpoints

### Authentication Endpoints

#### 🚀 Register User
- **Route:** `POST /api/v1/auth/register`
- **Request:**
  ```json
  {
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "password": "yourpassword",
    "provider": "",
    "provider_id": "",
    "role": "user"
  }
  ```
- **Response (201 Created):**
  - **Cookies Set:**
    - `accessToken`: A short-lived JWT stored as an HTTP-only cookie and also returned in the Authorization header as Bearer <token>.
    - `refreshToken`: A long-lived refresh token stored as an HTTP-only cookie.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```

#### 🔑 Login
- **Route:** `POST /api/v1/auth/login`
- **Request:**
  ```json
  {
    "email": "john@example.com",
    "password": "yourpassword"
  }
  ```
- **Response (201 Created):**
  - **Cookies Set:**
    - `accessToken`: A short-lived JWT stored as an HTTP-only cookie and also returned in the Authorization header as Bearer <token>.
    - `refreshToken`: A long-lived refresh token stored as an HTTP-only cookie.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```

#### 🔄 Refresh Token
- **Route:** `POST /api/v1/auth/refresh`
- **Request:** No body needed (uses HTTP-only cookie)
- **Response (201 Created):**
  - **Cookies Set:**
    - New `accessToken` and `refreshToken` are issued as HTTP-only cookies.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```

#### 🚪 Logout
- **Route:** `POST /api/v1/auth/logout`
- **Request:** No body needed (uses HTTP-only cookie)
- **Response (200 OK):**
  - **Cookies:** Clears authentication cookies
  - **Response Body:**
    ```json
    {
      "message": "User logged out successfully"
    }
    ```

#### 🔑 Google OAuth Login
- **Route:** `GET /api/v1/auth/google/login`
- **Response:** Redirects to Google authentication

#### 🔑 Google OAuth Callback
- **Route:** `GET /api/v1/auth/google/callback`
- **Response (201 Created):**
  - **Cookies Set:**
    - `accessToken` and `refreshToken` are issued as HTTP-only cookies.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```

### User Endpoints

#### 🔍 Get All Users (Admin only)
- **Route:** `GET /api/v1/admin/users`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "email": "john@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": "2",
      "email": "jane@example.com",
      "first_name": "Jane",
      "last_name": "Doe",
      "role": "admin",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
  ```

#### 🔍 Get Single User (Admin only)
- **Route:** `GET /api/v1/admin/users/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "role": "user",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

#### ✏️ Update User
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
  - **Cookies Set:** New `accessToken` and `refreshToken` are issued as HTTP-only cookies, and the access token is included in the Authorization header.
  - **Response Body:**
    ```json
    {
      "message": "User created successfully"
    }
    ```

#### ❌ Delete User
- **Route:** `DELETE /api/v1/users/{id}`
- **Response (204 No Content)**

### Poll Endpoints

#### 🚀 Create Poll
- **Route:** `POST /api/v1/polls`
- **Request:**
  ```json
  {
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Active"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "poll-uuid",
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "created_at": "2023-05-01T10:00:00Z",
    "updated_at": "2023-05-01T10:00:00Z",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Active"
  }
  ```

#### 🔍 Get Poll
- **Route:** `GET /api/v1/polls/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "poll-uuid",
    "userId": "user123",
    "title": "Sample Poll",
    "description": "This is a sample poll description",
    "created_at": "2023-05-01T10:00:00Z",
    "updated_at": "2023-05-01T10:00:00Z",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Active"
  }
  ```

#### 🔍 Get All Polls
- **Route:** `GET /api/v1/polls`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "poll-uuid-1",
      "userId": "user123",
      "title": "Sample Poll",
      "description": "This is a sample poll description",
      "created_at": "2023-05-01T10:00:00Z",
      "updated_at": "2023-05-01T10:00:00Z",
      "expiresAt": "2024-12-31T23:59:59Z",
      "status": "Active"
    },
    {
      "id": "poll-uuid-2",
      "userId": "user123",
      "title": "Another Sample Poll",
      "description": "This is another sample poll description",
      "created_at": "2023-05-02T10:00:00Z",
      "updated_at": "2023-05-02T10:00:00Z",
      "expiresAt": "2024-12-31T23:59:59Z",
      "status": "Active"
    }
  ]
  ```

#### ✏️ Update Poll
- **Route:** `PUT /api/v1/polls/{id}`
- **Request:**
  ```json
  {
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
    "id": "poll-uuid",
    "userId": "user123",
    "title": "Updated Poll",
    "description": "This is an updated poll description",
    "created_at": "2023-05-01T10:00:00Z",
    "updated_at": "2023-05-03T15:30:00Z",
    "expiresAt": "2024-12-31T23:59:59Z",
    "status": "Inactive"
  }
  ```

#### ❌ Delete Poll
- **Route:** `DELETE /api/v1/polls/{id}`
- **Response (204 No Content)**

### Poll Option Endpoints

#### 🚀 Create Options
- **Route:** `POST /api/v1/polls/{pollId}/options`
- **Request:**
  ```json
  {
    "options": [
      {
        "name": "Option 1",
        "value": "Value 1"
      },
      {
        "name": "Option 2",
        "value": "Value 2"
      }
    ]
  }
  ```
- **Response (201 Created):**
  ```json
  [
    {
      "id": "option-uuid-1",
      "name": "Option 1",
      "value": "Value 1",
      "poll_id": "poll-uuid",
      "created_at": "2023-05-01T10:10:00Z",
      "updated_at": "2023-05-01T10:10:00Z"
    },
    {
      "id": "option-uuid-2",
      "name": "Option 2",
      "value": "Value 2",
      "poll_id": "poll-uuid",
      "created_at": "2023-05-01T10:10:00Z",
      "updated_at": "2023-05-01T10:10:00Z"
    }
  ]
  ```

#### 🔍 Get Option
- **Route:** `GET /api/v1/polls/{pollId}/options/{optionId}`
- **Response (200 OK):**
  ```json
  {
    "id": "option-uuid",
    "name": "Option 1",
    "value": "Value 1",
    "poll_id": "poll-uuid",
    "created_at": "2023-05-01T10:10:00Z",
    "updated_at": "2023-05-01T10:10:00Z"
  }
  ```

#### 🔍 Get All Options for Poll
- **Route:** `GET /api/v1/polls/{pollId}/options`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "option-uuid-1",
      "name": "Option 1",
      "value": "Value 1",
      "poll_id": "poll-uuid",
      "created_at": "2023-05-01T10:10:00Z",
      "updated_at": "2023-05-01T10:10:00Z"
    },
    {
      "id": "option-uuid-2",
      "name": "Option 2",
      "value": "Value 2",
      "poll_id": "poll-uuid",
      "created_at": "2023-05-01T10:10:00Z",
      "updated_at": "2023-05-01T10:10:00Z"
    }
  ]
  ```

#### ✏️ Update Option
- **Route:** `PUT /api/v1/polls/{pollId}/options/{optionId}`
- **Request:**
  ```json
  {
    "id": "option-uuid",
    "name": "Updated Option",
    "value": "Updated Value"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": "option-uuid",
    "name": "Updated Option",
    "value": "Updated Value",
    "poll_id": "poll-uuid",
    "created_at": "2023-05-01T10:10:00Z",
    "updated_at": "2023-05-03T16:20:00Z"
  }
  ```

#### ❌ Delete Option
- **Route:** `DELETE /api/v1/polls/{pollId}/options/{optionId}`
- **Response (204 No Content)**

### Vote Endpoints

#### 🚀 Create Vote
- **Route:** `POST /api/v1/polls/{pollId}/votes`
- **Request:**
  ```json
  {
    "userId": "user123",
    "optionId": "option-uuid"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "vote-uuid",
    "pollId": "poll-uuid",
    "optionId": "option-uuid",
    "userId": "user123",
    "created_at": "2023-05-01T11:00:00Z"
  }
  ```

#### 🔍 Get Votes by Poll
- **Route:** `GET /api/v1/polls/{pollId}/votes`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "vote-uuid-1",
      "pollId": "poll-uuid",
      "optionId": "option-uuid-1",
      "userId": "user123",
      "created_at": "2023-05-01T11:00:00Z"
    },
    {
      "id": "vote-uuid-2",
      "pollId": "poll-uuid",
      "optionId": "option-uuid-2",
      "userId": "user456",
      "created_at": "2023-05-01T11:30:00Z"
    }
  ]
  ```

#### ❌ Delete Vote
- **Route:** `DELETE /api/v1/votes/{voteId}`
- **Response (204 No Content)**

## Technical Notes

- **Authentication:** The API uses JWT tokens for authentication, with both access and refresh tokens.
- **Cookie Security:** Authentication tokens are stored as HTTP-only cookies with appropriate security settings.
- **Database:** The API uses PostgreSQL with foreign key constraints and cascading deletes.
- **Transaction Support:** Critical operations like user creation and token management use database transactions to ensure data consistency.
- **Role-Based Access Control:** Certain endpoints are restricted to admin users.
- **OAuth Integration:** Google OAuth is supported for authentication.

## Deployment

The application can be deployed using Docker. See the Docker documentation for more details.

```bash
docker compose up --build
```

The API will be available at http://localhost:8080.
