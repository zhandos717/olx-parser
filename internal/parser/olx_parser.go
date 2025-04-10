package parser

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"olx-parser/internal/database"
	"olx-parser/internal/models"
)

// OLXParser управляет процессом парсинга
type OLXParser struct {
	collector *colly.Collector
	dbManager *database.DatabaseManager
}

// NewOLXParser создает новый OLX парсер
func NewOLXParser(dbManager *database.DatabaseManager) *OLXParser {
	return &OLXParser{
		collector: colly.NewCollector(
			colly.AllowedDomains("www.olx.kz", "olx.kz"),
			colly.MaxDepth(2),
		),
		dbManager: dbManager,
	}
}

// Parse запускает процесс парсинга
func (p *OLXParser) Parse(searchQuery string) error {
	// Настройка коллектора
	p.setupCollectorCallbacks()

	// Начало парсинга со страницы поиска
	startURL := fmt.Sprintf("https://www.olx.kz/d/list/q-%s/", searchQuery)
	return p.collector.Visit(startURL)
}

func (p *OLXParser) setupCollectorCallbacks() {
	// Добавляем задержку и ограничения
	p.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*olx.kz*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// Callback для страницы со списком объявлений
	p.collector.OnHTML("a.css-rc5s2f", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://www.olx.kz" + link
		}
		e.Request.Visit(link)
	})

	// Callback для страницы отдельного объявления
	p.collector.OnHTML("div.css-1bred0n", func(e *colly.HTMLElement) {
		ad := models.Advertisement{
			Type: "olx",
			Link: e.Request.URL.String(),
		}

		// Извлечение заголовка
		if title := e.DOM.Find("h1.css-r9ym6f").First(); title.Length() > 0 {
			ad.Title = strings.TrimSpace(title.Text())
		}

		// Извлечение описания
		if desc := e.DOM.Find("div.css-1t507yq").First(); desc.Length() > 0 {
			ad.Description = strings.TrimSpace(desc.Text())
		}

		// Извлечение количества просмотров
		if views := e.DOM.Find("div.css-1mzs4y").First(); views.Length() > 0 {
			viewsText := strings.TrimSpace(views.Text())
			fmt.Sscanf(strings.ReplaceAll(viewsText, " ", ""), "%d", &ad.ViewsCount)
		}

		// Извлечение города и даты публикации
		infoElements := e.DOM.Find("div.css-v352pn")
		if infoElements.Length() > 1 {
			ad.City = strings.TrimSpace(infoElements.Eq(0).Text())

			dateStr := strings.TrimSpace(infoElements.Eq(1).Text())
			ad.PublicationDate = p.parseOLXDate(dateStr)
		}

		// Извлечение информации о пользователе
		if user := e.DOM.Find("div.css-1s3v3en").First(); user.Length() > 0 {
			ad.UserName = strings.TrimSpace(user.Text())
		}

		// Сохранение в базу данных
		if err := p.dbManager.SaveAdvertisement(ad); err != nil {
			log.Printf("Ошибка сохранения объявления: %v", err)
		}
	})
}

// parseOLXDate преобразует строку с датой OLX в time.Time
func (p *OLXParser) parseOLXDate(dateStr string) time.Time {
	now := time.Now()
	dateStr = strings.ToLower(dateStr)

	switch {
	case strings.Contains(dateStr, "сегодня"):
		return now
	case strings.Contains(dateStr, "вчера"):
		return now.AddDate(0, 0, -1)
	default:
		parsedDate, err := time.Parse("15:04, 2 Jan", dateStr)
		if err == nil {
			return parsedDate.AddDate(now.Year(), 0, 0)
		}
		parsedDate, err = time.Parse("2 Jan 15:04", dateStr)
		if err == nil {
			return parsedDate
		}
	}
	return time.Time{} // нулевая дата при ошибке
}
