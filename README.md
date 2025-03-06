# DiscordBot

## 特色
1. 使用 Golang 打造的一個 Discord 機器人
2. 部屬在 Heroku
3. 使用 Heroku 提供的免費 PostgreSQL 與 Redis

## 功能
1. 查詢指定標的目前股價
2. 儲存所購買股票資訊
3. Redis 儲存觀察清單，當股價大幅波動時主動通知
4. 每日自動結算當日損益與總損益

## 專案結構

```
discordBot/
├── handler/            # Discord 命令處理
├── model/              # 資料模型
│   ├── dao/            # 資料訪問對象
│   ├── dto/            # 資料傳輸對象
│   ├── redis/          # Redis 連接和操作
│   └── postgresql/     # PostgreSQL 連接和操作
├── service/            # 業務邏輯服務
│   ├── discord/        # Discord 相關服務
│   ├── stock/          # 股票相關服務
│   └── exchange/       # 匯率相關服務
└── .github/            # GitHub 工作流配置
```

## 環境設置

1. 複製環境變數範例文件並填入必要資訊：
```
cp .env.example .env
```

2. 填寫必要的環境變數：
   - Discord Bot Token
   - Redis 連接資訊
   - PostgreSQL 連接資訊
   - Finnhub API Key

## 運行

```bash
go run main.go
```

## 測試

專案包含多個測試文件，可以運行：

```bash
go test ./...
```
