package main

import (
	"log"

	"olx-parser/internal/database"
	"olx-parser/internal/parser"
)

func main() {
	// Создаем менеджер базы данных
	dbManager := database.NewDatabaseManager()

	// Инициализируем базу данных
	err := dbManager.InitDatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer dbManager.Close()

	// Создаем парсер OLX
	olxParser := parser.NewOLXParser(dbManager)

	// Примеры поисковых запросов для строительных товаров
	searchQueries := []string{
		"строительные+инструменты",
		"аренда+строительного+оборудования",
		"прокат+строительной+техники",
	}

	// Парсим каждый запрос
	for _, query := range searchQueries {
		err := olxParser.Parse(query)
		if err != nil {
			log.Printf("Ошибка парсинга для запроса %s: %v", query, err)
		}
	}

	// Проверяем количество сохраненных объявлений
	count, err := dbManager.GetAdvertisementsCount()
	if err != nil {
		log.Printf("Ошибка подсчета объявлений: %v", err)
	} else {
		log.Printf("Всего сохранено объявлений: %d", count)
	}
}
