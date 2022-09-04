package main

import "fmt"
import (
    "net/http"
	"github.com/gorilla/mux"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "path"
)

var HTML_SEG_1 string = "";
var HTML_SEG_2 string = "";


func check(e error) {
    if e != nil {
        fmt.Println("There is an Error: ", e)
        panic(e)
    }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    w.Write([]byte(GetDefaultHtml()))
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
        HTML_SEG_1 = htmlDom[0:index+6]
        HTML_SEG_2 = htmlDom[index+6:len(htmlDom)]
    }else {
        fmt.Println(`no index := `)
        return
    }

    r := mux.NewRouter()
    r.HandleFunc("/products/{id}", ProductRenderer)
    r.PathPrefix("/").Handler(FileServerWithCustom404(http.Dir(distFolder())))
    http.ListenAndServe("0.0.0.0:3001", r)
}

func GetDefaultHtml() string{
	return HTML_SEG_1 + HTML_SEG_2
}
func SeoHtml(html string) string{
	return HTML_SEG_1 + html + HTML_SEG_2
}

type ProductData struct {
	APIResponseInfo struct {
		Message string `json:"message"`
	} `json:"api_response_info"`
	Data struct {
		UpdatedOn   string   `json:"updatedOn"`
		Name        string   `json:"name"`
		Domain      string   `json:"domain"`
		BaseCode    string   `json:"baseCode"`
		Tags        []string `json:"tags"`
		FeatureSets []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"featureSets"`
		Medias       []string `json:"medias"`
		SellingPrice int      `json:"sellingPrice"`
		IsActive     bool     `json:"isActive"`
		Description  string   `json:"description"`
		UUID         string   `json:"uuid"`
		MrpPrice     int      `json:"mrpPrice"`
		VarientList  []struct {
			VarientTypeName string `json:"varientTypeName"`
			UUID            string `json:"uuid"`
			Label           string `json:"label"`
			UpdatedOn       string `json:"updatedOn"`
			VarientOptions  []struct {
				PriceDiffrence        int           `json:"priceDiffrence"`
				IsActive              bool          `json:"isActive"`
				ParentVarientUUIDInfo string        `json:"parentVarientUUIDInfo"`
				Value                 string        `json:"value"`
				UUID                  string        `json:"uuid"`
				VarientTypeName       string        `json:"varientTypeName"`
				Label                 string        `json:"label"`
				Medias                []interface{} `json:"medias"`
				UpdatedOn             string        `json:"updatedOn"`
			} `json:"varientOptions"`
		} `json:"varientList"`
	} `json:"data"`
}



func ProductRenderer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        fmt.Println("id is missing in parameters")
    }
    fmt.Println(`id := `, id)
    response, err := http.Get("https://api.nrix.in/product-api/ui-resolver/domain/nrix/products/" + id)

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("err: ", err)
        log.Fatal(err)
    }

    var product ProductData
    json.Unmarshal(responseData, &product)

    fmt.Println("product name: ", product.Data.Name)

    rawHtml := []byte(`
        <div id="_SEO_SHIT" style="display: none">
            <p>Product Name: <span> %s </span></p>
            <p>Product Description: <span> %s </span></p>
            <p>Product Price: <span> %d </span></p>
            <p>Product ID: <span>%s</span></p>
        </div>
    `)
    s := fmt.Sprintf(string(rawHtml), product.Data.Name, product.Data.Description,product.Data.SellingPrice, id)

    
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    w.Write([]byte(SeoHtml(s)))
	
}