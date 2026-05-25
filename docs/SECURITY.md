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
- **Input sanitization**: `daoName` (đạo hiệu) phải trim, giới hạn độ dài, kiểm tra rune count.
- **No user-controlled logic**: Action names trong custom_id luôn là server-side constants, không phải text người dùng nhập.

## Economy Security

- Tất cả thay đổi currency dùng atomic MongoDB `$inc` + filter `balance >= deduct_amount`.
- Không bao giờ accept số lượng currency trực tiếp từ Discord user input.
- Currency âm bị chặn ở tầng repository (`adjustCurrency` filter).
- **Không gacha tiền thật**: Gacha chỉ nhận vé cơ duyên và linh ngọc trong game.
- **Không đấu giá tiền thật**: Toàn bộ thị trường chỉ dùng tài nguyên game.

## Anti-Race-Condition (Cần áp dụng từ v0.2+)

Các thao tác sau **phải** dùng atomic update hoặc idempotency key:

| Thao tác | Giải pháp |
|---|---|
| Gacha | Idempotency key + atomic pull count |
| Mua vật phẩm | Atomic $inc balance + item insert trong 1 session |
| Đấu giá | Optimistic lock + version field |
| Nhận thưởng boss | Idempotency key per (userId, bossId, round) |
| Đột phá cảnh giới | Atomic check-and-set realm |
| PvP | Combat session lock (userId không thể vào 2 PvP cùng lúc) |

## Content Safety

- **Đạo lữ/song tu**: Chỉ là buff đồng hành, nhiệm vụ đôi, kỹ năng hợp kích. Không tạo nội dung nhạy cảm hay tình dục.
- **NPC**: Chỉ tương tác game (nhiệm vụ, mua bán). Không roleplay nội dung nhạy cảm.
- **Chat**: Bot không lưu nội dung chat, không đọc tin nhắn người dùng.

## Database Security

- MongoDB URI chứa credentials — không log, không in ra stdout.
- Mọi query đều có context timeout (`database.NewContext()`).
- Index trên `(userId, guildId)` ngăn full-collection scan.
- TTL index tự xóa dữ liệu hết hạn (cooldown, session).
- Mỗi query **phải có** `guildId` filter để chống data leak giữa các server.

## Deployment Security

- Docker: chạy với non-root user (`botuser`).
- Render: đặt env vars trong Render dashboard, không trong code.
- VPS: dùng systemd `EnvironmentFile` hoặc vault, không file `.env` trên server.
- HTTPS: Render tự cấp TLS. VPS dùng Nginx + Let's Encrypt.
