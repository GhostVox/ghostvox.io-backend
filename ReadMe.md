# Ghostvox Backend RESTful API

This is the backend for the **Ghostvox app**. It is a RESTful API built using Go, hosted on Fly.io with a PostgreSQL database. The API handles storage and retrieval of  polls and their data for the Ghostvox app.

---

## Endpoints
&nbsp;
&nbsp;
### User Endpoints

#### 🚀 Create User ✅
- **Route:** `POST /api/v1/users`
- **Request:**
  ```json
  {
    "name": "John",
    "email": "john@example.com",
    "last_name": "Smith",
    "user_token": "abc123",
    "role": "user"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "1",
    "name": "John",
    "email": "john@example.com",
    "last_name": "Smith",
    "user_token": "abc123",
    "role": "user"
  }
  ```
&nbsp;
#### 🔍 Get All Users ✅
- **Route:** `GET /api/v1/users`
- **Response (200 OK):**
  ```json
  [
    {
      "id": "1",
      "name": "John",
      "email": "john@example.com",
      "last_name": "Smith",
      "user_token": "abc123",
      "role": "user"
    },
    {
      "id": "2",
      "name": "Jane",
      "email": "jane@example.com",
      "last_name": "Doe",
      "user_token": "def456",
      "role": "admin"
    }
  ]
  ```
&nbsp;
#### 🔍 Get Single User ✅
- **Route:** `GET /api/v1/users/{id}`
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "name": "John",
    "email": "john@example.com",
    "last_name": "Smith",
    "user_token": "abc123",
    "role": "user"
  }
  ```
&nbsp;
#### ✏️ Update User ✅
- **Route:** `PUT /api/v1/users/{id}`
- **Request:**
  ```json
  {
    "id": "1",
    "name": "John",
    "email": "john@example.com",
    "last_name": "Smith",
    "user_token": "abc123",
    "role": "user"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "id": "1",
    "name": "John",
    "email": "john@example.com",
    "last_name": "Smith",
    "user_token": "abc123",
    "role": "user"
  }
  ```
&nbsp;
#### ❌ Delete User ✅
- **Route:** `DELETE /api/v1/users/{id}`
- **Response (204 No Content)**

---
&nbsp;
&nbsp;
### Poll Endpoints
&nbsp;
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
&nbsp;
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
&nbsp;
### 🔍 Get Polls
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
&nbsp;
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
&nbsp;
#### ❌ Delete Poll ✅
- **Route:** `DELETE /api/v1/polls/{id}`
- **Response (204 No Content)**

&nbsp;
### Vote Endpoints
&nbsp;
#### 🚀 Create Vote
- **Route:** `POST /api/v1/polls/{id}/votes`
- **Request:**
  ```json
  {
    "userId": "user123",
    "optionId": "option1",

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
&nbsp;
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
&nbsp;
#### ✏️ Get Votes
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
&nbsp;
#### ❌ Delete Vote
- **Route:** `DELETE /api/v1/polls/{id}/votes/{id}`
- **Response (204 No Content)**

&nbsp;
### Options Endpoint
&nbsp;
#### 🚀 Create Option
- **Route:** `POST /api/v1/polls/{id}/options`
- **Request:**
  ```json
  {
    "options": [
      {
        "name": "Option Name",
        "value": "Option 1",
      },
      {
        "name": "Option Name",
        "value": "Option 2",
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
&nbsp;
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
&nbsp;
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
&nbsp;
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
&nbsp;
#### ❌ Delete Option
- **Route:** `DELETE /api/v1/polls/{id}/options/{id}`
- **Response (204 No Content)**
