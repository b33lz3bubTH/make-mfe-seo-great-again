package routes


import "seoCore/utils"
import (
    "net/http"
	"fmt"
	"os"
	"encoding/json"
    "io/ioutil"
    "log"
	"github.com/gorilla/mux"
)


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
    w.Write([]byte(utils.SeoHtml(s)))
	
}