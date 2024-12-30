package gorm

import (
	"github.com/go-gorm/caches/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func UseCachePlugin(db *gorm.DB, cacher caches.Cacher) error {
	log.Debug().Msg("setup database cache plugin")
	cachePlugin := &caches.Caches{
		Conf: &caches.Config{Cacher: cacher},
	}
	return db.Use(cachePlugin)
}
