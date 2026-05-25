# Version Plan — Vạn Pháp Tiên Nghịch

## v0.1 — Foundation ✅ (đã xong)
- Config từ env vars
- MongoDB Atlas connection + indexes
- Discord bot boot + graceful shutdown
- `/start` — đăng ký người chơi
- `/menu` — Main Menu UI
- User profile model
- Cultivation profile cơ bản
- Wallet (linh thạch, linh ngọc, vé cơ duyên)
- Cooldown infrastructure
- Menu session (chống người khác bấm nhầm)
- Health check endpoint (`GET /health`)
- Keepalive cho Render
- Structured logging (zap)
- Error handling chuẩn

## v0.2 — Cultivation Core ✅ (đã xong)
- Tĩnh tu (cooldown, exp gain, stamina cost)
- Bế quan (nhân hệ số tu vi, lock stamina)
- Tu vi và tiến độ cảnh giới
- Đột phá cảnh giới (chance-based, trừ linh thạch)
- Tâm cảnh ảnh hưởng hiệu suất tu luyện
- Chọn đạo lộ (kiếm tu, thể tu, linh tu, độc tu - select 1 lần)
- Cooldown chi tiết cho từng hành động

## v0.3 — Inventory & Equipment (hiện tại)
- Túi đồ (giới hạn ô)
- Item model (instance-based)
- Đan dược cơ bản (tăng exp, hồi stamina, buff tỉ lệ đột phá)
- Lò luyện đan dược
- Trang bị (vũ khí, áo giáp, linh khí, pháp bảo)
- Mặc/tháo trang bị
- Cường hóa cơ bản (tiêu thụ nguyên liệu)
- Phân giải vật phẩm → nguyên liệu

## v0.4 — Combat PvE
- Quái thường (mỗi zone một bộ quái)
- Vượt ải (tiến độ)
- Combat theo lượt với kỹ năng cơ bản
- Công thức sát thương (ATK, DEF, crit)
- Thưởng sau trận (exp, linh thạch, drop)
- Bản đồ thế giới cơ bản

## v0.5 — Gacha (Cơ Duyên)
- Banner hệ thống
- Tỉ lệ drop (theo tier)
- Pity system (đảm bảo tier cao sau N lần)
- Roll 1 / Roll 10
- Lịch sử quay
- Trùng vật phẩm → đổi mảnh
- **Chỉ dùng vé cơ duyên/linh ngọc trong game, không tiền thật**

## v0.6 — Linh Thú / Con Rối
- Gacha linh thú
- Nuôi dưỡng (cho ăn, tăng thân thiết)
- Tăng cấp, tăng sao linh thú
- Ra trận (passive buff hoặc active skill)
- Hệ thống linh thú hỗ trợ combat

## v0.7 — Boss & Dungeon
- Boss server theo giờ (mỗi X giờ xuất hiện)
- Phó bản bảo vật (giới hạn lần/ngày)
- Bảng xếp hạng sát thương boss
- Reward pool (chia theo đóng góp)
- Thông báo boss xuất hiện

## v0.8 — Market & Auction
- Cửa hàng hệ thống (NPC shop)
- Chợ người chơi (đăng bán vật phẩm)
- Đấu giá (bid timer, auto-close)
- Phí giao dịch (sink kinh tế)
- Transaction log (audit trail)
- Anti-race-condition cho mọi giao dịch

## v0.9 — PvP & Resource Contest
- PvP tự nguyện (duel)
- Bí cảnh tranh đoạt tài nguyên
- Bảo hộ tân thủ (không bị tấn công nếu quá yếu)
- Giới hạn cướp tài nguyên (anti-toxic)
- Mùa PvP (xếp hạng reset)

## v1.0 — Sect / Social / NPC
- Tông môn (tạo, gia nhập, buff thành viên)
- NPC (nhiệm vụ, mua bán đặc biệt)
- Nhiệm vụ hàng ngày / tuần / cốt truyện
- Đạo lữ (buff đồng hành, nhiệm vụ đôi, kỹ năng hợp kích)
  - **Không tạo nội dung nhạy cảm**
- Cốt truyện mùa (events theo thời gian thực)

## Ghi chú kỹ thuật theo version

| Version | Kỹ thuật mới cần thêm |
|---|---|
| v0.2 | Cooldown atomic check + set trong 1 transaction |
| v0.3 | Item instance ID, inventory lock khi equip |
| v0.4 | Combat state machine, session combat |
| v0.5 | Pity counter atomic, gacha idempotency key |
| v0.7 | Boss spawn scheduler, DPS leaderboard |
| v0.8 | Auction timer goroutine, market atomic buy |
| v0.9 | PvP matchmaking, resource lock |
| v1.0 | Guild data model, NPC dialogue tree |
