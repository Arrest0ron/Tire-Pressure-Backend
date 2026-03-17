package serializer

import (
	"metoda/internal/app/ds"
	"time"
)

// TirePressureListJSON — used for the list endpoint.
type TirePressureListJSON struct {
	TirePressureID   uint       `json:"tire_pressure_id"`
	Status           string     `json:"status"`
	DateCreate       time.Time  `json:"date_create"`
	DateFormed       *time.Time `json:"date_formed"`
	DateCompleted    *time.Time `json:"date_completed"`
	CreatorLogin     string     `json:"creator_login"`
	ModeratorLogin   *string    `json:"moderator_login"`
	AirTemperature   float64    `json:"air_temperature"`
	CarWeight        float64    `json:"car_weight"`
	TireEntriesCount int        `json:"tire_entries_count"`
}

// TirePressureDetailJSON — used for the single-record endpoint (includes tires/entries).
type TirePressureDetailJSON struct {
	TirePressureID uint                        `json:"tire_pressure_id"`
	Status         string                      `json:"status"`
	DateCreate     time.Time                   `json:"date_create"`
	DateFormed     *time.Time                  `json:"date_formed"`
	DateCompleted  *time.Time                  `json:"date_completed"`
	CreatorLogin   string                      `json:"creator_login"`
	ModeratorLogin *string                     `json:"moderator_login"`
	AirTemperature float64                     `json:"air_temperature"`
	CarWeight      float64                     `json:"car_weight"`
	Entries        []TirePressureEntryViewJSON `json:"entries"`
}

// TirePressureEntryViewJSON — one tire inside a tire-pressure detail.
type TirePressureEntryViewJSON struct {
	ID                       uint    `json:"id"`
	TireID                   uint    `json:"tire_id"`
	TireTitle                string  `json:"tire_title"`
	TireMaterialCoefficient  float64 `json:"tire_material_coefficient"`
	TireThicknessCoefficient float64 `json:"tire_thickness_coefficient"`
	Photo                    string  `json:"photo"`
	CoatingCoefficient       float64 `json:"coating_coefficient"`
	Pressure                 float64 `json:"pressure"`
}

// TirePressureUpdateJSON — поля, которые создатель может менять через PUT (черновик).
type TirePressureUpdateJSON struct {
	AirTemperature float64 `json:"air_temperature"`
	CarWeight      float64 `json:"car_weight"`
}

// FinishJSON — status sent by moderator to finish or reject.
type FinishJSON struct {
	Status string `json:"status"`
}

// CartJSON — cart icon response.
type CartJSON struct {
	TirePressureID uint  `json:"tire_pressure_id"`
	TiresCount     int64 `json:"tires_count"`
}

func TirePressureToListJSON(t ds.TirePressure, creatorLogin string, moderatorLogin string, entriesCount int) TirePressureListJSON {
	var dateFormed, dateCompleted *time.Time
	if t.DateFormed.Valid {
		dateFormed = &t.DateFormed.Time
	}
	if t.DateCompleted.Valid {
		dateCompleted = &t.DateCompleted.Time
	}
	var modLogin *string
	if moderatorLogin != "" {
		modLogin = &moderatorLogin
	}
	return TirePressureListJSON{
		TirePressureID:   t.TirePressureID,
		Status:           t.Status,
		DateCreate:       t.DateCreate,
		DateFormed:       dateFormed,
		DateCompleted:    dateCompleted,
		CreatorLogin:     creatorLogin,
		ModeratorLogin:   modLogin,
		AirTemperature:   t.AirTemperature,
		CarWeight:        t.CarWeight,
		TireEntriesCount: entriesCount,
	}
}
