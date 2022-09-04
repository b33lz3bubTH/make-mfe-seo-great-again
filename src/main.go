package main

import "fmt"
import (
    "net/http"
	"github.com/gorilla/mux"
    "os"
    "strings"
    "path"
)
import "seoCore/utils"
import "seoCore/routes"





func check(e error) {
    if e != nil {
        fmt.Println("There is an Error: ", e)
        panic(e)
    }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    w.Write([]byte(utils.GetDefaultHtml()))
}

func FileServerWithCustom404(fs http.FileSystem) http.Handler {
    // just send the index.html page for not found routes
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("r.URL.Path: ", r.URL.Path)
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			NotFound(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func main() {
    distFolder := func (relativeDoc ...string) string {
        absPath := "./../dist/"
        if len(relativeDoc) > 0 {
            return absPath + relativeDoc[0]
        }
        return absPath
    }

    fmt.Println("*** Started The Frontend SPA Server")
    html_, err := os.ReadFile(distFolder("index.html"))
    check(err)
    var htmlDom string = string(html_)
    var index int = strings.Index(htmlDom, "<body>")
    if index != -1{
        utils.HTML_SEG_1 = htmlDom[0:index+6]
        utils.HTML_SEG_2 = htmlDom[index+6:len(htmlDom)]
    }else {
        fmt.Println(`no index := `)
        return
    }

    r := mux.NewRouter()
    r.HandleFunc("/products/{id}", routes.ProductRenderer)
    r.PathPrefix("/").Handler(FileServerWithCustom404(http.Dir(distFolder())))
    http.ListenAndServe("0.0.0.0:3001", r)
}



