package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func (j StringArray) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *StringArray) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("database JSON value is not a string %v", value)
	}
	bytes := []byte(str)
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	return nil
}

func Upsert(db *gorm.DB, entry *Entry) error {
	fmt.Printf("upserting entry %v\n", entry)

	var existing Entry
	if err := db.Where("kanji = ?", entry.Kanji).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.Create(entry).Error
		}
		return err
	}
	existing.Kana = entry.Kana
	existing.English = entry.English
	existing.Labels = entry.Labels
	return db.Save(&existing).Error
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Entry{})
}
