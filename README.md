# DiscordBot

## 特色
1. 使用 **Go 1.25** 打造的現代 Discord 機器人
2. 部屬在 Heroku
3. 使用 Heroku 提供的免費 PostgreSQL 與 Redis
4. 支援連接池配置，提升效能和穩定性
5. 統一的 HTTP 客戶端，支援超時和重試機制
6. 結構化日誌系統 (log/slog)
7. **依賴注入與接口抽象** - 便於測試和維護
8. **統一命令路由** - 集中管理所有 Discord 命令
9. **優雅關閉機制** - 確保資源正確釋放

## 功能
1. 查詢指定標的目前股價
2. 儲存所購買股票資訊
3. Redis 儲存觀察清單，當股價大幅波動時主動通知
4. 每日自動結算當日損益與總損益
5. 顯示 ETH 即時價格

## 專案結構

```
discordBot/
├── handler/            # Discord 命令處理
├── model/              # 資料模型
│   ├── dao/            # 資料訪問對象 (Data Access Objects)
│   ├── dto/            # 資料傳輸對象 (Data Transfer Objects)
│   ├── redis/          # Redis 連接和操作
│   └── postgresql/     # PostgreSQL 連接和操作
├── pkg/                # 共用套件
│   ├── config/         # 配置管理
│   └── logger/         # 日誌系統
├── service/            # 業務邏輯服務
│   ├── client/         # HTTP 客戶端
│   ├── crypto/         # 加密貨幣相關服務
│   ├── discord/        # Discord 相關服務
│   ├── exchange/       # 匯率相關服務
│   └── stock/          # 股票相關服務
└── .github/            # GitHub 工作流配置
```

## 環境設置

### 必要環境變數

1. 複製環境變數範例文件並填入必要資訊：
```bash
cp .env.example .env
```

2. 填寫必要的環境變數：

| 變數 | 說明 | 必需 |
|------|------|------|
| `DCToken` | Discord Bot Token | 是 |
| `APIKey` | Finnhub API Key | 是 |
| `REDIS_URL` | Redis 連接 URL | 是 |
| `DATABASE_Host` | PostgreSQL 主機 | 是 |
| `DATABASE_Port` | PostgreSQL 埠 | 是 |
| `DATABASE_Name` | PostgreSQL 資料庫名稱 | 是 |
| `DATABASE_User` | PostgreSQL 用戶 | 是 |
| `DATABASE_Password` | PostgreSQL 密碼 | 是 |

### 建議配置環境變數

| 變數 | 說明 | 必需 |
|------|------|------|
| `CRYPTO_PRICE_CHANNEL_ID` | 加密貨幣價格更新頻道 ID | 是 |
| `WATCH_LIST_CHANNEL_ID` | 股票觀察清單頻道 ID | 是 |
| `PROFIT_REPORT_CHANNEL_ID` | 損益報告頻道 ID | 是 |
| `DEFAULT_USER_ID` | 預設用戶 ID（用於提及） | 是 |

### 可選環境變數

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `ENV` | 執行環境 (development/production) | development |
| `LOG_LEVEL` | 日誌級別 (DEBUG/INFO/WARN/ERROR) | INFO |

### 連接池配置環境變數

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `REDIS_POOL_SIZE` | Redis 連接池大小 | 10 |
| `REDIS_MIN_IDLE_CONNS` | Redis 最小空閒連接數 | 5 |
| `REDIS_MAX_RETRIES` | Redis 最大重試次數 | 3 |
| `REDIS_READ_TIMEOUT` | Redis 讀取超時 | 10s |
| `REDIS_WRITE_TIMEOUT` | Redis 寫入超時 | 10s |
| `DB_MAX_OPEN_CONNS` | 資料庫最大開啟連接數 | 25 |
| `DB_MAX_IDLE_CONNS` | 資料庫最大空閒連接數 | 10 |
| `DB_CONN_MAX_LIFETIME` | 資料庫連接最大生命週期 | 5m |
| `DB_CONN_MAX_IDLE_TIME` | 資料庫連接最大空閒時間 | 10m |

