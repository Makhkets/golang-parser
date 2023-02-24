package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	URL        = "http://m-classic.ru/"
	CATEGORIES = []string{
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D0%B8/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%BA%D1%83%D1%85%D0%BD%D0%B8/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%B2%D0%B8%D1%82%D1%80%D0%B8%D0%BD%D1%8B-%D1%81-%D1%82%D0%B2-%D1%82%D1%83%D0%BC%D0%B1%D0%B0%D0%BC%D0%B8/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D1%81%D1%82%D0%BE%D0%BB%D0%BE%D0%B2%D1%8B%D0%B5/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%BC%D1%8F%D0%B3%D0%BA%D0%B0%D1%8F-%D0%BC%D0%B5%D0%B1%D0%B5%D0%BB%D1%8C/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%B3%D0%BE%D1%81%D1%82%D0%B8%D0%BD%D1%8B%D0%B5/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%BC%D0%B0%D0%BB%D1%8B%D0%B5-%D1%84%D0%BE%D1%80%D0%BC%D1%8B/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%B2%D1%8B%D1%81%D1%82%D0%B0%D0%B2.-%D0%BE%D0%B1%D1%80%D0%B0%D0%B7%D1%86%D1%8B-%D1%82%D1%8E%D0%BC%D0%B5%D0%BD%D1%8C/",
		"http://m-classic.ru/%D0%BA%D0%B0%D1%82%D0%B0%D0%BB%D0%BE%D0%B3/%D0%BB%D0%B8%D0%BA%D0%B2%D0%B8%D0%B4%D0%B0%D1%86%D0%B8%D1%8F/",
	}
)

// COLORS
const (
	Red    string = "\033[31m"
	Green  string = "\033[32m"
	Yellow string = "\033[33m"
	Blue   string = "\033[34m"
	Purple string = "\033[35m"
	Cyan   string = "\033[36m"
	Reset  string = "\033[0m"
)

func main() {
	excelInit()
	start()
}

// Запускаем парсер
func start() {
	var wg sync.WaitGroup
	for _, category := range CATEGORIES {
		wg.Add(1)
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
		go func() {
			defer wg.Done()
			parse(category)
		}()
	}

	wg.Wait()
}

// Запускает парсинг каждой страницы
func parse(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	// Парсим номер последней страницы (PAGINATION)
	var lastPageNumber int = 0
	element := doc.Find(".mse2_pagination .pagination ul li:last-of-type")
	lastPage, exists := element.Find("a").Attr("href")
	if exists {
		lastPageNumber, err = strconv.Atoi(strings.Split(lastPage, "page=")[1])
		if err != nil {
			log.Printf("Не смог спарсить номер последней страницы при запросе на url: %s", url)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i <= lastPageNumber; i++ {
		wg.Add(1)
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
		go func() {
			defer wg.Done()
			uri := url + "?page=2"
			parsePage(uri)
		}()
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
	}
	wg.Wait()
	log.Println("Parse Done")

}

// Парсит страницу с карточками
func parsePage(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	var wg sync.WaitGroup
	doc.Find(".item-info").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		wg.Add(1)
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		go func() {
			defer wg.Done()
			parseCardPage(URL + href)
		}()
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)

	})
	wg.Wait()
}

// Парсит саму карточку
func parseCardPage(cardUrl string) {
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	res, err := http.Get(cardUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	breadcrumb := doc.Find(".breadcrumbs")

	category := breadcrumb.Find("li").Eq(1).Find("a").Text()
	description := doc.Find(".desc.std").Text()
	title := doc.Find(".item-title").Eq(0).Find("h1").Text()
	price := doc.Find(".item-price .price").Eq(0).Text()
	images := ""
	producer := ""
	description2 := ""

	//Search for the img tag
	elements, _ := doc.Find("#msGallery").Html()

	// Создаем регулярное выражение для поиска ссылок на картинки с расширениями .png, .jpg и .jpeg
	re := regexp.MustCompile(`(?i)<a.*?href=["']([^"']+\.(?:png|jpe?g))["'].*?>`)

	// Находим все ссылки на картинки внутри тегов a с классом fotorama
	matches := re.FindAllStringSubmatch(elements, -1)

	// Сохраняем найденные ссылки
	for _, match := range matches {
		images += strings.TrimSpace(match[1]) + "\n"
	}

	// Находим элемент с информацией о производителе
	var manufacturer *goquery.Selection
	doc.Find(".desc2 .form-group").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("label").Text())
		description2 += label + " "
		description2 += s.Find("strong").Text() + "\n"
		if label == "Производитель:" {
			manufacturer = s.Find("strong")
			return
		}
	})

	if manufacturer == nil {
		log.Println("Производитель не найден: ", cardUrl)
	} else {
		// Извлекаем текст из элемента с информацией о производителе
		producer = strings.TrimSpace(manufacturer.Text())
		fmt.Println(producer)
	}

	data := Data{
		Category:     category,
		Title:        title,
		Price:        price,
		Producer:     producer,
		Description:  description,
		Description2: description2,
		Url:          cardUrl,
		Images:       images,
	}
	data.addDataToExcel()
}
