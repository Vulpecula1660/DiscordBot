# AGENTS.md

## Build Commands

```bash
# Build the project
go build -o bin/discordBot

# Build for Heroku (production)
go build -o bin/discordBot -ldflags="-s -w"

# Run the application (requires .env or env vars set)
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

# Run a single test by name
go test -run Test_Quote ./service/stock

# Run a single sub-test (table-driven)
go test -run Test_Quote/valid_stock_quote ./service/stock

# Run tests for a specific package
go test ./service/stock
go test ./service/discord
go test ./service/exchange

# Run tests with race detector (also used in CI)
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Lint (CI uses default linters, no .golangci.yml)
golangci-lint run --timeout=5m
```

CI runs with `APIKey=""` and `DCToken=""` -- tests must not depend on real credentials.

## Code Style Guidelines

### Language
- Comments, log messages, and business error messages use **Traditional Chinese**.
- Infrastructure error wrapping uses English (`"failed to get connection: %w"`).
- Log key names use English (`logger.Error("查詢失敗", "symbol", sym, "error", err)`).

### Imports
- Three groups separated by blank lines: **stdlib**, **third-party**, **project internal**.
- Use blank import for side effects: `_ "github.com/lib/pq"`.
- Alias packages only to avoid conflicts (`stockdao "discordBot/model/dao/stock"`).

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
- **Module**: `discordBot` (Go 1.25)
- **Packages**: lowercase, single word (`handler`, `stock`, `redis`)
- **Types**: PascalCase (`SendMessageInput`, `TaskConfig`)
- **Functions**: PascalCase exported, camelCase unexported
- **Variables**: camelCase (`channelID`, `userID`)
- **Constants**: camelCase or PascalCase (`maxRetries`, `coingeckoAPIURL`)
- **Interfaces**: PascalCase; defined near the consumer, not the implementer

### Function Comments
- Exported functions: comment must start with function name.
- Both `// FuncName : description` (older) and `// FuncName description` (newer) are accepted.

```go
// Quote : 查詢標的
func Quote(ctx context.Context, message string) (string, error) {
```

### Error Handling
- Always check errors; return early to reduce nesting.
- Wrap with context using `fmt.Errorf("...: %w", err)` (use `%w`, not `%v`).
- Business errors in Chinese; infrastructure errors in English.

```go
conn, err := GetConn("Redis")
if err != nil {
	return fmt.Errorf("failed to get redis connection: %w", err)
}

if err := stock.Ins(ctx, nil, input); err != nil {
	return err
}
```

### Context Usage
- Pass `context.Context` as first parameter in service/model functions.
- Handler functions create `context.Background()` at entry points.
- Use `context.WithTimeout` for all external calls.

### Logging
- Use `pkg/logger` (wraps `log/slog`), never `fmt.Println`.
- Messages in Chinese, structured key-value pairs with English keys.

```go
logger.Info("查詢股票價格", "symbol", symbol)
logger.Error("計算損益失敗", "symbol", sym, "error", err)
```

### Testing
- Test files: `*_test.go` in the same package.
- Test functions: `Test_FunctionName` with underscore (project convention).
- Use **table-driven tests** with `t.Run`.
- Assertions use `t.Errorf` directly (no assertion library).
- Float comparisons use epsilon tolerance.
- Mocks: manual `MockXxx` structs in `mock_test.go`, implementing interfaces.
- Setup/teardown: `SetDefaultClient(mock)` / `defer ResetDefaultClient()`.
- Tests must pass without real API keys or Discord tokens.

```go
tests := []struct {
	name    string
	want    string
	wantErr bool
}{...}
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		if got != tt.want {
			t.Errorf("Quote() = %v, want %v", got, tt.want)
		}
	})
}
```

### Dependency Injection
- Define interfaces near the consumer, not the implementation.
- Compile-time check: `var _ Session = (*discordgo.Session)(nil)`.
- For testability, provide `FuncWithDeps(deps Interface)` variants or global Set/Reset.

### Project Structure
```
handler/         # Discord command handlers (router dispatches here)
model/
  dao/           # Data Access Objects (SQL queries)
  dto/           # Data Transfer Objects
  postgresql/    # PostgreSQL connection pool
  redis/         # Redis connection pool + operations
pkg/
  config/        # Environment variable configuration (GetTaskConfig, etc.)
  logger/        # Structured logging (slog wrapper)
service/
  client/        # HTTP client with retry + functional options
  crypto/        # CoinGecko ETH price service
  discord/       # Session interface, message sending, error reporter
  exchange/      # Currency exchange rate service
  stock/         # Stock quote, profit calculation, watch list alerts
```

## General Rules
- Go version: 1.25; deploy target: Heroku worker dyno
- No naked returns
- No hardcoded channel IDs or user IDs -- use `pkg/config`
- Prefer `strings.Fields()` over `strings.Split(s, " ")` for command parsing
- Command routing uses longest-prefix-first matching (`handler/router.go`)
- Bot self-message filtering is in the router only; handlers must not duplicate it
- Use `sync.RWMutex` for connection pools; never hold locks during network I/O
- Concurrency limiting uses channel-based semaphore + `sync.WaitGroup`
- Use `defer rows.Close()` for SQL query results
