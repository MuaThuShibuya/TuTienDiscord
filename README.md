# Vạn Pháp Tiên Nghịch — Discord Bot Tu Tiên RPG

Bot game tu tiên dạng RPG all-in-app trên Discord, viết bằng Go.

## Yêu cầu

- Go 1.22+
- MongoDB Atlas account (free tier đủ dùng)
- Discord Application + Bot token
- (Optional) Docker

## Chạy local

### 1. Clone và cài dependencies

```sh
git clone <repo-url>
cd tu-tien-bot
go mod download
```

### 2. Tạo file .env

```sh
cp .env.example .env
```

Điền các giá trị vào `.env`:

```env
DISCORD_TOKEN=your_bot_token_here
DISCORD_APP_ID=your_app_id_here
DISCORD_GUILD_ID=your_test_server_id
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/
```

### 3. Chạy bot

```sh
go run ./cmd/bot
```

Bot sẽ:
- Kết nối MongoDB Atlas
- Tạo indexes
- Đăng ký slash commands (guild mode — tức thì)
- Bật HTTP server tại `http://localhost:8080`
- Sẵn sàng nhận lệnh Discord

### 4. Test health check

```sh
curl http://localhost:8080/health
```

Kết quả mong đợi:
```json
{"status":"ok","app":"tu-tien-discord-bot","version":"0.1.0","database":"connected"}
```

## Deploy lên Render

### 1. Tạo Web Service trên Render

- **Build Command**: `go build -o bot ./cmd/bot`
- **Start Command**: `./bot`
- **Instance Type**: Free (hoặc Starter nếu cần)

### 2. Thiết lập Environment Variables

Trên Render Dashboard → Environment, thêm tất cả biến từ `.env.example`:

| Biến | Giá trị |
|---|---|
| `DISCORD_TOKEN` | Bot token từ Discord Developer Portal |
| `DISCORD_APP_ID` | Application ID |
| `DISCORD_GUILD_ID` | (để trống cho global commands) |
| `MONGODB_URI` | MongoDB Atlas connection string |
| `MONGODB_DATABASE` | `tu_tien_bot` |
| `COMMAND_REGISTER_MODE` | `global` (cho production) |
| `KEEPALIVE_ENABLED` | `true` |
| `KEEPALIVE_URL` | `https://your-app-name.onrender.com/health` |
| `LOG_LEVEL` | `info` |

### 3. Health Check URL

Render → Settings → Health Check Path: `/health`

### 4. Keepalive

Bot tự ping `/health` mỗi 10 phút để tránh Render free tier sleep.
Nếu dùng Render paid plan hoặc VPS, set `KEEPALIVE_ENABLED=false`.

## Deploy lên VPS (Docker)

```sh
# Build và chạy
sh scripts/docker_run.sh
```

Hoặc với docker-compose (tạo sau):

```sh
docker-compose up -d
```

Sau khi lên VPS, tắt keepalive (`KEEPALIVE_ENABLED=false`) vì systemd/Docker
restart policy đảm bảo process luôn chạy.

## Cấu trúc dự án

Xem chi tiết: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

```
cmd/bot/main.go          ← Entrypoint
internal/
  config/                ← Load env vars
  logger/                ← Structured logging
  errors/                ← Sentinel errors
  database/              ← MongoDB connection + indexes
  game/                  ← Game domain (profile, cultivation, economy, cooldown)
    */model.go           ← Data structs
    */repository.go      ← DB access (interface + MongoDB impl)
    */service.go         ← Business logic
  discord/
    bot.go               ← Bot lifecycle
    commands.go          ← Slash command definitions
    router.go            ← Top-level interaction dispatcher
    handlers/            ← Controllers (nhận Discord interaction)
    menu/                ← Menu session + UI builders + menu router
    ui/                  ← Reusable UI components
  server/                ← HTTP health check
  scheduler/             ← Keepalive, session cleanup
pkg/utils/               ← Time, ID, number helpers
```

## Roadmap

Xem [docs/VERSION_PLAN.md](docs/VERSION_PLAN.md) để biết kế hoạch từng version.

## Bảo mật

Xem [docs/SECURITY.md](docs/SECURITY.md) để biết các quy tắc bảo mật.

## MongoDB Schema Migration

- MongoDB không ép schema cứng.
- Nếu đổi field từ `string` sang `int` (ví dụ: `mindState` ở v0.2), dữ liệu cũ có thể decode lỗi trong Go (`cannot decode string into an integer type`).
- Bắt buộc phải chạy script migration (như `scripts/mongo_fix_cultivation_schema.js`) khi đổi model.
- Runtime không được dùng dữ liệu giả (fake data chỉ dùng trong mock testing).
- Database cũ/cùng cluster phải tách `MONGODB_DATABASE` rõ ràng để tránh xung đột schema.

---

> Dự án game tu tiên thuần giải trí. Không có giao dịch tiền thật. Gacha chỉ dùng tài nguyên trong game.
