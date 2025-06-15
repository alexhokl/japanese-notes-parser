package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	var existing Entry

	if entry.Kanji == "" {
		fmt.Printf("upserting entry %s %v %v\n", entry.Kana, entry.English, entry.Labels)
		if err := db.Where("kana = ?", entry.Kana).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return db.Create(entry).Error
			}
			return err
		}
		existing.Kana = entry.Kana
	} else {
		fmt.Printf("upserting entry %s %s %v %v\n", entry.Kanji, entry.Kana, entry.English, entry.Labels)
		if err := db.Where("kanji = ?", entry.Kanji).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return db.Create(entry).Error
			}
			return err
		}
	}

	existing.English = entry.English
	existing.Labels = entry.Labels
	return db.Save(&existing).Error
}

func ListPartsOfSpeech(db *gorm.DB) ([]string, error) {
	var entries []Entry
	if err := db.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("failed to list parts of speech: %w", err)
	}

	partOfSpeechMap := make(map[string]struct{})
	for _, entry := range entries {
		for _, label := range entry.Labels {
			partOfSpeechMap[label] = struct{}{}
		}
	}

	var partOfSpeechList []string
	for partOfSpeech := range partOfSpeechMap {
		partOfSpeechList = append(partOfSpeechList, partOfSpeech)
	}

	return partOfSpeechList, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Entry{})
}
