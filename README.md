# omnitrix.sh

a claude code alternative with no vendor lock-in, think of it as cursor for terminal.

## What is this?

omnitrix.sh is a CLI tool that lets you chat with AI about your code directly in the terminal. It supports multiple LLM providers (OpenAI, Anthropic, Google) and integrates with Language Server Protocol (LSP) to understand your codebase better.

## Features

- **Multi-provider support**: Works with OpenAI, Anthropic Claude, and Google's Gemini
- **LSP integration**: Understands your code through LSP servers (like gopls)
- **Terminal UI**: Clean, mouse-enabled interface built with Bubble Tea
- **Session management**: Keeps track of your conversations with the AI
- **File watching**: Monitors workspace changes in real-time
- **MCP tools support**: Extensible with Model Context Protocol tools
- **Database-backed**: Stores sessions and messages locally in SQLite

## Getting started

### Prerequisites

- Go 1.24 or later
- An API key for your preferred LLM provider (OpenAI, Anthropic, or Google)

### Installation

```bash
go build -o omnitrix .
```

### Configuration

Create a `.env` file in the project directory:

```bash
OPENAI_API_KEY=your-key-here
# or ANTHROPIC_API_KEY=your-key-here
# or GOOGLE_API_KEY=your-key-here
```

You can also customize the tool with a `.omnitrix.json` file:

```json
{
  "model": {
    "coder": "claude-3.7-sonnet",
    "coderMaxTokens": 20000
  },
  "lsp": {
    "gopls": {
      "command": "gopls"
    }
  }
}
```

### Usage

Run the tool in your project directory:

```bash
./omnitrix
```

Or with debug logging:

```bash
./omnitrix --debug
```

The TUI will launch, and you can start chatting with the AI about your code.

## Development

### Build

```bash
go build ./...
```

### Lint

```bash
go fmt ./...
go vet ./...
```

### Testing

```bash
go test ./...
```

## Project structure

- `cmd/` - CLI entry point and command setup
- `internals/app/` - Core application logic and service orchestration
- `internals/config/` - Configuration management
- `internals/database/` - SQLite database layer
- `internals/llm/` - LLM provider integrations and agents
- `internals/lsp/` - Language Server Protocol client
- `internals/tui/` - Terminal UI components
- `internals/sessions/` - Session management
- `internals/message/` - Message handling
- `internals/permissions/` - Permission system
- `internals/logging/` - Logging infrastructure

## License

Check the repository for license information.

## Contributing

Contributions are welcome. Please follow the Go code style guidelines and run the linters before submitting.
