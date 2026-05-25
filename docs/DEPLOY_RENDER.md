# Hướng dẫn Deploy lên Render — Vạn Pháp Tiên Nghịch

## Yêu cầu trước khi deploy

- [ ] Tài khoản [Render](https://render.com) (free tier đủ dùng)
- [ ] Tài khoản [MongoDB Atlas](https://cloud.mongodb.com) với cluster đã tạo
- [ ] Discord Bot đã tạo tại [Discord Developer Portal](https://discord.com/developers/applications)
- [ ] Repository trên GitHub/GitLab (Render cần kết nối git để auto-deploy)

---

## Bước 1 — Chuẩn bị MongoDB Atlas

1. Tạo cluster Free (M0) tại MongoDB Atlas.
2. Tạo database user: **Database Access** → Add New User → chọn password auth.
3. Whitelist IP: **Network Access** → Add IP Address → `0.0.0.0/0` (allow all — Render dùng IP động).
4. Lấy connection string: **Connect** → Drivers → copy `mongodb+srv://...`.
5. **Quan trọng**: Đặt database name là `tu_tien_bot` — không dùng chung với bot khác.

---

## Bước 2 — Chuẩn bị Discord Bot

1. Vào [Discord Developer Portal](https://discord.com/developers/applications) → chọn ứng dụng.
2. **Bot** tab → lấy Token (Reset Token nếu cần).
3. **OAuth2** → chọn scope `bot` + `applications.commands` → lấy invite URL.
4. Mời bot vào server của bạn qua invite URL.
5. Lấy Guild ID: bật Developer Mode trong Discord → chuột phải server → Copy ID.

---

## Bước 3 — Tạo Web Service trên Render

1. Render Dashboard → **New** → **Web Service**.
2. Kết nối repository GitHub/GitLab.
3. Cấu hình:

| Trường | Giá trị |
|---|---|
| **Name** | `van-phap-tien-nghich` |
| **Region** | Singapore (gần nhất với VN) |
| **Branch** | `main` |
| **Root Directory** | *(để trống)* |
| **Runtime** | Go |
| **Build Command** | `go build -o bot ./cmd/bot` |
| **Start Command** | `./bot` |
| **Plan** | Free |

---

## Bước 4 — Đặt Environment Variables

Vào **Environment** tab của Web Service, thêm từng biến:

### Bắt buộc

```
DISCORD_TOKEN          = <token từ Developer Portal>
DISCORD_GUILD_IDS      = <guild ID, phân cách bằng dấu phẩy nếu nhiều server>
DISCORD_OWNER_IDS      = <Discord user ID của chủ bot>
MONGODB_URI            = mongodb+srv://<user>:<pass>@<cluster>.mongodb.net/?appName=tu-tien
MONGODB_DATABASE       = tu_tien_bot
```

### Cấu hình ứng dụng

```
APP_ENV                = production
APP_NAME               = tu-tien-discord-bot
APP_VERSION            = 0.1.2
COMMAND_REGISTER_MODE  = guild
PORT                   = 8080
MENU_SESSION_TTL_MINUTES = 15
```

### Logging — production

```
LOG_LEVEL   = info
LOG_FORMAT  = json
LOG_COLOR   = false
LOG_CALLER  = true
```

### Keepalive — bắt buộc trên Render free tier

Render free tier sleep sau ~15 phút không có traffic. Bot phải tự ping để tránh.

```
KEEPALIVE_ENABLED          = true
KEEPALIVE_URL              = https://<ten-app>.onrender.com/health
KEEPALIVE_INTERVAL_SECONDS = 600
```

> **Lưu ý**: Thay `<ten-app>` bằng tên Web Service thực tế của bạn trên Render.

---

## Bước 5 — Deploy

1. Click **Create Web Service** → Render tự build và deploy.
2. Theo dõi build log — thường mất 2-3 phút.
3. Khi thấy `Deploy live` → vào Discord, gõ `/start` để kiểm tra.

---

## Theo dõi Log trên Render

Render có tab **Logs** real-time. Vì `LOG_FORMAT=json`, mỗi dòng log có dạng:

```json
{"level":"info","time":"2026-05-25T14:16:07+07:00","logger":"discord.bot","msg":"Bot đã sẵn sàng","username":"Vạn Pháp Tu Tiên","guilds":1}
```

Tìm lỗi nhanh với filter:
- `"level":"error"` — lỗi nghiêm trọng
- `"level":"warn"` — cảnh báo
- `"logger":"discord.bot"` — log từ Discord layer

---

## Chuyển COMMAND_REGISTER_MODE sang global

Khi đã kiểm tra ổn định trên ít nhất 1 server:

1. Đổi `COMMAND_REGISTER_MODE=global` trong Render Environment.
2. Redeploy (hoặc bot tự restart sau khi save).
3. Chờ tối đa 1 giờ để Discord propagate lệnh toàn cầu.

> **Lưu ý**: Xóa guild commands cũ bằng cách gọi `COMMAND_REGISTER_MODE=guild` + deploy lần cuối rồi mới đổi sang global.

---

## Troubleshooting

### Bot không phản hồi lệnh

- Kiểm tra `DISCORD_TOKEN` có đúng không.
- Kiểm tra bot đã được mời vào server với scope `applications.commands`.
- Xem log có dòng `"msg":"Bot đã sẵn sàng"` chưa.

### Lỗi MongoDB kết nối

```
database: ping thất bại: ...
```

- Kiểm tra `MONGODB_URI` — đặc biệt password không được chứa ký tự đặc biệt chưa encode (`@` → `%40`).
- Kiểm tra Network Access Atlas đã whitelist `0.0.0.0/0`.

### Bot bị sleep và không tự wake

- Kiểm tra `KEEPALIVE_ENABLED=true`.
- Kiểm tra `KEEPALIVE_URL` trỏ đúng URL của app trên Render.
- Thử curl thủ công: `curl https://<ten-app>.onrender.com/health`.

### Lệnh không xuất hiện trong server

- Nếu `COMMAND_REGISTER_MODE=guild`: kiểm tra `DISCORD_GUILD_IDS` có đúng Guild ID không.
- Guild commands cập nhật tức thì; global commands cần tối đa 1 giờ.

---

## Upgrade lên VPS (tùy chọn)

Khi cần uptime 100% và không muốn dùng keepalive, chuyển sang VPS:

```bash
# Trên VPS Ubuntu/Debian
go build -o /opt/tu-tien-bot/bot ./cmd/bot

# /etc/systemd/system/tu-tien-bot.service
[Unit]
Description=Van Phap Tien Nghich Discord Bot
After=network.target

[Service]
Type=simple
User=botuser
WorkingDirectory=/opt/tu-tien-bot
EnvironmentFile=/etc/tu-tien-bot/env
ExecStart=/opt/tu-tien-bot/bot
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# env file /etc/tu-tien-bot/env — chỉ root đọc được
chmod 600 /etc/tu-tien-bot/env

# Thêm biến production vào đây (không commit lên git!)
LOG_FORMAT=json
LOG_COLOR=false
KEEPALIVE_ENABLED=false
...
```

```bash
systemctl enable tu-tien-bot
systemctl start tu-tien-bot
journalctl -u tu-tien-bot -f  # theo dõi log
```
