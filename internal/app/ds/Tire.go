package ds

type Tire struct {
	ID              int     `gorm:"primaryKey" json:"id"`
	Title           string  `gorm:"type:varchar(255);not null" json:"title"`
	Description     string  `gorm:"type:text" json:"description"`
	Photo           string  `gorm:"type:varchar(255)" json:"photo"`
	Video           string  `gorm:"type:varchar(255)" json:"video"`
	TireCoefficient float64 `gorm:"not null" json:"tire_coefficient"`
	IsDelete        bool    `gorm:"type:boolean;default:false" json:"is_delete"`
}