## 運行

```bash
# 開發模式
go run main.go

# 生產模式（需要先設置環境變數）
ENV=production go run main.go
```

## 建置

```bash
# 建置可執行檔案
go build -o discordBot

# 建置用於 Heroku 的可執行檔案
go build -o bin/discordBot
```

## 測試

專案包含多個測試文件，可以運行：

```bash
# 運行所有測試
go test ./...

# 運行測試並顯示詳細資訊
go test -v ./...

# 運行特定包的測試
go test ./service/stock
```

## 升級記錄

### 2026-02 架構優化

- **Bug 修復**:
  - 修復 SQL 參數索引錯誤（`model/dao/stock/get.go` 中兩個條件都使用 `$1` 的問題）

- **架構改進**:
  - **依賴注入**: 新增股票服務和 Discord 服務的接口抽象（`service/stock/interface.go`, `service/discord/interface.go`）
  - **命令路由**: 重構命令處理為統一路由器模式（`handler/router.go`）
  - **優雅關閉**: 實現 graceful shutdown，確保資源正確釋放
  - **配置管理**: 新增 `TaskConfig` 統一管理頻道 ID 和用戶 ID
  - **日誌整合**: 全項目統一使用 `pkg/logger`，替代 `fmt.Println`

- **新增文件**:
  - `service/stock/interface.go` - 股票服務接口
  - `service/discord/interface.go` - Discord 消息接口
  - `handler/router.go` - 命令路由器
  - `OPTIMIZATION_SUMMARY.md` - 優化詳細報告

- **環境變數更新**:
  - 新增 `CRYPTO_PRICE_CHANNEL_ID`, `WATCH_LIST_CHANNEL_ID`, `PROFIT_REPORT_CHANNEL_ID`
  - 新增 `ENV` 和 `LOG_LEVEL`
  - 新增多個 Redis 和 PostgreSQL 連接池配置選項

### 2025-02 重大重構

- **Go 版本**: 從 1.19 升級到 1.25
- **主要依賴更新**:
  - discordgo: v0.26.1 → v0.29.0
  - go-redis: v8 → v9.17.3
  - finnhub-go: v2.0.15 → v2.0.22
  - lib/pq: v1.10.7 → v1.10.9
  - godotenv: v1.4.0 → v1.5.1

- **架構改進**:
  - 新增結構化日誌系統 (使用 log/slog)
  - 新增統一的 HTTP 客戶端，支援超時和重試
  - 改進 Redis 和 PostgreSQL 連接池配置
  - 改進錯誤處理，使用錯誤包裝
  - 消除硬編碼的 Channel ID 和 User ID
  - 修復過時的 API 使用 (ioutil.ReadAll → io.ReadAll)

## 架構說明

### 依賴注入與接口

專案使用接口來抽象外部依賴，便於單元測試和維護：

```go
// 股票服務接口
service/stock/interface.go
- FinnhubClient: 股票報價 API 客戶端抽象

// Discord 服務接口
service/discord/session.go
- Session: Discord 消息發送抽象
```

### 命令路由

使用統一的命令路由器集中管理所有 Discord 命令：

```go
router := handler.NewCommandRouter()
router.Register("$+", handler.Quote)
router.Register("$set_stock", handler.SetStock)
// ... 其他命令
dg.AddHandler(router.Handle)
```

### 優雅關閉

實現 graceful shutdown 機制，確保資源正確釋放：

1. 停止定時任務
2. 關閉數據庫連線
3. 關閉 Redis 連線
4. 關閉 Discord Session

### 配置管理

統一配置管理，支援環境變量和預設值：

```go
// 獲取配置
taskConfig := config.GetTaskConfig()
logger.Init()  // 根據環境初始化日誌
```
