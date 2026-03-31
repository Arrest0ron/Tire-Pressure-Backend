package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"metoda/internal/app/ds"
	"metoda/internal/app/serializer"
)

// ─── helpers ────────────────────────────────────────────────────────────────

func (r *Repository) ClearDraftTirePressuresOnStartup() {
	r.db.Exec("UPDATE tire_pressures SET status = ? WHERE status = ?", ds.StatusDeleted, ds.StatusDraft)
}

func (r *Repository) GetCreatorLogin(creatorID uint) string {
	var u ds.Users
	r.db.Where("id = ?", creatorID).First(&u)
	return u.Login
}

func (r *Repository) GetModeratorLogin(moderatorID *uint) string {
	if moderatorID == nil {
		return ""
	}
	var u ds.Users
	r.db.Where("id = ?", *moderatorID).First(&u)
	return u.Login
}

// GetTireEntriesCount — число записей м-м, в которых рассчитываемое поле результата (давление) не пустое.
func (r *Repository) GetTireEntriesCount(tirePressureID uint) int {
	var count int64
	r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ?", tirePressureID).
		Where("pressure IS NOT NULL AND pressure > 0").
		Count(&count)
	return int(count)
}

// ─── HTML-layer methods (kept for lab2 compatibility) ───────────────────────

func (r *Repository) GetDraftTirePressure(creatorID uint) (*ds.TirePressure, error) {
	var t ds.TirePressure
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *Repository) GetTirePressureWithEntries(tirePressureID uint) (*ds.TirePressure, []TirePressureEntryView, error) {
	var t ds.TirePressure
	err := r.db.First(&t, tirePressureID).Error
	if err != nil {
		return nil, nil, err
	}

	var items []ds.TirePressureEntry
	err = r.db.Where("tire_pressure_id = ?", tirePressureID).
		Preload("Tire").
		Order("id").
		Find(&items).Error
	if err != nil {
		return nil, nil, err
	}

	// ✅ ОТЛАДКА: выводим в консоль
	logrus.Infof("GetTirePressureWithEntries: Found %d items", len(items))
	for i, item := range items {
		logrus.Infof("  Item %d: TireID=%d, TireTitle=%s, Photo=%s",
			i, item.TireID, item.Tire.TireTitle, item.Tire.Photo)
	}

	var views []TirePressureEntryView
	for _, item := range items {
		views = append(views, TirePressureEntryView{
			ID:                       item.ID,
			TireTitle:                item.Tire.TireTitle,
			TireMaterialCoefficient:  item.Tire.TireMaterialCoefficient,
			TireThicknessCoefficient: item.Tire.TireThicknessCoefficient,
			Photo:                    item.Tire.Photo,
			CoatingCoefficient:       item.CoatingCoefficient,
			Pressure:                 item.Pressure,
		})
	}
	return &t, views, nil
}

type TirePressureEntryView struct {
	ID                       uint
	TireTitle                string
	TireMaterialCoefficient  float64
	TireThicknessCoefficient float64
	Photo                    string
	CoatingCoefficient       float64
	Pressure                 float64
}

func (r *Repository) GetCalculatedPressures(tirePressureID uint) []float64 {
	var pressures []float64
	err := r.db.Raw(`
		SELECT pressure
		FROM tire_pressure_entries
		WHERE tire_pressure_id = ?
		  AND pressure IS NOT NULL AND pressure > 0
	`, tirePressureID).Scan(&pressures).Error
	if err != nil {
		logrus.Errorf("GetCalculatedPressures: %v", err)
	}
	return pressures
}

func (r *Repository) AddTireToTirePressure(tireID uint, creatorID uint) error {
	var t ds.TirePressure
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&t).Error
	if err != nil {
		t = ds.TirePressure{
			Status:     ds.StatusDraft,
			DateCreate: time.Now(),
			CreatorID:  creatorID,
		}
		if err := r.db.Create(&t).Error; err != nil {
			return err
		}
	}
	tp := ds.TirePressureEntry{
		TirePressureID:     t.TirePressureID,
		TireID:             tireID,
		CoatingCoefficient: 1.0,
	}
	return r.db.Create(&tp).Error
}

