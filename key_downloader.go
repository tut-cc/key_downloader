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
	Fun  func(int) bool
}

type ZipList struct {
	List []ZipFile
}

func (z *ZipList) Index(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("index.html")

	if err != nil {
		panic(err)
	}

	err = tmple.Execute(w, z)

	if err != nil {
		panic(err)
	}
}

func (z *ZipFile) ItemIndex(w http.ResponseWriter, r *http.Request) {
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

	if z.Fun(i) {
		file, err := ioutil.ReadFile(z.Path)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	} else {
		http.Redirect(w, r, "/"+z.Name, http.StatusFound)
	}

}

func (z *ZipFile) DownloadPath() string {
	buf := z.Name
	buf = "/" + buf + ".zip"
	return buf
}

func check(i int) bool {
	return i == 300
}

func main() {
	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("./layout"))))

	zip := ZipFile{"test", "test.zip", check}
	zipList := ZipList{
		[]ZipFile{zip},
	}

	http.HandleFunc("/", zipList.Index)

	for i := 0; i < len(zipList.List); i++ {
		http.HandleFunc("/"+zipList.List[i].Name, zipList.List[i].ItemIndex)
		http.HandleFunc(zipList.List[i].DownloadPath(), zipList.List[i].DownloadPage)
	}

	http.ListenAndServe(":8080", nil)
}
