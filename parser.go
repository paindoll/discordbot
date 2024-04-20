package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerData struct {
	PlayersCount int    `json:"playersCount"`
	Available    bool   `json:"available"`
	Name         string `json:"name"`
}

func mergeJSONFiles(locations []string, urls []string, filepath string, interval time.Duration) {
	for {
		var mergedData []ServerData

		for index, url := range urls {
			siteName := locations[index]
			response, err := http.Get(url)
			if err != nil {
				log.Printf("Ошибка при получении данных с %s: %v\n", siteName, err)
				continue
			}
			// defer response.Body.Close() // defers in this infinite loop will never run
			var data ServerData
			if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
				log.Printf("Ошибка при декодировании данных с %s: %v\n", siteName, err)
				continue
			}

			data.Name = siteName

			mergedData = append(mergedData, data)
		}

		jsonData, err := json.MarshalIndent(mergedData, "", "    ")
		if err != nil {
			log.Printf("Ошибка при маршалинге данных: %v\n", err)
			continue
		}

		if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
			log.Printf("Ошибка при записи в файл: %v\n", err)
			continue
		}

		log.Println("Данные обновлены успешно.")
		time.Sleep(interval)
	}
}

func parser() {
	locations := []string{
		"New York",
		"Detroit",
		"Chicago",
		"San Francisco",
		"Atlanta",
		"San Diego",
		"Los Angeles",
		"Miami",
		"Las Vegas",
		"Washington",
	}

	urls := []string{
		"https://api.alt-mp.com/servers/SvgjrFK",
		"https://api.alt-mp.com/servers/y8vgQbo",
		"https://api.alt-mp.com/servers/0kIhphA",
		"https://api.alt-mp.com/servers/qL4KKep",
		"https://api.alt-mp.com/servers/AwCvQP2",
		"https://api.alt-mp.com/servers/0FcLWkD",
		"https://api.alt-mp.com/servers/y5mEO6p",
		"https://api.alt-mp.com/servers/BnVfDtR",
		"https://api.alt-mp.com/servers/YMRctiN",
		"https://api.alt-mp.com/servers/GK6qeYO",
	}

	filepath := "./servers.json"
	interval := 10 * time.Second

	mergeJSONFiles(locations, urls, filepath, interval)
}