func (r *Repository) DeleteTirePressureBySQL(tirePressureID uint) error {
	return r.db.Exec("UPDATE tire_pressures SET status = '"+ds.StatusDeleted+"' WHERE tire_pressure_id = ?", tirePressureID).Error
}

// UpdateTirePressureItem обновляет коэффициент типа покрытия
// для одной строки связи заявки с шиной (HTML-форма на странице корзины).
func (r *Repository) UpdateTirePressureItem(itemID uint, coatingCoefficient float64) error {
	var item ds.TirePressureEntry
	if err := r.db.First(&item, itemID).Error; err != nil {
		return err
	}

	if coatingCoefficient <= 0 {
		coatingCoefficient = 1.0
	}

	updates := map[string]interface{}{
		"coating_coefficient": coatingCoefficient,
	}

	return r.db.Model(&item).Updates(updates).Error
}

// FormTirePressure рассчитывает давление и меняет статус на "сформирован"
func (r *Repository) FormTirePressure(tirePressureID uint) error {
	// ✅ Сначала рассчитываем давление для всех шин
	var entries []ds.TirePressureEntry
	if err := r.db.Where("tire_pressure_id = ?", tirePressureID).Find(&entries).Error; err != nil {
		return err
	}

	var tp ds.TirePressure
	if err := r.db.First(&tp, tirePressureID).Error; err != nil {
		return err
	}

	for _, entry := range entries {
		var tire ds.Tire
		if err := r.db.First(&tire, entry.TireID).Error; err != nil {
			continue
		}

		// Формула расчёта
		basePressure := 200.0
		tempCorrection := (tp.AirTemperature - 20.0) * 2.0
		weightCorrection := (tp.CarWeight - 1500.0) * 0.1

		pressure := (basePressure + tempCorrection + weightCorrection) *
			tire.TireMaterialCoefficient *
			tire.TireThicknessCoefficient *
			entry.CoatingCoefficient

		r.db.Model(&entry).Updates(map[string]interface{}{
			"pressure": pressure,
		})
	}

	// ✅ Затем меняем статус
	now := time.Now()
	updates := map[string]interface{}{
		"status":      ds.StatusFormed,
		"date_formed": now,
	}
	return r.db.Model(&ds.TirePressure{}).Where("tire_pressure_id = ?", tirePressureID).Updates(updates).Error
}

func (r *Repository) GetCartCount(creatorID uint) int64 {
	var tirePressureID uint
	var count int64
	err := r.db.Model(&ds.TirePressure{}).Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).Select("tire_pressure_id").First(&tirePressureID).Error
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.TirePressureEntry{}).Where("tire_pressure_id = ?", tirePressureID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records:", err)
	}
	return count
}

func (r *Repository) GetDraftTirePressureID(creatorID uint) uint {
	var tirePressureID uint
	err := r.db.Model(&ds.TirePressure{}).Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).Select("tire_pressure_id").First(&tirePressureID).Error
	if err != nil {
		return 0
	}
	return tirePressureID
}

// ─── API methods ─────────────────────────────────────────────────────────────

// GetTirePressureByID returns a non-deleted tire-pressure by ID.
func (r *Repository) GetTirePressureByID(id int) (ds.TirePressure, error) {
	var t ds.TirePressure
	err := r.db.Where("tire_pressure_id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.TirePressure{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.TirePressure{}, err
	}

	if t.Status == ds.StatusDeleted {
		return ds.TirePressure{}, fmt.Errorf("%w: заявка удалена", ErrNotFound)
	}

	return t, nil
}

