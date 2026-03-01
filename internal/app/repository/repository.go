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
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Width            int     `json:"width"`
	Profile          int     `json:"profile"`
	Diameter         int     `json:"diameter"`
	Type             string  `json:"type"`
	TireCoefficient  float64 `json:"tire_coefficient"`
	ServiceDescription string  `json:"service_description"`
	Photo            string  `json:"photo"`
	Video            string  `json:"video"`
}

// Request - структура заявки (словарь)
type Request struct {
	ID               int
	TireIDs          []int // список ID шин в заявке
	AirTemperature   float64
	CarWeight        float64
	SurfaceCoefficient float64
	PressureResult   float64 // результат вычислений
}

func (r *Repository) GetTires() ([]Tire, error) {
	tires := []Tire{
		{
			ID:               1,
			Title:            "Michelin Pilot Sport 4S",
			Type:             "summer",
			Width:            225,
			Profile:          45,
			Diameter:         17,
			Photo:            "Michelin Pilot Sport 4S.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Высокопроизводительная летняя шина для спортивных автомобилей. Отличная управляемость на сухом и мокром покрытии, низкий уровень шума, высокий ресурс. Идеально подходит для мощных седанов и купе класса премиум.",
			TireCoefficient:  2.2,
		},
		{
			ID:               2,
			Title:            "Nokian Hakkapeliitta R5",
			Type:             "winter",
			Width:            205,
			Profile:          55,
			Diameter:         16,
			Photo:            "Nokian Hakkapeliitta R5.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Профессиональная зимняя шина с уникальным рисунком протектора для максимального сцепления на льду и снегу. Отличная курсовая устойчивость, короткий тормозной путь, комфортная езда. Подходит для суровых зимних условий.",
			TireCoefficient:  2.3,
		},
		{
			ID:               3,
			Title:            "Continental AllSeasonContact",
			Type:             "allseason",
			Width:            195,
			Profile:          65,
			Diameter:         15,
			Photo:            "Continental AllSeasonContact.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Всесезонная шина премиум-класса для умеренного климата. Сбалансированное сочетание летних и зимних характеристик, низкое сопротивление качению, отличное сцепление на мокрой дороге. Идеальный выбор для повседневной эксплуатации круглый год.",
			TireCoefficient:  2.2,
		},
		{
			ID:               4,
			Title:            "Bridgestone Potenza RE003",
			Type:             "summer",
			Width:            235,
			Profile:          40,
			Diameter:         18,
			Photo:            "Bridgestone Potenza RE003.jpeg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Спортивная летняя шина для мощных седанов и автомобилей бизнес-класса. Высокая курсовая устойчивость, точное управление, отличное сцепление на сухом покрытии. Подходит для агрессивного стиля вождения и высоких скоростей.",
			TireCoefficient:  2.3,
		},
		{
			ID:               5,
			Title:            "Goodyear UltraGrip Ice 2",
			Type:             "winter",
			Width:            215,
			Profile:          60,
			Diameter:         16,
			Photo:            "Goodyear UltraGrip Ice 2.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Зимняя шина для суровых зимних условий. Трехмерный ламелированный рисунок протектора обеспечивает отличное сцепление на льду и снегу. Надежное торможение, устойчивость на скользких участках, комфортная езда.",
			TireCoefficient:  2.4,
		},
		{
			ID:               6,
			Title:            "Vredestein Quatrac Pro",
			Type:             "allseason",
			Width:            205,
			Profile:          50,
			Diameter:         17,
			Photo:            "Vredestein Quatrac Pro.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Всесезонная шина премиум класса с инновационной резиновой смесью. Отличное сцепление в любых погодных условиях, низкий уровень шума, высокий ресурс. Подходит для современных автомобилей с высокими требованиями к комфорту.",
			TireCoefficient:  2.2,
		},
		{
			ID:               7,
			Title:            "Yokohama Advan Neova",
			Type:             "summer",
			Width:            245,
			Profile:          35,
			Diameter:         19,
			Photo:            "Yokohama Advan Neova.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Спортивная летняя шина для высокопроизводительных автомобилей. Улучшенная управляемость и точность рулевого управления, отличное сцепление на сухом покрытии, низкое сопротивление качению. Подходит для спортивного вождения и динамичной езды.",
			TireCoefficient:  2.4,
		},
		{
			ID:               8,
			Title:            "Pirelli Winter Cinturato",
			Type:             "winter",
			Width:            185,
			Profile:          60,
			Diameter:         14,
			Photo:            "Pirelli Winter Cinturato.jpg",
			Video:            "yellow_black.mp4",
			ServiceDescription: "Зимняя шина для городских автомобилей и компактных автомобилей. Отличное сцепление на снегу и льду, короткий тормозной путь, комфортная и тихая езда. Идеальный выбор для безопасной зимней эксплуатации в городских условиях.",
			TireCoefficient:  2.1,
		},
	}
	if len(tires) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return tires, nil
}

