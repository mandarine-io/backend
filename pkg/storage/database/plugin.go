package database

import (
	"github.com/go-gorm/caches/v4"
	"gorm.io/gorm"
	"log/slog"
)

func UseCachePlugin(db *gorm.DB, cacher caches.Cacher) error {
	slog.Info("Setup database cache plugin")
	cachePlugin := &caches.Caches{
		Conf: &caches.Config{Cacher: cacher},
	}
	return db.Use(cachePlugin)
}
