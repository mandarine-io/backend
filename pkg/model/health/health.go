package health

type HealthOutput struct {
	Name string `json:"name" binding:"required"`
	Pass bool   `json:"pass" binding:"required"`
}
