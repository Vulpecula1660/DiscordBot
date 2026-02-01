# AGENTS.md

## Build Commands

```bash
# Build the project
go build -o bin/discordBot

# Build for Heroku (production)
go build -o bin/discordBot -ldflags="-s -w"

# Run the application
go run main.go

# Clean build artifacts
rm -rf bin/
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -run Test_Quote ./service/stock

# Run tests for a specific package
go test ./service/stock
go test ./model/dao/stock

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

## Code Style Guidelines

### Imports
- Group imports: stdlib, third-party, project internal
- Use blank import for side effects (e.g., `_ "github.com/lib/pq"`)
- Import project packages as: `"discordBot/service/stock"`

```go
import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordBot/service/discord"
	"discordBot/service/stock"
)
```

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `handler`, `stock`, `redis`)
- **Types**: PascalCase (e.g., `SendMessageInput`, `StockConfig`)
- **Functions**: PascalCase for exported, camelCase for internal
- **Variables**: camelCase (e.g., `channelID`, `userID`)
- **Constants**: PascalCase or camelCase (e.g., `maxRetries`, `coingeckoAPIURL`)
- **Interfaces**: PascalCase with -er suffix (e.g., `Reader`, `Writer`)

### Error Handling
- Always check errors and wrap with context using `fmt.Errorf("...: %w", err)`
- Return early on errors to reduce nesting
- Use `if err := fn(); err != nil { return err }` pattern for single-line checks

```go
conn, err := GetConn("Redis")
if err != nil {
	return fmt.Errorf("failed to get redis connection: %w", err)
}

// Or for single-line:
if err := stock.Ins(ctx, nil, input); err != nil {
	return err
}
```

### Function Comments
- All exported functions must have a comment starting with function name
- Use `// FunctionName : description` format (colon style seen in codebase)
- Keep comments concise but descriptive

```go
// Quote : 查詢標的
func Quote(ctx context.Context, message string) (string, error) {
```

### Project Structure
```
handler/         # Discord command handlers
model/           # Data models
  dao/           # Data Access Objects
  dto/           # Data Transfer Objects
  postgresql/    # PostgreSQL connection
  redis/         # Redis connection
pkg/             # Shared packages
  config/        # Configuration
  logger/        # Structured logging
service/         # Business logic
  client/        # HTTP client utilities
  crypto/        # Crypto price service
  discord/       # Discord service
  exchange/      # Exchange rate service
  stock/         # Stock service
```

### Testing
- Test files: `*_test.go` suffix in same package
- Test functions: `Test_FunctionName` (e.g., `Test_Quote`)
- Use table-driven tests where applicable
- Mock external dependencies; don't call real APIs in tests

### Environment Variables
- Use `pkg/config` for all configuration access
- Environment variables use UPPER_CASE (e.g., `DCToken`, `DATABASE_Host`)
- Always provide defaults via `getEnv()` helper functions

### Context Usage
- Pass `context.Context` as first parameter to functions
- Use `context.Background()` at entry points
- Respect context cancellation in HTTP requests and DB operations

### Database/Redis
- Always use connection pool configuration
- Handle connection errors gracefully
- Use `defer rows.Close()` for query results

## General Rules
- Go version: 1.25
- No naked returns
- Avoid global state where possible
- Prefer composition over inheritance
- Use `log/slog` for structured logging (via `pkg/logger`)
