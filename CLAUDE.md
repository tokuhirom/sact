# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`sact` is a TUI (Terminal User Interface) application for managing Sakura Cloud servers. The name comes from さくっと (sakutto - "quickly/easily"). The primary use case is operational tasks like starting/stopping servers, with the assumption that server provisioning is handled by Terraform.

## Git Workflow

- The `main` branch is protected. Never commit directly to main.
- Always create a feature branch and submit a pull request.
- Before committing, always run `go fmt ./...` to format the code.

## Architecture

### Application Structure

The application uses the Bubble Tea framework (bubbletea) for the TUI, following the Elm architecture pattern:

- **Model** (`internal/model.go`): Contains the application state including current zone, server list, cursor position, loading states
- **Update** (`internal/model.go:81`): Handles messages (key presses, async data loads) and updates the model
- **View** (`internal/model.go:122`): Renders the UI based on current model state
- **Commands**: Async operations like `loadServers` that fetch data from Sakura Cloud API

### Key Components

- `cmd/sact/main.go`: Entry point that handles logging setup, config loading, and TUI initialization
- `internal/client.go`: Wrapper around `github.com/sacloud/iaas-api-go` for Sakura Cloud API operations
- `internal/model.go`: Bubble Tea TUI implementation
- `internal/config.go`: Configuration loading from `~/.config/sact/config.toml`

### Authentication

Credentials are loaded from environment variables via the sacloud API client:
- `SAKURA_ACCESS_TOKEN`
- `SAKURA_ACCESS_TOKEN_SECRET`

These must be set before running the application.

### Supported Zones

The application supports 5 Sakura Cloud zones: `tk1a`, `tk1b`, `is1a`, `is1b`, `is1c`

## Development Commands

### Build and Run

```bash
# Build the binary
make build

# Build and run
make run

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

### Testing and Code Quality

```bash
# Run all tests
make test

# Format code
make fmt

# Run go vet
make vet

# Full pre-commit check (fmt + vet + build)
make all
```

### Running the Application

```bash
# Run with default settings (logs to stderr)
./sact

# Run with file logging for debugging
./sact --log=/tmp/sact.log
```

## Code Patterns

### Bubbletea Message Flow

When adding new async operations:
1. Define a message type (e.g., `serversLoadedMsg`)
2. Create a command function that returns `tea.Cmd`
3. Handle the message in `Update()` method
4. Update model state and return new commands if needed

Example pattern from `internal/model.go:47-53`:
```go
type myDataLoadedMsg struct {
    data []MyData
    err  error
}

func loadMyData(client *SakuraClient) tea.Cmd {
    return func() tea.Msg {
        // async operation
        return myDataLoadedMsg{data: result, err: err}
    }
}
```

### Zone Switching

When switching zones (`internal/model.go:90-99`):
1. Update cursor position
2. Update current zone in model
3. Call `client.SetZone()`
4. Set loading state
5. Return `loadServers` command

### Logging

Use `log/slog` for structured logging. Log files are configured via `--log` flag. Important events to log:
- API operations (zone switches, server fetches)
- User interactions (key presses for zone switching, refresh, quit)
- Errors

## Testing

Tests use `github.com/stretchr/testify` for assertions. See `internal/client_test.go` for examples.

## Future Extensions

Per README.md, after server list functionality is solid, the plan is to add support for:
- switch+router management
- switch management
- DNS management
- Database appliances