func (r *Repository) GetAllTirePressures(from, to time.Time, status string) ([]ds.TirePressure, error) {
	var list []ds.TirePressure
	sub := r.db.Where("status != ? AND status != ?", ds.StatusDeleted, ds.StatusDraft)
	if !from.IsZero() {
		sub = sub.Where("date_formed >= ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_formed <= ?", to.Add(24*time.Hour))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("tire_pressure_id").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetTirePressureEntriesAPI returns the m-m records for a tire-pressure with tire data.
func (r *Repository) GetTirePressureEntriesAPI(id uint) ([]serializer.TirePressureEntryViewJSON, error) {
	var items []ds.TirePressureEntry
	err := r.db.Where("tire_pressure_id = ?", id).Preload("Tire").Order("id").Find(&items).Error
	if err != nil {
		return nil, err
	}
	views := make([]serializer.TirePressureEntryViewJSON, 0, len(items))
	for _, item := range items {
		views = append(views, serializer.TirePressureEntryViewJSON{
			ID:                       item.ID,
			TireID:                   item.TireID,
			TireTitle:                item.Tire.TireTitle,
			TireMaterialCoefficient:  item.Tire.TireMaterialCoefficient,
			TireThicknessCoefficient: item.Tire.TireThicknessCoefficient,
			Photo:                    item.Tire.Photo,
			CoatingCoefficient:       item.CoatingCoefficient,
			Pressure:                 item.Pressure,
		})
	}
	return views, nil
}

// GetCartInfo returns the draft tire-pressure ID and tires count for the given user.
func (r *Repository) GetCartInfo(creatorID uint) (uint, int64, error) {
	var t ds.TirePressure
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	var count int64
	r.db.Model(&ds.TirePressureEntry{}).Where("tire_pressure_id = ?", t.TirePressureID).Count(&count)
	return t.TirePressureID, count, nil
}

// UpdateTirePressureFields обновляет редактируемые поля заявки (air_temperature, car_weight). Только для черновика и только создатель.
func (r *Repository) UpdateTirePressureFields(id int, j serializer.TirePressureUpdateJSON) (ds.TirePressure, error) {
	t, err := r.GetTirePressureByID(id)
	if err != nil {
		return ds.TirePressure{}, err
	}
	if t.Status != ds.StatusDraft {
		return ds.TirePressure{}, fmt.Errorf("%w: можно менять только черновик", ErrNotAllowed)
	}
	// ✅ SINGLETON: замена r.userID → GetUserID()
	if int(t.CreatorID) != GetUserID() {
		return ds.TirePressure{}, fmt.Errorf("%w: только создатель может редактировать заявку", ErrNotAllowed)
	}

	updates := map[string]interface{}{}
	if j.AirTemperature != 0 {
		updates["air_temperature"] = j.AirTemperature
	}
	if j.CarWeight != 0 {
		updates["car_weight"] = j.CarWeight
	}

	if len(updates) > 0 {
		if err := r.db.Model(&t).Updates(updates).Error; err != nil {
			return ds.TirePressure{}, err
		}
	}

	// Re-fetch to get updated values
	return r.GetTirePressureByID(id)
}

// ✅ FormTirePressureAPI рассчитывает давление и меняет статус на "сформирован"
func (r *Repository) FormTirePressureAPI(id int) (ds.TirePressure, error) {
	t, err := r.GetTirePressureByID(id)
	if err != nil {
		return ds.TirePressure{}, err
	}
	if t.Status != ds.StatusDraft {
		return ds.TirePressure{}, fmt.Errorf("нельзя сформировать заявку со статусом %s", t.Status)
	}

	// ✅ SINGLETON: замена r.userID → GetUserID()
	if int(t.CreatorID) != GetUserID() {
		return ds.TirePressure{}, fmt.Errorf("%w: только создатель может сформировать заявку", ErrNotAllowed)
	}

	// ✅ Проверка: хотя бы одна шина
	var count int64
	r.db.Model(&ds.TirePressureEntry{}).Where("tire_pressure_id = ?", t.TirePressureID).Count(&count)
	if count == 0 {
		return ds.TirePressure{}, fmt.Errorf("нельзя сформировать пустую заявку: добавьте хотя бы одну шину")
	}

	// ✅ РАССЧИТЫВАЕМ ДАВЛЕНИЕ ДЛЯ КАЖДОЙ ШИНЫ
	var entries []ds.TirePressureEntry
	if err := r.db.Where("tire_pressure_id = ?", t.TirePressureID).Find(&entries).Error; err != nil {
		return ds.TirePressure{}, err
	}

	for _, entry := range entries {
		var tire ds.Tire
		if err := r.db.First(&tire, entry.TireID).Error; err != nil {
			continue
		}

		// ✅ Формула расчёта
		basePressure := 200.0
		tempCorrection := (t.AirTemperature - 20.0) * 2.0
		weightCorrection := (t.CarWeight - 1500.0) * 0.1

		pressure := (basePressure + tempCorrection + weightCorrection) *
			tire.TireMaterialCoefficient *
			tire.TireThicknessCoefficient *
			entry.CoatingCoefficient

		// ✅ СОХРАНЯЕМ ДАВЛЕНИЕ В БД
		if err := r.db.Model(&entry).Update("pressure", pressure).Error; err != nil {
			return ds.TirePressure{}, err
		}
	}

	// ✅ Меняем статус
	now := time.Now()
	updates := map[string]interface{}{
		"status":      ds.StatusFormed,
		"date_formed": now,
	}

	if err := r.db.Model(&t).Updates(updates).Error; err != nil {
		return ds.TirePressure{}, err
	}

	return r.GetTirePressureByID(id)
}

// FinishTirePressureAPI lets a moderator set status to "завершён" or "отклонён".
func (r *Repository) FinishTirePressureAPI(id int, status string) (ds.TirePressure, error) {
	if status != ds.StatusCompleted && status != ds.StatusRejected {
		return ds.TirePressure{}, fmt.Errorf("недопустимый статус: ожидается '%s' или '%s'", ds.StatusCompleted, ds.StatusRejected)
	}

	// ✅ SINGLETON: замена r.userID → GetUserID()
	moderator, err := r.GetUserByID(GetUserID())
	if err != nil {
		return ds.TirePressure{}, err
	}
	if !moderator.IsModerator {
		return ds.TirePressure{}, fmt.Errorf("%w: только модератор может завершить или отклонить заявку", ErrNotAllowed)
	}

	t, err := r.GetTirePressureByID(id)
	if err != nil {
		return ds.TirePressure{}, err
	}
	if t.Status != ds.StatusFormed {
		return ds.TirePressure{}, fmt.Errorf("завершить или отклонить можно только сформированную заявку; текущий статус — %s (сначала PUT .../form)", t.Status)
	}

	now := time.Now()
	moderatorID := uint(moderator.ID)
	if err := r.db.Model(&t).Updates(map[string]interface{}{
		"status":         status,
		"date_completed": now,
		"moderator_id":   moderatorID,
	}).Error; err != nil {
		return ds.TirePressure{}, err
	}

	return r.GetTirePressureByID(id)
}

// DeleteTirePressureAPI soft-deletes a draft tire-pressure (creator only).
func (r *Repository) DeleteTirePressureAPI(id int) error {
	t, err := r.GetTirePressureByID(id)
	if err != nil {
		return err
	}
	if t.Status != ds.StatusDraft {
		return fmt.Errorf("%w: только черновик может быть удалён создателем", ErrNotAllowed)
	}
	// ✅ SINGLETON: замена r.userID → GetUserID()
	if int(t.CreatorID) != GetUserID() {
		return fmt.Errorf("%w: только создатель может удалить заявку", ErrNotAllowed)
	}
	return r.db.Model(&t).Update("status", ds.StatusDeleted).Error
}

// AddTireToCartAPI adds a tire to the user's draft (creates draft if needed).
func (r *Repository) AddTireToCartAPI(tireID uint, creatorID uint) (ds.TirePressure, bool, error) {
	// Verify tire exists and is not deleted
	var t ds.Tire
	if err := r.db.Where("tire_id = ? AND is_delete = false", tireID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.TirePressure{}, false, fmt.Errorf("%w: шина с id %d", ErrNotFound, tireID)
		}
		return ds.TirePressure{}, false, err
	}

	// Get or create draft
	var tp ds.TirePressure
	created := false
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&tp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tp = ds.TirePressure{
			Status:     ds.StatusDraft,
			DateCreate: time.Now(),
			CreatorID:  creatorID,
		}
		if err := r.db.Create(&tp).Error; err != nil {
			return ds.TirePressure{}, false, err
		}
		created = true
	} else if err != nil {
		return ds.TirePressure{}, false, err
	}

	// Check for duplicate
	var existing ds.TirePressureEntry
	res := r.db.Where("tire_pressure_id = ? AND tire_id = ?", tp.TirePressureID, tireID).First(&existing)
	if res.Error == nil {
		return ds.TirePressure{}, false, fmt.Errorf("%w: шина %d уже добавлена в заявку %d", ErrAlreadyExists, tireID, tp.TirePressureID)
	}

	tpe := ds.TirePressureEntry{
		TirePressureID:     tp.TirePressureID,
		TireID:             tireID,
		CoatingCoefficient: 1.0,
	}
	if err := r.db.Create(&tpe).Error; err != nil {
		return ds.TirePressure{}, false, err
	}
	return tp, created, nil
}

