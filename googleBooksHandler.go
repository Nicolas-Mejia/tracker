package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//takes the name of a book and returns the corresponding Google Books object
func getBookData(bookTitle string){
	bookTitle = strings.Replace(bookTitle, " ", "+",-1)

	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s", bookTitle)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil{
		fmt.Println(err.Error())
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

 	m := GoogleBook{}
	//m.Items[0].VolumeInfo.Title
 	err = json.NewDecoder(res.Body).Decode(&m)
 	if err != nil {
		fmt.Println(err.Error())
 	}	
	fmt.Println(m.Items[0].VolumeInfo.Title)
}
