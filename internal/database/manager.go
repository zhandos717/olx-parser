package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"olx-parser/internal/models"
)

// DatabaseManager управляет операциями с базой данных
type DatabaseManager struct {
	db *sql.DB
}

// NewDatabaseManager создает новый экземпляр менеджера базы данных
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

// InitDatabase создает и настраивает SQLite базу данных
func (dm *DatabaseManager) InitDatabase() error {
	var err error
	dm.db, err = sql.Open("sqlite3", "./advertisements.db")
	if err != nil {
		return err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS advertisements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT UNIQUE,
		description TEXT,
		type TEXT,
		views_count INTEGER,
		city TEXT,
		publication_date DATETIME,
		user_name TEXT,
		user_phone TEXT,
		user_online_date DATETIME,
		user_id TEXT
	)`

	_, err = dm.db.Exec(createTableSQL)
	return err
}

// SaveAdvertisement сохраняет parsed объявление в базу данных
func (dm *DatabaseManager) SaveAdvertisement(ad models.Advertisement) error {
	if dm.db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	insertSQL := `
	INSERT OR REPLACE INTO advertisements 
	(title, link, description, type, views_count, city, 
	publication_date, user_name, user_phone, user_online_date, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := dm.db.Exec(
		insertSQL,
		ad.Title,
		ad.Link,
		ad.Description,
		ad.Type,
		ad.ViewsCount,
		ad.City,
		ad.PublicationDate,
		ad.UserName,
		ad.UserPhone,
		ad.UserOnlineDate,
		ad.UserID,
	)
	return err
}

// Close закрывает соединение с базой данных
func (dm *DatabaseManager) Close() {
	if dm.db != nil {
		if err := dm.db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}
}

// GetAdvertisementsCount возвращает количество сохраненных объявлений
func (dm *DatabaseManager) GetAdvertisementsCount() (int, error) {
	if dm.db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	var count int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM advertisements").Scan(&count)
	return count, err
}