// DeleteTireFromCartAPI removes a tire from a tire-pressure (must be draft).
func (r *Repository) DeleteTireFromCartAPI(tireID, tirePressureID int) (ds.TirePressure, error) {
	t, err := r.GetTirePressureByID(tirePressureID)
	if err != nil {
		return ds.TirePressure{}, err
	}
	if t.Status != ds.StatusDraft {
		return ds.TirePressure{}, fmt.Errorf("%w: нельзя изменить не черновик", ErrNotAllowed)
	}
	err = r.db.Where("tire_id = ? AND tire_pressure_id = ?", tireID, tirePressureID).Delete(&ds.TirePressureEntry{}).Error
	if err != nil {
		return ds.TirePressure{}, err
	}
	return t, nil
}

// UpdateTireInCartAPI edits coating_coefficient in m-m.
func (r *Repository) UpdateTireInCartAPI(tireID, tirePressureID int, j serializer.TirePressureEntryUpdateJSON) (ds.TirePressureEntry, error) {
	t, err := r.GetTirePressureByID(tirePressureID)
	if err != nil {
		return ds.TirePressureEntry{}, err
	}
	if t.Status != ds.StatusDraft {
		return ds.TirePressureEntry{}, fmt.Errorf("%w: нельзя изменить не черновик", ErrNotAllowed)
	}

	var item ds.TirePressureEntry
	err = r.db.Where("tire_id = ? AND tire_pressure_id = ?", tireID, tirePressureID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.TirePressureEntry{}, fmt.Errorf("%w: связь не найдена", ErrNotFound)
		}
		return ds.TirePressureEntry{}, err
	}

	updates := map[string]interface{}{}
	if j.CoatingCoefficient > 0 {
		updates["coating_coefficient"] = j.CoatingCoefficient
	}

	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return ds.TirePressureEntry{}, err
	}

	// Re-fetch
	r.db.Where("tire_id = ? AND tire_pressure_id = ?", tireID, tirePressureID).First(&item)
	return item, nil
}

// GetTirePressureEntryByID получает запись м-м по ID
func (r *Repository) GetTirePressureEntryByID(id uint) (ds.TirePressureEntry, error) {
	var item ds.TirePressureEntry
	err := r.db.First(&item, id).Error
	return item, err
}
