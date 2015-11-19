package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	logger *log.Logger
)

type ZipFile struct {
	Name string
	Path string
	Key  *list.List
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

	if z.check(num) {
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

func (z *ZipFile) check(key string) bool {

	for e := z.Key.Front(); e != nil; e = e.Next() {
		if key == fmt.Sprintf("%s", e.Value) {
			logger.Printf("Succces full : %s", key)
			return true
		}
	}
	logger.Printf("Fail : %s", key)
	return false
}

func logSetUp() {
	f, err := os.OpenFile("/tmp/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	logger = log.New(f, "logger: ", log.Lshortfile)

}

func KeyListGen(file string) *list.List {
	key := list.New()
	fp, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		key.PushBack(line)
	}
	return key
}

func main() {
	key := KeyListGen("./download_key")
	zip := ZipFile{"test", "test.zip", key}
	zipList := ZipList{
		[]ZipFile{zip},
	}

	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("./layout"))))
	http.HandleFunc("/", zipList.Index)

	for i := 0; i < len(zipList.List); i++ {
		http.HandleFunc("/"+zipList.List[i].Name, zipList.List[i].ItemIndex)
		http.HandleFunc(zipList.List[i].DownloadPath(), zipList.List[i].DownloadPage)
	}

	http.ListenAndServe(":8080", nil)
}
