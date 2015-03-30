package main

import (
	"fmt"
	"github.com/kimitoboku/go-PollarRho"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type ZipFile struct {
	Name string
	Path string
	Fun  func(string) bool
}

type ZipList struct {
	List []ZipFile
}

func (z *ZipList) Index(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("index.html")

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	err = tmple.Execute(w, z)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}

func (z *ZipFile) ItemIndex(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("layout.html")

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	err = tmple.Execute(w, z)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}

func (z *ZipFile) DownloadPage(w http.ResponseWriter, r *http.Request) {
	num := r.FormValue("num")

	if z.Fun(num) {
		file, err := ioutil.ReadFile(z.Path)
		if err != nil {
			fmt.Errorf(err.Error())
			return
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

func check(num string) bool {
	i, err := strconv.Atoi(num)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	return i == 300
}

func checkPrim(num string) bool {
	i, err := strconv.Atoi(num)
	if err != nil {
		fmt.Errorf(err.Error())
		return false
	}

	list := pollarrho.Factor(i)
	if len(list) == 3 {
		for _, num := range list {
			if num > 10000 || num < 999 {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func main() {
	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("./layout"))))

	zip := ZipFile{"test", "test.zip", checkPrim}
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
