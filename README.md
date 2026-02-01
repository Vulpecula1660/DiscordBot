# DiscordBot

## 特色
1. 使用 **Go 1.25** 打造的現代 Discord 機器人
2. 部屬在 Heroku
3. 使用 Heroku 提供的免費 PostgreSQL 與 Redis
4. 支援連接池配置，提升效能和穩定性
5. 統一的 HTTP 客戶端，支援超時和重試機制
6. 結構化日誌系統

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
| `CHANNEL_PROFIT_REPORT` | 損益報告頻道 ID | 建議 |
| `CHANNEL_WATCH_LIST` | 觀察清單頻道 ID | 建議 |
| `CHANNEL_CRYPTO_UPDATE` | 加密貨幣更新頻道 ID | 建議 |
| `DEFAULT_USER_ID` | 預設用戶 ID | 建議 |

### 選填環境變數（連接池配置）

- `REDIS_POOL_SIZE` - Redis 連接池大小 (預設: 10)
- `REDIS_MIN_IDLE_CONNS` - Redis 最小空閒連接數 (預設: 5)
- `REDIS_MAX_RETRIES` - Redis 最大重試次數 (預設: 3)
- `DB_MAX_OPEN_CONNS` - 資料庫最大開啟連接數 (預設: 25)
- `DB_MAX_IDLE_CONNS` - 資料庫最大空閒連接數 (預設: 10)
- `LOG_LEVEL` - 日誌級別: DEBUG, INFO, WARN, ERROR (預設: INFO)

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
