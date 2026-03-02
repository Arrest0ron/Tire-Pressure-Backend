package repository

import (
	"fmt"
	"strings"
)

// Repository — хранилище данных (Lab 1: данные в массивах, без БД).
type Repository struct {
}

// NewRepository создаёт новый экземпляр репозитория.
func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// Tire — услуга: шина с коэффициентами.
type Tire struct {
	ID              int
	Title           string  // название шины
	Description     string  // описание шины
	Photo           string  // ключ изображения в Minio
	Video           string  // ключ видео в Minio
	TireCoefficient float64 // коэффициент шины
}

// TirePressure — заявка: параметры для расчёта давления в шинах (Tire_pressure).
type TirePressure struct {
	ID            int
	Title         string
	Description   string
	Entries       []TirePressureEntry
	EntryCount    int
}

// TirePressureEntry — связь м-м: шина + параметры, результат — давление.
type TirePressureEntry struct {
	Tire           Tire
	CoatingCoeff   float64 // коэффициент покрытия (м-м поле)
	AirTemperature float64 // температура воздуха (поле Заявка)
	CarWeight      float64 // вес авто (поле Заявка)
	Pressure       float64 // результат: давление в шине
}

// GetTires возвращает все шины (услуги).
func (r *Repository) GetTires() ([]Tire, error) {
	tires := []Tire{
		{
			ID:              1,
			Title:           "Michelin Pilot Sport 4S",
			Description:     "Высокопроизводительная летняя шина для спортивных автомобилей. Отличная управляемость на сухом и мокром покрытии.",
			Photo:           "Michelin Pilot Sport 4S.jpg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.2,
		},
		{
			ID:              2,
			Title:           "Nokian Hakkapeliitta R5",
			Description:     "Профессиональная зимняя шина с уникальным рисунком протектора для максимального сцепления на льду и снегу.",
			Photo:           "Nokian Hakkapeliitta R5.jpg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.3,
		},
		{
			ID:              3,
			Title:           "Continental AllSeasonContact",
			Description:     "Всесезонная шина премиум-класса для умеренного климата. Сбалансированное сочетание характеристик.",
			Photo:           "Continental AllSeasonContact.jpg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.2,
		},
		{
			ID:              4,
			Title:           "Bridgestone Potenza RE003",
			Description:     "Спортивная летняя шина для мощных седанов. Высокая курсовая устойчивость и точное управление.",
			Photo:           "Bridgestone Potenza RE003.jpeg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.3,
		},
		{
			ID:              5,
			Title:           "Goodyear UltraGrip Ice 2",
			Description:     "Зимняя шина для суровых условий. Трехмерный ламелированный рисунок для сцепления на льду.",
			Photo:           "Goodyear UltraGrip Ice 2.jpg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.4,
		},
		{
			ID:              6,
			Title:           "Vredestein Quatrac Pro",
			Description:     "Всесезонная шина премиум класса с инновационной резиновой смесью для любых погодных условий.",
			Photo:           "Vredestein Quatrac Pro.jpg",
			Video:           "yellow_black.mp4",
			TireCoefficient: 2.2,
		},
	}
	if len(tires) == 0 {
		return nil, fmt.Errorf("массив шин пуст")
	}
	return tires, nil
}

// GetTire возвращает шину по ID.
func (r *Repository) GetTire(id int) (Tire, error) {
	tires, err := r.GetTires()
	if err != nil {
		return Tire{}, err
	}
	for _, t := range tires {
		if t.ID == id {
			return t, nil
		}
	}
	return Tire{}, fmt.Errorf("шина не найдена")
}

// GetTiresByTitle возвращает шины, содержащие подстроку в названии или описании.
func (r *Repository) GetTiresByTitle(query string) ([]Tire, error) {
	tires, err := r.GetTires()
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var result []Tire
	for _, t := range tires {
		if strings.Contains(strings.ToLower(t.Title), q) || strings.Contains(strings.ToLower(t.Description), q) {
			result = append(result, t)
		}
	}
	return result, nil
}

// CalculatePressure вычисляет давление по коэффициентам и параметрам нагрузки.
func CalculatePressure(tireCoeff, coatingCoeff float64, airTemp, carWeight float64) float64 {
	if coatingCoeff <= 0 {
		coatingCoeff = 1
	}
	tempFactor := (20.0 - airTemp) * 0.05
	weightFactor := (carWeight - 1500) * 0.001
	return (tempFactor + weightFactor + 1) * tireCoeff * coatingCoeff * 10
}

// buildTirePressure собирает заявку Tire_pressure из записей м-м.
func (r *Repository) buildTirePressure(id int, title, description string, entries []struct {
	TireID         int
	CoatingCoeff   float64
	AirTemperature float64
	CarWeight      float64
}) (TirePressure, error) {
	tires, err := r.GetTires()
	if err != nil {
		return TirePressure{}, err
	}
	tireMap := make(map[int]Tire)
	for _, t := range tires {
		tireMap[t.ID] = t
	}
	var calcEntries []TirePressureEntry
	for _, e := range entries {
		tire, ok := tireMap[e.TireID]
		if !ok {
			continue
		}
		pressure := CalculatePressure(tire.TireCoefficient, e.CoatingCoeff, e.AirTemperature, e.CarWeight)
		calcEntries = append(calcEntries, TirePressureEntry{
			Tire:           tire,
			CoatingCoeff:   e.CoatingCoeff,
			AirTemperature: e.AirTemperature,
			CarWeight:      e.CarWeight,
			Pressure:       pressure,
		})
	}
	return TirePressure{
		ID:            id,
		Title:         title,
		Description:   description,
		Entries:       calcEntries,
		EntryCount:    len(calcEntries),
	}, nil
}

// GetTirePressures возвращает все заявки на расчёт давления (Tire_pressure).
func (r *Repository) GetTirePressures() ([]TirePressure, error) {
	entries := []struct {
		TireID         int
		CoatingCoeff   float64
		AirTemperature float64
		CarWeight      float64
	}{
		{1, 1.0, 20, 1500},
		{2, 1.1, 15, 1800},
		{3, 0.9, 25, 1600},
		{4, 1.2, 10, 1700},
		{5, 1.0, 5, 2000},
		{6, 1.1, 30, 1400},
	}
	req, err := r.buildTirePressure(
		1,
		"Tire_pressure: стандартные условия",
		"Расчёт давления в шинах для легкового авто: асфальт, умеренный климат, средняя загрузка. Все заявки типа Tire_pressure.",
		entries,
	)
	if err != nil {
		return nil, err
	}
	return []TirePressure{req}, nil
}

// GetTirePressure возвращает заявку Tire_pressure по ID.
func (r *Repository) GetTirePressure(id int) (TirePressure, error) {
	requests, err := r.GetTirePressures()
	if err != nil {
		return TirePressure{}, err
	}
	for _, req := range requests {
		if req.ID == id {
			return req, nil
		}
	}
	return TirePressure{}, fmt.Errorf("заявка Tire_pressure не найдена")
}

// GetTirePressureForTire ищет заявку Tire_pressure, содержащую данную шину, и возвращает запись м-м.
func (r *Repository) GetTirePressureForTire(tireID int) (*TirePressureEntry, error) {
	requests, err := r.GetTirePressures()
	if err != nil {
		return nil, err
	}
	for _, req := range requests {
		for _, entry := range req.Entries {
			if entry.Tire.ID == tireID {
				return &entry, nil
			}
		}
	}
	return nil, fmt.Errorf("шина не найдена в заявках Tire_pressure")
}