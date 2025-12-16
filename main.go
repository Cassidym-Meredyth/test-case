package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Структура, в которой есть лишь одно поле - Current с тэгом - json:current
// Current - анонимная вложенная структура
// Ну то есть, структура будет иметь вид такого JSON файла:
//
//	{
//		"current": {
//			"temp_c": 4.2
//		}
//	}
//
// После этого, данная структура идет в json.NewDecoder(resp.Body).Decode(&wr)
// Значение объекта current попадает в поле wr.Current
type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// Структура данных, полученных после GET-запроса
type Data struct {
	City  string
	TempC float64
}

func getWeather(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("API key is not set")
	}

	u, err := url.Parse("http://api.weatherapi.com/v1/current.json") // URL-ссылка на API
	if err != nil {
		return 0, err
	}

	// Составление запроса
	q := u.Query()
	q.Set("key", apiKey) // Ввод API ключа
	q.Set("q", city)     // Ввод города для получения температуры
	q.Set("aqi", "no")   // Отключение aqi (пыльца)
	u.RawQuery = q.Encode()

	// Отправка GET Request
	resp, err := http.Get(u.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Декодирование запроса
	var wr WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
		return 0, err
	}

	return wr.Current.TempC, nil
}

// Обработка главной страницы
func index(w http.ResponseWriter, r *http.Request) {
	// query параметр. Нужен для получения температуры в необходимом городе
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Moscow"
	}

	// Получение температуры с помощью GET-запроса
	temp, err := getWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Структура для вывода информации о городе и текущей температуре
	data := Data{City: city, TempC: temp}

	// Загрузка и проверка HTML-шаблона
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Объединение самого HTML-шаблона с данными, полученными после GET-запроса
	if err := t.ExecuteTemplate(w, "index", data); err != nil {
		log.Println("template execute: ", err)
	}
}

// Обработка страниц
func handleFunc() {
	http.HandleFunc("/", index)       // Главная страница
	http.ListenAndServe(":8080", nil) // Порт для локального сервера
}

func main() {
	// Проверка файла .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	handleFunc()
}
