# Photo Slider

[English](README.md)

一個輕量級、開箱即用的本地照片瀏覽器。只需指定本地照片目錄，即可在瀏覽器中瀏覽幻燈片 — 無需資料庫，無需雲端服務，單一執行檔即可運行。

<p align="center">
  <img src="https://img.shields.io/github/v/release/Aquila-f/photo-slider" alt="Release">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="License">
</p>

## 功能特色

- 掃描多個目錄中的照片（JPEG、PNG、WebP、GIF）
- 依照資料夾結構自動組織相簿
- 即時圖片壓縮（最大 1920 px，JPEG 品質 80）並提供記憶體 LRU 快取
- 顯示 EXIF 中繼資料（相機型號、拍攝日期）
- 支援鍵盤、滑鼠及觸控/滑動操作
- 全螢幕模式，附帶資訊疊加層
- 隨機播放與自動播放，可調整間隔時間（1–30 秒）
- 透過 Web 介面管理照片來源目錄
- 單一執行檔，內嵌 Web 資源 — 無外部相依性
- 提供 REST API 供程式化存取

## 快速開始

### 從原始碼編譯

```bash
go build -o photo-slider ./cmd
cp config.example.yaml config.yaml   # 編輯並填入你的照片路徑
./photo-slider
```

### 使用預先編譯版本

從 [Releases](https://github.com/Aquila-f/photo-slider/releases) 頁面下載對應平台的執行檔，編輯 `config.example.yaml` 後執行：

```bash
./photo-slider -config config.yaml
```

然後開啟 **http://localhost:8080**。

### 命令列參數

| 參數 | 預設值 | 說明 |
|------|--------|------|
| `-config` | `config.yaml` | 設定檔路徑 |
| `-port` | `8080` | HTTP 伺服器連接埠 |

## 設定

建立 `config.yaml`（參考 [config.example.yaml](config.example.yaml)）：

```yaml
sources:
  - /path/to/your/photos
  - /another/photo/directory
```

每個項目必須是已存在的可讀目錄。相對路徑在啟動時會自動解析為絕對路徑。

你也可以在執行期間透過 Web 介面新增或移除照片來源 — 點選頁面頂部的 **Sources** 面板即可操作。

## 操控方式

### 鍵盤快捷鍵

| 按鍵 | 功能 |
|------|------|
| `←` / `→` | 上一張 / 下一張 |
| `Space` | 切換自動播放 |
| `f` | 切換全螢幕 |
| `Esc` | 離開全螢幕或停止播放 |

### 觸控操作

左右滑動切換照片。

### 幻燈片播放

點選 ▶ 開始自動播放。使用間隔滑桿（1–30 秒）調整速度。勾選 **Shuffle** 啟用隨機播放。

## API 介面

| 方法 | 路徑 | 說明 |
|------|------|------|
| `GET` | `/api/sources` | 列出已設定的照片來源目錄 |
| `POST` | `/api/sources` | 新增照片來源目錄 |
| `DELETE` | `/api/sources` | 移除照片來源目錄 |
| `GET` | `/api/albums` | 列出所有相簿 |
| `GET` | `/api/albums/:key` | 列出相簿中的照片（`?shuffle=true` 啟用隨機排序） |
| `GET` | `/photos/:album/:key` | 取得壓縮後的照片 |

相簿和照片識別碼使用 Base64 URL 編碼。當 EXIF 資料可用時，照片回應會包含 `X-Photo-Taken-At`（RFC 3339 格式）和 `X-Photo-Model` 回應標頭。

## 專案結構

```
cmd/
  main.go             進入點，相依性組裝
  static/             內嵌 Web 介面（HTML、JS、CSS，使用 Alpine.js）

internal/
  config/             YAML 設定載入器
  domain/             核心型別、介面、錯誤定義
  handler/            Gin HTTP 處理器與路由
  mapper/             Base64 編碼/解碼器
  photo/              圖片壓縮器、環形緩衝 LRU 快取、EXIF 擷取器
  service/            業務邏輯（相簿同步、來源目錄管理）
  storage/            本地檔案系統提供器
  strategy/           相簿產生策略與照片清單策略
```

## 建置

```bash
# 開發建置
go build -o photo-slider ./cmd

# 執行測試
go test ./...

# 跨平台發佈（需要 GoReleaser）
goreleaser release --snapshot --clean
```

## 授權條款

MIT
