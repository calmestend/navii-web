# Navii Web Server

This is the web server for the Navii project, written in Go.
It uses a custom router and starts an HTTP server configurable via environment variables.

------------------------------------------------------------

## Requirements

- Go 1.25+ 
- Git

------------------------------------------------------------

## Project Structure (simplified)

```bash
.
├── cmd
│   └── web
├── internal
│   ├── handlers
│   └── router
├── main
├── resources.http
└── web
    ├── index.html
    ├── privacy_policies.html
    ├── static
    └── user_delete.html
```

------------------------------------------------------------

## Environment Variables

The server is configured using environment variables.

PORT
  Description: Port where the server runs
  Default: 8080

------------------------------------------------------------

## Getting Started 

1. Clone the repository

```bash
    git clone https://github.com/calmestend/navii-web
    cd navii-web
```

2. Install dependencies

```bash
    go mod download
```

3. Run the server

```bash
    PORT=8080 go run cmd/web/main.go
```

------------------------------------------------------------

## Testing

The project includes a `resources.http` file for testing endpoints. 

Available endpoints:
- `GET  /` - Home page
- `GET  /static/` - Static files
- `GET  /user/delete` - Delete account form
- `POST /user/delete` - Submit delete account request
- `GET  /privacy_policies` - Privacy policies page

See `resources.http` for complete request examples.

------------------------------------------------------------

## External Dependencies

Logging

[charmbracelet/log](https://github.com/charmbracelet/log) for structured logging.

------------------------------------------------------------
