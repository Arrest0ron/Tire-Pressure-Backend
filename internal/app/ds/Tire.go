package ds

type Tire struct {
	TireID                   uint    `gorm:"primaryKey"`
	TireTitle                string  `gorm:"type:varchar(100);not null"`
	Description              string  `gorm:"type:text"`
	Photo                    string  `gorm:"type:varchar(255)"`
	Video                    string  `gorm:"type:varchar(255)"`
	TireMaterialCoefficient  float64 `gorm:"not null"`
	TireThicknessCoefficient float64 `gorm:"not null"`
	IsDelete                 bool    `gorm:"type:boolean;not null;default:false"`
}

// Done
