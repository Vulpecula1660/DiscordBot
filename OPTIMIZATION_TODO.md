# 專案優化代辦清單

## Phase 1 - 穩定性與可靠性（高優先）

### 目標
先解決容易中斷服務、造成排程失效、或外部 API 被打爆的問題。

### Tasks
- [x] 修正首次執行 `CalculateProfit` 時，`yesterdayTotalValue` 為空字串導致 `ParseFloat` 失敗（`service/stock/cauculate_profit.go`）。
- [x] 為所有 `taskCron.AddFunc(...)` 加上錯誤處理與記錄，避免 cron 註冊失敗被忽略（`handler/task.go`）。
- [x] 在 `CheckChange` 與 `CalculateProfit` 加入併發上限（worker pool 或 semaphore），避免大量標的時觸發 API rate limit（`service/stock/check_change.go`, `service/stock/cauculate_profit.go`）。
- [x] 為排程內外部呼叫補上可控 timeout context（目前多處使用 `context.Background()`）。
- [x] 針對排程錯誤回報訊息做節流/去重，避免錯誤風暴時大量重複通知。

### Done 條件
- [ ] 空 Redis 狀態下可正常跑完一次收益計算。
- [ ] cron 註冊失敗可在 logs 明確看到。
- [ ] watch list 大量資料時，系統可穩定運行且不會瞬間打滿外部 API。

---

## Phase 2 - 測試穩定與可持續 CI（高優先）

### 目標
讓測試可離線、可重複、可在 CI 穩定通過。

### Tasks
- [ ] 將會打真實外部服務的測試改為 mock/stub（`service/stock/stock_test.go`, `service/exchange/convert_test.go`）。
- [ ] 將會發送真實 Discord 訊息的測試改為 interface mock（`service/discord/send_message_test.go`）。
- [ ] 移除對即時價格的固定值斷言，改為可控 fixture。
- [ ] 修正 table-driven test 錯誤斷言模板（如 `if err != nil && tt.wantErr`），避免漏判。
- [ ] 補齊關鍵邊界測試：空 watch list、空持倉、Redis miss、匯率 API 異常、訊息發送失敗。
- [ ] 新增 CI workflow（至少包含 `go test ./...`，建議加 `go test -race ./...`）。

### Done 條件
- [ ] 本機與 CI 在無 API Key/Discord Token 下仍可穩定執行單元測試。
- [ ] 測試重跑多次結果一致，不依賴外部即時資料。

---

## Phase 3 - 架構收斂與可維護性（中優先）

### 目標
降低重複邏輯與未使用抽象，讓後續功能開發更快、更安全。

### Tasks
- [x] 將命令參數解析由 `strings.Split(..., " ")` 改為 `strings.Fields(...)`，提升容錯（`handler/redis_string.go`, `handler/redis_list.go`, `service/stock/set_stock.go`）。
- [x] 讓命令路由匹配順序可預期（註冊順序或最長前綴優先），避免 map 迭代不穩定（`handler/router.go`）。
- [x] 收斂 bot 自己訊息過濾邏輯，避免 Router 與 handler 重複判斷（`handler/router.go`, `handler/stock.go`, `handler/redis_*.go`）。
- [x] 統一 Discord 訊息發送抽象，整合 `service/discord/session.go` 與 `service/discord/send_message.go`。
- [x] 統一 config key 命名並規劃舊 key 遷移（`CHANNEL_*` vs `*_CHANNEL_ID`），移除預設硬編碼頻道 ID。
- [x] 清理未使用的 config 型別/方法與介面，降低死碼（`pkg/config/config.go`, `service/stock/interface.go`）。
- [x] 優化 DB/Redis 連線池鎖粒度，避免在 lock 內做耗時操作（`model/postgresql/create_conn.go`, `model/redis/create_conn.go`）。

### Done 條件
- [x] 命令解析對多空白、尾隨空白有一致行為。
- [x] 路由結果在每次啟動都一致可預期。
- [x] 設定檔與 env 命名規則一致，無重複來源。

---

## Phase 4 - 文件與工程治理（中低優先）

### 目標
補足長期維運與交接可讀性。

### Tasks
- [ ] 更新 README：新增優化後的命令/設定說明與 migration note。
- [ ] 增加故障排查指南（排程失敗、Redis/DB 連線失敗、API 額度問題）。
- [ ] 將 GitHub Actions 舊版 action（如 CodeQL `v1`）升級到目前建議版本。
- [ ] 補上「部署前檢查清單」：環境變數、資料庫連線、Redis、Discord 權限。

### Done 條件
- [ ] 新成員可依文件在 30 分鐘內完成本機啟動與基本驗證。
- [ ] 發版前檢查有明確標準化流程可執行。