// GetRequests - словарь заявок (вторая коллекция)
func (r *Repository) GetRequests() (map[int]Request, error) {
	requests := map[int]Request{
		1: {
			ID:               1,
			TireIDs:          []int{1, 5},
			AirTemperature:   20,
			CarWeight:        1500,
			SurfaceCoefficient: 1.0,
			PressureResult:   2.2,
		},
		2: {
			ID:               2,
			TireIDs:          []int{2, 4},
			AirTemperature:   15,
			CarWeight:        1800,
			SurfaceCoefficient: 0.9,
			PressureResult:   2.4,
		},
		3: {
			ID:               3,
			TireIDs:          []int{3, 6},
			AirTemperature:   25,
			CarWeight:        1600,
			SurfaceCoefficient: 1.1,
			PressureResult:   2.1,
		},
		4: {
			ID:               4,
			TireIDs:          []int{7, 8},
			AirTemperature:   10,
			CarWeight:        1700,
			SurfaceCoefficient: 0.8,
			PressureResult:   2.5,
		},
		5: {
			ID:               5,
			TireIDs:          []int{1, 3, 5},
			AirTemperature:   18,
			CarWeight:        1900,
			SurfaceCoefficient: 1.2,
			PressureResult:   2.3,
		},
		6: {
			ID:               6,
			TireIDs:          []int{2, 6, 8},
			AirTemperature:   5,
			CarWeight:        2000,
			SurfaceCoefficient: 0.7,
			PressureResult:   2.6,
		},
		7: {
			ID:               7,
			TireIDs:          []int{4, 7},
			AirTemperature:   30,
			CarWeight:        1400,
			SurfaceCoefficient: 1.3,
			PressureResult:   2.0,
		},
		8: {
			ID:               8,
			TireIDs:          []int{1, 2, 3, 4},
			AirTemperature:   0,
			CarWeight:        2100,
			SurfaceCoefficient: 0.6,
			PressureResult:   2.7,
		},
	}
	if len(requests) == 0 {
		return nil, fmt.Errorf("Словарь пустой")
	}

	return requests, nil
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

// GetRequestByID - получить заявку по ID из словаря
func (r *Repository) GetRequestByID(id int) (Request, error) {
	requests, err := r.GetRequests()
	if err != nil {
		return Request{}, err
	}

	request, exists := requests[id]
	if !exists {
		return Request{}, fmt.Errorf("Заявка не найдена")
	}

	return request, nil
}

// GetRequestTires - получить шины для конкретной заявки
func (r *Repository) GetRequestTires(requestID int) ([]Tire, error) {
	request, err := r.GetRequestByID(requestID)
	if err != nil {
		return []Tire{}, err
	}

	tires, err := r.GetTires()
	if err != nil {
		return []Tire{}, err
	}

	var result []Tire
	for _, tireID := range request.TireIDs {
		for _, tire := range tires {
			if tire.ID == tireID {
				result = append(result, tire)
				break
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Заявка пуста")
	}

	return result, nil
}

// CalculatePressure - рассчитать давление
func CalculatePressure(basePressure float64, temperature float64, weight float64, surfaceCoeff float64) float64 {
	tempDiff := 20.0 - temperature
	tempAdjustment := tempDiff * 0.01

	weightAdjustment := 0.0
	if weight > 1500 {
		weightAdjustment = (weight - 1500) / 500 * 0.1
	}

	surfaceAdjustment := (1.0 - surfaceCoeff) * 0.2

	pressure := basePressure + tempAdjustment + weightAdjustment - surfaceAdjustment

	if pressure < 1.8 {
		pressure = 1.8
	}
	if pressure > 3.0 {
		pressure = 3.0
	}

	return pressure
}
