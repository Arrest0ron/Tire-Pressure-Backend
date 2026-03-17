package serializer

import "metoda/internal/app/ds"

type TireJSON struct {
	TireID                   uint    `json:"tire_id"`
	TireTitle                string  `json:"tire_title"`
	TireMaterialCoefficient  float64 `json:"tire_material_coefficient"`
	TireThicknessCoefficient float64 `json:"tire_thickness_coefficient"`
	Description              string  `json:"description"`
	Photo                    string  `json:"photo"`
	Video                    string  `json:"video"`
	IsDelete                 bool    `json:"is_delete"`
}

func TireToJSON(t ds.Tire) TireJSON {
	return TireJSON{
		TireID:                   t.TireID,
		TireTitle:                t.TireTitle,
		TireMaterialCoefficient:  t.TireMaterialCoefficient,
		TireThicknessCoefficient: t.TireThicknessCoefficient,
		Description:              t.Description,
		Photo:                    t.Photo,
		Video:                    t.Video,
		IsDelete:                 t.IsDelete,
	}
}

func TireFromJSON(j TireJSON) ds.Tire {
	return ds.Tire{
		TireTitle:                j.TireTitle,
		TireMaterialCoefficient:  j.TireMaterialCoefficient,
		TireThicknessCoefficient: j.TireThicknessCoefficient,
		Description:              j.Description,
		Photo:                    j.Photo,
		Video:                    j.Video,
		IsDelete:                 j.IsDelete,
	}
}
