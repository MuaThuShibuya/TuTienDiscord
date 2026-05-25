// File: scripts/mongo_fix_cultivation_schema.js
// Phiên bản: v0.2
// Mục đích: Sửa dữ liệu cultivation_profiles cũ bị sai kiểu dữ liệu do nâng cấp version.
// Bảo mật: Không chứa MongoDB URI, username, password. Chạy script này trong mongosh hoặc MongoDB Atlas UI.
// Hướng dẫn chạy:
//   1. Mở MongoDB Atlas UI -> Database -> Collections
//   2. Chọn db "tu_tien_bot", mở mongosh (hoặc dán query này vào MongoDB Compass mongosh)
//   3. Chạy toàn bộ script dưới đây.

const collection = db.getCollection("cultivation_profiles");

print("Bắt đầu chạy migration cho collection cultivation_profiles...");

const result = collection.updateMany(
  {
    $or: [
      { mindState: { $type: "string" } },
      { realmLevel: { $type: "string" } },
      { cultivationExp: { $type: "string" } },
      { cultivationExpRequired: { $type: "string" } },
      { combatPower: { $type: "string" } },
      { stamina: { $type: "string" } },
      { maxStamina: { $type: "string" } }
    ]
  },
  [
    {
      $set: {
        mindState: {
          $convert: { input: "$mindState", to: "int", onError: 50, onNull: 50 }
        },
        realmLevel: {
          $convert: { input: "$realmLevel", to: "int", onError: 1, onNull: 1 }
        },
        cultivationExp: {
          $convert: { input: "$cultivationExp", to: "long", onError: 0, onNull: 0 }
        },
        cultivationExpRequired: {
          $convert: { input: "$cultivationExpRequired", to: "long", onError: 200, onNull: 200 }
        },
        combatPower: {
          $convert: { input: "$combatPower", to: "long", onError: 100, onNull: 100 }
        },
        stamina: {
          $convert: { input: "$stamina", to: "int", onError: 100, onNull: 100 }
        },
        maxStamina: {
          $convert: { input: "$maxStamina", to: "int", onError: 100, onNull: 100 }
        }
      }
    }
  ]
);

print(`Migration hoàn tất. Số documents được cập nhật: ${result.modifiedCount}`);