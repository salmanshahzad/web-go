# Go Web Server

-   Web server written in Go using Fiber.
-   Connects to PostgreSQL and Redis.
-   Includes middleware for CORS and cookie-based sessions.
-   Serves static files from the `public` directory.

## Usage

Ensure environment variables are set in `.env`.
See required environment variables in `.env.example`.

### Local

```bash
make run
```

### Docker

```bash
docker build -t web-go .
docker run [-d] [-p] [-v] web-go
```

## Endpoints

-   `GET /api/health` returns `204`
-   `GET /api/user` return signed in user
-   `POST /api/user` create new user
-   `PUT /api/user/username` change username
-   `PUT /api/user/password` change password
-   `DELETE /api/user` delete signed in user
-   `POST /api/session` sign in
-   `DELETE /api/session` sign out
