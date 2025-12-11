# Twist CLI

A command-line interface for interacting with the [Twist](https://twist.com) API. Manage your workspaces, channels, and conversations from your terminal.

## Features

- Personal access token authentication
- List workspaces
- Secure token storage in `~/.config/twist/config.json`

## Installation

### From Source

```bash
go install github.com/robertmeta/twist-cli@latest
```

### Build Locally

```bash
git clone https://github.com/robertmeta/twist-cli.git
cd twist-cli
go build
```

## Getting Your Twist API Token

1. Go to [https://twist.com/integrations](https://twist.com/integrations)
2. Create a new integration or select an existing one
3. Copy your personal access token from the OAuth section
4. Run any `twist` command and you'll be prompted to enter your token

Your token will be securely stored in `~/.config/twist/config.json` for future use.

## Usage

### List Workspaces

View all workspaces you have access to:

```bash
twist workspaces list
```

Output:
```
ID      NAME              PLAN
--      ----              ----
12345   My Team           unlimited
67890   Personal Space    free
```

### Help

Get help on available commands:

```bash
twist --help
twist workspaces --help
```

### Version

Check the CLI version:

```bash
twist --version
```

## Configuration

The CLI stores your configuration in `~/.config/twist/config.json`:

```json
{
  "token": "your-api-token-here"
}
```

You can manually edit this file if needed, or delete it to reset your authentication.

## Project Structure

```
twist-cli/
├── cmd/              # Cobra command definitions
├── pkg/
│   ├── api/         # Twist API client
│   └── config/      # Configuration management
└── internal/
    └── auth/        # Authentication helpers
```

## Development

### Prerequisites

- Go 1.21 or higher

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o twist
```

## Contributing

Contributions are welcome! This project follows standard Go conventions to make it easy for developers to contribute.

### Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following Go best practices
4. Write tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style

- Follow standard Go formatting (`gofmt`, `go vet`)
- Write clear, descriptive commit messages
- Add comments for exported functions and types
- Keep functions focused and modular

## API Documentation

This CLI uses the [Twist API v3](https://developer.twist.com/v3/). For more information about available endpoints and data structures, see the official documentation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Roadmap

Future enhancements may include:

- Channel management commands
- Conversation and message operations
- Thread management
- User and team management
- Webhooks configuration
- OAuth2 flow support
- Shell completion

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses the [Twist API](https://developer.twist.com/v3/)
