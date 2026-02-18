package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Tire struct {
	ID           int
	Title        string
	Type         string // summer, winter, allseason
	Width        int    // mm
	Profile      int    // %
	Diameter     int    // inches
	Photo        string
	Video        string
	Description  string
	BasePressure float64 // base pressure in bar
}

func (r *Repository) GetTires() ([]Tire, error) {
	tires := []Tire{
		{
			ID:           1,
			Title:        "Michelin Pilot Sport 4S",
			Type:         "summer",
			Width:        225,
			Profile:      45,
			Diameter:     17,
			Photo:        "Michelin Pilot Sport 4S.jpg",
			Video:        "edit.mp4",
			Description:  "Летняя шина для спортивных автомобилей",
			BasePressure: 2.2,
		},
		{
			ID:           2,
			Title:        "Nokian Hakkapeliitta R5",
			Type:         "winter",
			Width:        205,
			Profile:      55,
			Diameter:     16,
			Photo:        "Nokian Hakkapeliitta R5.jpg",
			Video:        "edit.mp4",
			Description:  "Зимняя шина с отличным сцеплением на льду",
			BasePressure: 2.3,
		},
		{
			ID:           3,
			Title:        "Continental AllSeasonContact",
			Type:         "allseason",
			Width:        195,
			Profile:      65,
			Diameter:     15,
			Photo:        "Continental AllSeasonContact.jpg",
			Video:        "edit.mp4",
			Description:  "Всесезонная шина для умеренного климата",
			BasePressure: 2.2,
		},
		{
			ID:           4,
			Title:        "Bridgestone Potenza RE003",
			Type:         "summer",
			Width:        235,
			Profile:      40,
			Diameter:     18,
			Photo:        "Bridgestone Potenza RE003.jpeg",
			Video:        "edit.mp4",
			Description:  "Летняя шина для мощных седанов",
			BasePressure: 2.3,
		},
		{
			ID:           5,
			Title:        "Goodyear UltraGrip Ice 2",
			Type:         "winter",
			Width:        215,
			Profile:      60,
			Diameter:     16,
			Photo:        "Goodyear UltraGrip Ice 2.jpg",
			Video:        "edit.mp4",
			Description:  "Зимняя шина для суровых зим",
			BasePressure: 2.4,
		},
		{
			ID:           6,
			Title:        "Vredestein Quatrac Pro",
			Type:         "allseason",
			Width:        205,
			Profile:      50,
			Diameter:     17,
			Photo:        "Vredestein Quatrac Pro.jpg",
			Video:        "edit.mp4",
			Description:  "Всесезонная шина премиум класса",
			BasePressure: 2.2,
		},
		{
			ID:           7,
			Title:        "Yokohama Advan Neova",
			Type:         "summer",
			Width:        245,
			Profile:      35,
			Diameter:     19,
			Photo:        "Yokohama Advan Neova.jpg",
			Video:        "edit.mp4",
			Description:  "Спортивная летняя шина",
			BasePressure: 2.4,
		},
		{
			ID:           8,
			Title:        "Pirelli Winter Cinturato",
			Type:         "winter",
			Width:        185,
			Profile:      60,
			Diameter:     14,
			Photo:        "Pirelli Winter Cinturato.jpg",
			Video:        "edit.mp4",
			Description:  "Зимняя шина для городских автомобилей",
			BasePressure: 2.1,
		},
	}
	if len(tires) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return tires, nil
}

func (r *Repository) GetTire(id int) (Tire, error) {
	tires, err := r.GetTires()
	if err != nil {
		return Tire{}, err
	}

	for _, tire := range tires {
		if tire.ID == id {
			return tire, nil
		}
	}
	return Tire{}, fmt.Errorf("Шина не найдена")
}

func (r *Repository) GetTireByTitle(title string) ([]Tire, error) {
	tires, err := r.GetTires()
	if err != nil {
		return []Tire{}, err
	}

	var result []Tire
	for _, tire := range tires {
		if strings.Contains(strings.ToLower(tire.Title), strings.ToLower(title)) {
			result = append(result, tire)
		}
	}
	return result, nil
}

func (r *Repository) GetRequest() ([]Tire, error) {
	tires, err := r.GetTires()
	if err != nil {
		return []Tire{}, err
	}

	var result []Tire
	for _, tire := range tires {
		if tire.ID == 1 || tire.ID == 5 {
			result = append(result, tire)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return result, nil
}

// CalculatePressure calculates optimal tire pressure based on parameters
// temperature: ambient temperature in Celsius
// weight: vehicle weight with load in kg
// surfaceCoeff: road surface coefficient (asphalt=1.0, gravel=0.9, snow=0.8, ice=0.7)
func CalculatePressure(basePressure float64, temperature float64, weight float64, surfaceCoeff float64) float64 {
	// Temperature adjustment: +0.1 bar per 10°C below 20°C
	tempDiff := 20.0 - temperature
	tempAdjustment := tempDiff * 0.01 // 0.01 bar per 1°C

	// Weight adjustment: base is 1500kg, +0.1 bar per 500kg over
	weightAdjustment := 0.0
	if weight > 1500 {
		weightAdjustment = (weight - 1500) / 500 * 0.1
	}

	// Surface coefficient adjustment
	surfaceAdjustment := (1.0 - surfaceCoeff) * 0.2

	// Calculate final pressure
	pressure := basePressure + tempAdjustment + weightAdjustment - surfaceAdjustment

	// Keep pressure in reasonable range (1.8 - 3.0 bar)
	if pressure < 1.8 {
		pressure = 1.8
	}
	if pressure > 3.0 {
		pressure = 3.0
	}

	return pressure
}
