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

// Tire — услуга: шина с коэффициентами и медиа.
type Tire struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`            // название для каталога
	Description     string  `json:"description"`      // описание для каталога
	Photo           string  `json:"photo"`            // ключ фото в Minio
	Video           string  `json:"video"`            // ключ видео в Minio
	TireCoefficient float64 `json:"tire_coefficient"` // коэффициент шины (услуга)
}

// Tire_pressure — заявка: параметры расчёта + результаты.
type Tire_pressure struct {
	ID             int                 `json:"id"`              // единственная идентификация
	AirTemperature float64             `json:"air_temperature"` // температура воздуха (заявка)
	CarWeight      float64             `json:"car_weight"`      // вес авто (заявка)
	Entries        []TirePressureEntry `json:"entries"`         // связи м-м
	EntryCount     int                 `json:"entry_count"`     // [ДОБАВЛЕНО] количество записей
}

// TirePressureEntry — связь м-м: шина + покрытие → давление.
type TirePressureEntry struct {
	TireID       int     `json:"tire_id"`       // ссылка на Tire.ID
	CoatingCoeff float64 `json:"coating_coeff"` // коэффициент покрытия (м-м поле)
	Pressure     float64 `json:"pressure"`      // результат: давление в шине (хранится)
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

// buildTirePressure собирает заявку Tire_pressure из записей м-м.
// Давление (Pressure) берётся из хранимых данных, расчёт не производится.
func (r *Repository) buildTirePressure(id int, entries []TirePressureEntry) (Tire_pressure, error) {
	return Tire_pressure{
		ID:             id,
		AirTemperature: 20, // значение по умолчанию, если не указано в заявке
		CarWeight:      1500,
		Entries:        entries,
		EntryCount:     len(entries), // [ДОБАВЛЕНО] заполняем количество записей
	}, nil
}

// GetTirePressures возвращает все заявки на расчёт давления (Tire_pressure).
func (r *Repository) GetTirePressures() ([]Tire_pressure, error) {
	entries := []TirePressureEntry{
		{TireID: 1, CoatingCoeff: 1.0, Pressure: 44.0},
		{TireID: 2, CoatingCoeff: 1.1, Pressure: 48.3},
		{TireID: 3, CoatingCoeff: 0.9, Pressure: 39.6},
		{TireID: 4, CoatingCoeff: 1.2, Pressure: 52.4},
		{TireID: 5, CoatingCoeff: 1.0, Pressure: 48.0},
		{TireID: 6, CoatingCoeff: 1.1, Pressure: 44.2},
	}
	req, err := r.buildTirePressure(1, entries)
	if err != nil {
		return nil, err
	}
	return []Tire_pressure{req}, nil
}

// GetTirePressure возвращает заявку Tire_pressure по ID.
func (r *Repository) GetTirePressure(id int) (Tire_pressure, error) {
	requests, err := r.GetTirePressures()
	if err != nil {
		return Tire_pressure{}, err
	}
	for _, req := range requests {
		if req.ID == id {
			return req, nil
		}
	}
	return Tire_pressure{}, fmt.Errorf("заявка Tire_pressure не найдена")
}

// GetTirePressureForTire ищет заявку Tire_pressure, содержащую данную шину, и возвращает запись м-м.
func (r *Repository) GetTirePressureForTire(tireID int) (*TirePressureEntry, error) {
	requests, err := r.GetTirePressures()
	if err != nil {
		return nil, err
	}
	for _, req := range requests {
		for _, entry := range req.Entries {
			if entry.TireID == tireID {
				return &entry, nil
			}
		}
	}
	return nil, fmt.Errorf("шина не найдена в заявках Tire_pressure")
}
