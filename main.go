package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Структура, в которой есть лишь одно поле - Current с тэгом - json:current
// Current - анаонимная вложенная структура
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

// Обработка главной страницы
func index(w http.ResponseWriter, r *http.Request) {
	// Для отображения текущей погоды использовался бесплатный API - WeatherAPI (https://www.weatherapi.com/docs/)
	// Для
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Moscow"
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key is not set", http.StatusInternalServerError)
		return
	}

	u, err := url.Parse("http://api.weatherapi.com/v1/current.json") // URL-ссылка на API
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Декодирование запроса
	var wr WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Структура для вывода информации о городе и текущей температуре
	data := struct {
		City  string
		TempC float64
	}{
		City:  city,
		TempC: wr.Current.TempC,
	}

	//
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	handleFunc()
}
