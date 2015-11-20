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

// ZipFile is download item struct
// ダウンロードするアイテムを格納する構造体です
type ZipFile struct {
	Name string
	Path string
	Key  *list.List
}

// ZipList is download items contain
// ダウンロードするアイテムのリストを格納します
type ZipList struct {
	List []ZipFile
}

// Index is ZipList http.Handler
// トップページを生成します
func (z *ZipList) Index(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("index.html")

	if err != nil {
		logger.Println(err)
		return
	}

	err = tmple.Execute(w, z)

	if err != nil {
		logger.Println(err)
		return
	}
}

// ItemIndex is ZipFile http.Handler
// それぞれのアイテムをダウンロードするページを生成します
func (z *ZipFile) ItemIndex(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("layout.html")

	if err != nil {
		logger.Println(err)
		return
	}

	err = tmple.Execute(w, z)

	if err != nil {
		logger.Println(err)
		return
	}
}

// DownloadPage is ZipFile http.Handler
// アイテムをダウンロードを行うページです。認証を行います
func (z *ZipFile) DownloadPage(w http.ResponseWriter, r *http.Request) {
	num := r.FormValue("num")

	if z.check(num) {
		file, err := ioutil.ReadFile(z.Path)
		if err != nil {
			logger.Print(err)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(file)
		if err != nil {
			logger.Print(err)
			return
		}
	} else {
		http.Redirect(w, r, "/"+z.Name+"?error=1", http.StatusFound)
	}

}

// DownloadPath is create item download path
// アイテムをダウンロードするパスを生成します。
func (z *ZipFile) DownloadPath() string {
	buf := z.Name
	buf = "/" + buf + ".zip"
	return buf
}

// check return certified success or fail
// ダウンロードの鍵の確認を行う関数
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

// KeyListGen is create key list from keyfile
// 鍵の一覧のファイルから鍵のリストを生成します
func KeyListGen(file string) *list.List {
	key := list.New()
	fp, err := os.Open(file)
	if err != nil {
		fmt.Print(err)
	}

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Print(err)
		}
		key.PushBack(line)
	}
	return key
}

func main() {
	f, err := os.OpenFile("./log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Println(err)
	}

	logger = log.New(f, "logger: ", log.Lshortfile)

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

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Print(err)
	}
}
