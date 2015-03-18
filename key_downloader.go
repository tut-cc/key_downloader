package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type ZipFile struct {
	Name string
	Path string
}

func (z *ZipFile) Index(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("layout.html")

	if err != nil {
		panic(err)
	}

	err = tmple.Execute(w, z)

	if err != nil {
		panic(err)
	}
}

func (z *ZipFile) DownloadPage(w http.ResponseWriter, r *http.Request) {
	num := r.FormValue("num")
	i, err := strconv.Atoi(num)

	if err != nil {
		panic(err)
	}

	if i == 300 {
		file, err := ioutil.ReadFile(z.Path)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func main() {
	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("./layout"))))

	zip := ZipFile{"test", "test.zip"}
	http.HandleFunc("/", zip.Index)
	http.HandleFunc("/"+zip.Path, zip.DownloadPage)

	http.ListenAndServe(":8080", nil)
}
