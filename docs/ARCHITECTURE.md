# Architecture — Vạn Pháp Tiên Nghịch Discord Bot

## Tổng quan tầng (Strict Layer Separation)

Bot được tổ chức theo 5 tầng nghiêm ngặt. **Không tầng nào được bỏ qua tầng khác.**

```
Discord Interaction
       │
       ▼
┌──────────────────────────────────┐
│  discord/router.go               │  Top-level Router
│  Phân loại: slash command /       │  Nhận interaction từ Discord,
│  component interaction / modal    │  dispatch đến đúng Handler
└──────────────┬───────────────────┘
               │
       ┌───────▼────────┐
       │    Handler      │  discord/handlers/*.go
       │   (Controller)  │  Nhận interaction → validate input
       │                 │  → gọi Service → gọi UI Builder
       │                 │  → trả về response cho Discord
       └───┬─────────┬───┘
           │         │
    ┌──────▼──┐  ┌───▼──────────────┐
    │ Service  │  │   UI Builder     │
    │          │  │                  │
    │ Business │  │ discord/ui/      │
    │ logic    │  │ discord/menu/    │
    │ Gọi Repo │  │ *_menu.go        │
    └──────┬───┘  │ Chỉ build embed  │
           │      │ và component     │
    ┌──────▼───┐  └──────────────────┘
    │ Repository│
    │           │
    │ DB access │
    │ only      │
    │ MongoDB   │
    └───────────┘
```

## Chi tiết từng tầng

### 1. Repository (`internal/game/*/repository.go`)
- **Trách nhiệm**: Chỉ nói chuyện với MongoDB
- **Không được**: Chứa business logic, gọi service khác
- **Có interface**: Dễ mock khi test
- **Ví dụ**: `profile.NewRepository(db)`, `economy.NewRepository(db)`

### 2. Service (`internal/game/*/service.go`)
- **Trách nhiệm**: Business logic thuần túy
- **Không được**: Gọi Discord API, build embed, nhận discordgo objects
- **Phụ thuộc**: Chỉ vào Repository interface
- **Ví dụ**: `profile.NewService(repo)`, `economy.NewService(repo)`

### 3. Handler / Controller (`internal/discord/handlers/`)
- **Trách nhiệm**: Nhận Discord interaction → validate → gọi Service → gọi UI Builder → respond
- **Không được**: Chứa SQL/MongoDB trực tiếp, chứa embed HTML/template logic
- **Đây là tầng duy nhất nhận `*discordgo.InteractionCreate`**
- **Ví dụ**: `start_handler.go`, `menu_handler.go`

### 4. UI Builder (`internal/discord/ui/`, `internal/discord/menu/*_menu.go`)
- **Trách nhiệm**: Nhận data struct thuần → trả về embed và component
- **Không được**: Gọi DB, gọi service, có side effects
- **Hàm Build* chỉ nhận data và trả về `*discordgo.InteractionResponseData`**
- **Ví dụ**: `BuildMainMenuResponse(data)`, `BuildProfileMenuResponse(data)`

### 5. Router (`internal/discord/router.go`, `internal/discord/menu/router.go`)
- **Trách nhiệm**: Phân loại và điều phối interactions đến đúng Handler
- **Không được**: Chứa business logic, chứa DB access
- **Top-level router**: Phân loại slash command vs. component interaction
- **Menu router**: Phân loại button/select trong menu → gọi PageLoader

## Menu Session Security Flow

```
User nhấn button
      │
      ▼
menu/router.go
      │ Parse custom_id: "domain:action:sessionId:extra"
      │ ValidateOwner(sessionId, userId)
      │   ├── Không tìm thấy session → "Giao diện đã hết hạn"
      │   ├── Session expired → "Giao diện đã hết hạn"
      │   └── userId không khớp → "Đây không phải giao diện của đạo hữu"
      │
      ▼ (session hợp lệ)
      Refresh TTL
      Dispatch đến action handler
      Call PageLoader → Service → UI Builder
      Edit existing message (không spam embed mới)
```

## Quy tắc custom_id

Format: `<domain>:<action>:<sessionId>[:<extra>]`

| Ví dụ | Ý nghĩa |
|---|---|
| `nav:refresh:abc123:main` | Làm mới trang main |
| `nav:back:abc123:main` | Quay về trang main |
| `nav:close:abc123` | Đóng menu |
| `menu:nav:abc123` | Select menu chọn category |
| `profile:rename:abc123` | Đổi đạo hiệu |
| `cultivation:meditate:abc123` | Tĩnh tu |

## Cấu trúc thư mục

```
cmd/bot/main.go                         ← Entrypoint, DI wiring

internal/
  config/config.go                      ← Load env vars
  logger/logger.go                      ← Structured logging (zap)
  errors/errors.go                      ← Sentinel errors + AppError

  database/
    mongo.go                            ← DB connect/disconnect
    indexes.go                          ← Index creation

  game/                                 ← Pure game domain (no Discord)
    profile/
      model.go       ← Struct Player
      repository.go  ← Interface + MongoDB impl
      service.go     ← Business logic
    cultivation/
      model.go
      repository.go
      service.go
    economy/
      model.go
      repository.go
      service.go
    cooldown/
      model.go
      repository.go
      service.go

  discord/
    bot.go           ← Connect, register commands, lifecycle
    commands.go      ← Slash command schemas
    router.go        ← Top-level interaction dispatcher

    handlers/        ← CONTROLLER LAYER
      start_handler.go
      menu_handler.go

    menu/            ← MENU SESSION + UI
      session.go     ← Session model + repository
      navigation.go  ← Session service
      router.go      ← Menu component interaction dispatcher
      main_menu.go   ← UI builder (main page)
      profile_menu.go
      cultivation_menu.go

    ui/              ← REUSABLE UI PRIMITIVES
      emojis.go      ← Custom Discord emoji registry
      colors.go      ← Color palette
      embeds.go      ← Standard embed builders
      components.go  ← Button/select builders
      messages.go    ← Vietnamese string constants

  server/
    health.go        ← GET /health handler
    http_server.go   ← HTTP server lifecycle

  scheduler/
    keepalive.go     ← Render self-ping
    cleanup_sessions.go ← Expired session cleanup

pkg/utils/
  time.go            ← Duration formatting, Discord timestamps
  id.go              ← Cryptographic session ID generation
  numbers.go         ← Number formatting, progress bar
```
