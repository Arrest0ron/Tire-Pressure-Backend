package ds

import (
	"time"
)

type TirePressure struct {
	ID             int                 `gorm:"primaryKey" json:"id"`
	AirTemperature float64             `gorm:"not null" json:"air_temperature"`
	CarWeight      float64             `gorm:"not null" json:"car_weight"`
	DateCreate     time.Time           `gorm:"autoCreateTime" json:"date_create"`
	DateUpdate     time.Time           `gorm:"autoUpdateTime" json:"date_update"`
	Status         string              `gorm:"type:varchar(20);not null;check:status IN ('draft', 'deleted', 'formed', 'completed', 'rejected')" json:"status"`
	CreatorID      uint                `gorm:"not null" json:"creator_id"`
	Entries        []TirePressureEntry `gorm:"foreignKey:TirePressureID" json:"entries"`
	EntryCount     int                 `gorm:"default:0" json:"entry_count"`
	TotalPressure  float64             `gorm:"type:float;default:0" json:"total_pressure"`

	Creator Users `gorm:"foreignKey:CreatorID"`
}
