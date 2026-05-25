# Security Guidelines — Vạn Pháp Tiên Nghịch

## Secrets Management

- **Không bao giờ hardcode**: Discord token, MongoDB URI, Guild ID, Owner ID, webhook URL, secret key.
- Tất cả secret đọc từ environment variables qua `internal/config/config.go`.
- File `.env` chỉ dùng local, không commit vào git (có trong `.gitignore`).
- Chỉ commit `.env.example` với các giá trị trống.
- Không log token, URI, password, private key ở bất kỳ log level nào.

## Discord Interaction Security

- **Guild-only**: Tất cả commands chỉ hoạt động trong server (guild), không trong DM.
- **Menu ownership**: Mỗi menu session gắn với một `userId`. Người khác bấm menu của bạn nhận lỗi ephemeral.
- **Session ID**: Dùng `crypto/rand` (16 bytes = 32 hex chars). Không thể đoán được.
- **custom_id validation**: Mọi button/select đều phải parse và validate sessionId + userId trước khi xử lý.
- **Input sanitization**: `daoName` (đạo hiệu) phải trim, giới hạn rune count — dùng `utf8.RuneCountInString`, không dùng `len()` (sai với UTF-8 tiếng Việt).
- **No user-controlled logic**: Action names trong custom_id luôn là server-side constants, không phải text người dùng nhập.

## Economy Security

### Ngăn chặn Double-Spend

Lỗ hổng double-spend xảy ra khi user gửi nhiều request đồng thời để tiêu cùng một khoản tiền:

```
Goroutine A: read balance=500 → 500 >= 300 ✓ → write balance=200
Goroutine B: read balance=500 → 500 >= 300 ✓ → write balance=200   ← DOUBLE SPEND!
```

**Giải pháp bắt buộc cho MongoDB repository**:

```go
// SAI: đọc trước, check, rồi update → TOCTOU race condition
wallet, _ := col.FindOne(filter)
if wallet.SpiritStones >= amount { col.UpdateOne(..., $inc) }

// ĐÚNG: atomic conditional update — chỉ update khi balance đủ
result := col.FindOneAndUpdate(
    bson.M{"userId": uid, "guildId": gid, "spiritStones": bson.M{"$gte": amount}},
    bson.M{"$inc": bson.M{"spiritStones": -amount}},
    options.FindOneAndUpdate().SetReturnDocument(options.After),
)
if result.Err() == mongo.ErrNoDocuments {
    return nil, apperrors.ErrInsufficientFunds  // balance không đủ HOẶC user không tồn tại
}
```

TTL của câu lệnh: đây là một round trip DB duy nhất, atomic theo MongoDB's document-level locking.

### Currency Rules

- Tất cả thay đổi currency dùng atomic MongoDB `$inc` kết hợp filter `balance >= amount`.
- Không bao giờ accept số lượng currency trực tiếp từ Discord user input.
- Số âm bị chặn tại tầng service (`EarnSpiritStones` từ chối `amount <= 0`).
- **Không gacha tiền thật**: Gacha chỉ nhận vé cơ duyên (FateTickets) và linh ngọc (SpiritJades).
- **Không đấu giá tiền thật**: Thị trường chỉ dùng tài nguyên in-game.

## Anti-Race-Condition

### Đã bảo vệ (v0.1.x)

| Thao tác | Cơ chế bảo vệ |
|---|---|
| Chi tiêu linh thạch | Atomic `$inc` với filter `balance >= amount` |
| Chi tiêu linh ngọc | Atomic `$inc` với filter `balance >= amount` |
| Chi tiêu vé cơ duyên | Atomic `$inc` với filter `fateTickets >= amount` |
| Menu session ownership | SessionID ngẫu nhiên + validate userId mỗi interaction |
| Cooldown | TTL index tự xóa + Upsert thay vì insert |

### Cần bảo vệ từ v0.2+

| Thao tác | Giải pháp đề xuất |
|---|---|
| Gacha nhiều vé | Idempotency key + atomic decrement vé trước khi roll |
| Nhận thưởng nhiệm vụ | `insertOne` với unique index `(userId, questId, round)` |
| Đột phá cảnh giới | Atomic compare-and-swap: `{"realm": currentRealm, "exp": {"$gte": required}}` |
| PvP | Combat session lock: unique index `(userId, "in_pvp")` — xóa khi combat kết thúc |
| Mua vật phẩm | MongoDB transaction hoặc atomic item+balance trong 1 session |
| Nhận thưởng boss | Idempotency key `(userId, bossId, spawnRound)` |

## Content Safety

- **Đạo lữ/song tu**: Chỉ là buff đồng hành, nhiệm vụ đôi, kỹ năng hợp kích. Không tạo nội dung nhạy cảm.
- **NPC**: Chỉ tương tác game (nhiệm vụ, mua bán). Không roleplay nội dung không phù hợp.
- **Chat**: Bot không lưu nội dung chat, không đọc tin nhắn người dùng.

## Database Security

- MongoDB URI chứa credentials — không log, không in ra stdout dù ở level debug.
- Mọi query đều có context timeout (`database.NewContext()` hoặc `context.WithTimeout`).
- Index trên `(userId, guildId)` ngăn full-collection scan.
- TTL index tự xóa dữ liệu hết hạn (cooldown, menu_session).
- Mọi query **bắt buộc có** `guildId` filter để chống data leak giữa các Discord server.
- Database name riêng biệt (`tu_tien_bot`) — không dùng chung với bot khác.

## Logging Security

- **Không log**: `DISCORD_TOKEN`, `MONGODB_URI`, password, API key.
- **LOG_FORMAT=json** khi production: structured log không chứa ANSI escape code.
- **LOG_COLOR=false** khi production/CI: tránh làm ô nhiễm log aggregator (Datadog, Loki, CloudWatch).
- Log level `info` cho production — tránh `debug` vì có thể tiết lộ internal state.

## Deployment Security

- Docker: chạy với non-root user (`botuser`).
- Render: đặt env vars trong Render Dashboard → Environment, không trong code hoặc Dockerfile.
- VPS: dùng systemd `EnvironmentFile=/etc/bot/env` hoặc Vault, không file `.env` trên server.
- HTTPS: Render tự cấp TLS. VPS dùng Nginx reverse proxy + Let's Encrypt (`certbot`).
