package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Genre struct {
	GenreId    string `json:"genreId"`
	GenreName  string `json:"genreName"`
	GenreOrder string `json:"genreOrder"`
	URL        string `json:"url"`
	Count      int    `json:"count"`
	HasChild   bool   `json:"hasChild"`
}

type Data struct {
	Array []Genre `json:"array"`
}

func main() {
	// Open the file
	file, err := os.Open("test/counter/data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file's content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal JSON into struct
	var data Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Use a map to store unique counts by GenreId
	countsMap := make(map[string]int)

	for _, genre := range data.Array {
		if _, exists := countsMap[genre.GenreId]; !exists {
			// If the GenreId doesn't exist in the map, add it
			countsMap[genre.GenreId] = genre.Count
		}
	}

	// Calculate total count from unique GenreId counts
	totalCount := 0
	for _, count := range countsMap {
		totalCount += count
	}

	fmt.Println("Total count:", totalCount)
}

/*
Api Urls
https://books.rakuten.co.jp/json/genre/000
https://books.rakuten.co.jp/json/genre/001
https://books.rakuten.co.jp/json/genre/002
https://books.rakuten.co.jp/json/genre/003
https://books.rakuten.co.jp/json/genre/004
https://books.rakuten.co.jp/json/genre/005
https://books.rakuten.co.jp/json/genre/006
https://books.rakuten.co.jp/json/genre/007
https://books.rakuten.co.jp/json/genre/101
*/
