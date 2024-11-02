package postgres

import "fmt"

func GetDSN(cfg *GormConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		cfg.DBName,
	)
}
