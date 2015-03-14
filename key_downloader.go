package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type Page struct {
	Title string
}

func handler(w http.ResponseWriter, r *http.Request) {
	page := Page{"ZipDownloader"}
	tmple, err := template.ParseFiles("layout.html")

	if err != nil {
		panic(err)
	}

	err = tmple.Execute(w, page)

	if err != nil {
		panic(err)
	}

}

func zipDownload(w http.ResponseWriter, r *http.Request) {
	num := r.FormValue("num")
	i, err := strconv.Atoi(num)

	if err != nil {
		panic(err)
	}

	if i == 300 {
		file, err := ioutil.ReadFile("test.zip")
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/zip", zipDownload)
	http.ListenAndServe(":8080", nil)
}
