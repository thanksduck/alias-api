# Authentication API Endpoints

## Endpoints

---

### 1. Register a New User

**Endpoint:**

- `POST /api/v2/auth/signup`
- `POST /api/v2/auth/register`

**Description:**

Registers a new user account with the provided information.

**Request Body:**

```json
{
  "username": "string",        // Required, unique username (min 4 characters)
  "name": "string",            // Required, full name (min 4 characters)
  "email": "string",           // Required, valid email address
  "password": "string",        // Required, min 8 characters
  "passwordConfirm": "string", // Required, should match password
  "avatar": "string"           // Optional, URL to avatar image
}
```

### Response

- Success (201 Created)

```json
{
  "message": "User Created Successfully",
  "status": "success",
  "user": {
    "username": "string",
    "name": "string",
    "email": "string",
    "avatar": "string"
  }
}
```

**Error Responses:**

- `400 Bad Request` if validation fails.
- `409 Conflict` if the user already exists.
- `500 Internal Server Error` for server errors.

### 2. User Login

**Endpoint:**

- `POST /api/v2/auth/login`

**Description:**

Authenticates a user using their username or email and password.

**Request Body:**

```json
{
    "username": "string", // Optional if email is provided
    "email": "string",    // Optional if username is provided
    "password": "string"  // Required
}
```

**Note: Either `username` or `email` must be provided.**

**Response:**

- Success (200 OK)

```json
{
    "message": "Login Successful",
    "status": "success",
    "user": {
        "username": "string",
        "name": "string",
        "email": "string",
        "avatar": "string"
    }
}
```

**Error Responses:**

- `400 Bad Request` if required fields are missing.
- `401 Unauthorized` if credentials are invalid.
- `500 Internal Server Error` for server errors.
