# Ghostvox backend RESTFUL API
>This is the backend for the ***Ghostvox app***. It is a RESTFUL API that is built using Go.
It is hosted on fly.io and uses a Postgres database.
The API is used to store and retrieve audio files and metadata for the Ghostvox app.

## Endpoints
>### CURRENT/POLLS
  >- GET v1/current/polls
  >- GET v1/current/polls/{id}
  >- POST v1/current/polls
  >- PUT v1/current/polls/{id}
  >- DELETE v1/current/polls/{id}
  #### RESPONSE
  ```json
  {
    "id": 1,
    "question": "What is your favorite color?",
    "options": [
      {
        "id": 1,
        "text": "Red"
      },
      {
        "id": 2,
        "text": "Blue"
      },
      {
        "id": 3,
        "text": "Green"
      }
    ]


  }
  ```

  ## FINISHED/POLLS
  >- GET v1/finished/polls
  >- GET v1/finished/polls/{id}

### RESPOSNSE
  ```json
  {
    "id": 1,
    "question": "What is your favorite color?",
    "options": [
      {
        "id": 1,
        "text": "Red",
        "votes": 5
      },
      {
        "id": 2,
        "text": "Blue",
        "votes": 3
      },
      {
        "id": 3,
        "text": "Green",
        "votes": 2
      }
    ],
  "winning_result": "option_id"
  }
  ```

  ## VOTES
  >- POST v1/votes
  #### REQUEST
  ```json
  {
    "poll_id": 1,
    "option_id": 1,
  "user_token": "1234567890"
  }
  ```
  #### RESPONSE
  ```json
  {
    "id": 1,
    "poll_id": 1,
    "option_id": 1
  "status_code": code
  }
  ```

  ## USERS
  >- GET v1/users
  >- GET v1/users/{id}
  >- POST v1/users
  >- PUT v1/users/{id}
  >- DELETE v1/users/{id}
  ### REQUEST
  ```json
  {
    "id": 1,
    "name": "Brent Harrington",
    "email": "bob@gmail.com",
    "user_token": "password",
    "ip_address": "",
    "user_agent": "",
    "created_at": "2021-07-01T00:00:00Z",
  "updated_at": "2021-07-01T00:00:00Z"
  }
  ```

  #### RESPONSE
  ```json
  {
    "id": 1,
    "name": "Brent Harrington",
    "errors":[
      {
        "field": "email",
        "message": "Email is already taken"
      }
    ],
    "status": code,
  }
    ```

### Users
  >- GET v1/users
  >- GET v1/users/{id}
  >- POST v1/users
  >- PUT v1/users/{id}
  >- DELETE v1/users/{id}
  #### RESPONSE
  ```json
  {
    "id": 1,
    "name": "Brent Harrington",
    "email": "bobs@gmail.com",
    "user_token": "1234567890"
  }
  ```
