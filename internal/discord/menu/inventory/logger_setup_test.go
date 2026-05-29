package invmenu

import (
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

func init() {
	_ = logger.Init(logger.Options{Level: "error", Format: "json"})
}
