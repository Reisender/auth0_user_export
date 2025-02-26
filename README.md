# Auth0 User Export CLI Tool

A command-line tool to export Auth0 users to CSV format, including their app_metadata. The tool outputs to STDOUT, allowing for easy redirection to files or piping to other commands.

## Features

- Export Auth0 users to CSV format
- Include app_metadata in the export
- Customize which user fields to include
- Paginated requests to handle large user bases
- Environment variable support for authentication credentials

## Installation

```
go install github.com/auth0/auth0-user-export@latest
```

### Prerequisites

- Go 1.13 or later
- Auth0 account with Management API credentials

### Building from source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/auth0-user-export.git
   cd auth0-user-export
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o auth0-user-export
   ```

## Configuration

The tool requires Auth0 Management API credentials, which can be provided via command-line flags or environment variables:

- Auth0 Domain
- Client ID
- Client Secret

To set up these credentials:

1. Log in to your Auth0 dashboard
2. Navigate to Applications > APIs > Auth0 Management API
3. Authorize a client application or create a new one
4. Ensure the client has the `read:users` permission

### Basic Usage 

With env vars set:
```bash
export AUTH0_DOMAIN=<your-auth0-domain>
export AUTH0_CLIENT_ID=<your-client-id>
export AUTH0_CLIENT_SECRET=<your-client-secret>
auth0-user-export > ./users.csv
```

Using flags:
```bash
auth0-user-export --domain <your-auth0-domain> --client-id <your-client-id> --client-secret <your-client-secret> > ./users.csv
```

