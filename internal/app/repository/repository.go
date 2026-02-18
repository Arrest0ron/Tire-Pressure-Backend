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
	MM           string  // количество/порядок/главный/комментарий
}

// Request - структура заявки (словарь)
type Request struct {
	ID          int
	TireIDs     []int // список ID шин в заявке
	Temperature float64
	Weight      float64
	Surface     float64
	Result      float64 // результат вычислений
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
			MM:           "4 шт., основная",
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
			MM:           "4 шт.",
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
			MM:           "4 шт.",
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
			MM:           "4 шт.",
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
			MM:           "2 шт., задняя ось",
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
			MM:           "4 шт.",
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
			MM:           "4 шт.",
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
			MM:           "4 шт.",
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
			ID:          1,
			TireIDs:     []int{1, 5},
			Temperature: 20,
			Weight:      1500,
			Surface:     1.0,
			Result:      2.2,
		},
		2: {
			ID:          2,
			TireIDs:     []int{2, 4},
			Temperature: 15,
			Weight:      1800,
			Surface:     0.9,
			Result:      2.4,
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
