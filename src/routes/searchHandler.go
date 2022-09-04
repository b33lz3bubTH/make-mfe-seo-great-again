package routes


import "seoCore/utils"
import (
    "net/http"
	"fmt"
	"encoding/json"
    "io/ioutil"
    "log"
	// "github.com/gorilla/mux"
)

type ProductListApiResponse struct {
	APIResponseInfo struct {
		Message string `json:"message"`
	} `json:"api_response_info"`
	Data struct {
		Products []struct {
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
		} `json:"products"`
		PageSize     int `json:"pageSize"`
		PageNumber   int `json:"pageNumber"`
		TotalResults int `json:"totalResults"`
		TotalPage    int `json:"totalPage"`
	} `json:"data"`
}


func SearchProductHandler(w http.ResponseWriter, r *http.Request) {
	// dummy URL: http://localhost:3001/search?q=men

	queryParamQ := r.URL.Query().Get("q")
	queryParamTags := r.URL.Query().Get("tags")
	queryTerm := ""


	var apiUrl string = "https://api.nrix.in/product-api/ui-resolver/domain/nrix/products/search?pageSize=100&" 


	if queryParamQ == "" {
		apiUrl += "tags=" + queryParamTags
		queryTerm = queryParamTags
	} else if queryParamTags == "" {
		apiUrl += "q=" + queryParamQ
		queryTerm = queryParamQ
	}

	response, err := http.Get(apiUrl)

	if err != nil {
        fmt.Print(err.Error())
        // os.Exit(1)
    }

	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("err: ", err)
        log.Fatal(err)
    }



    var productListApiResponse ProductListApiResponse
    json.Unmarshal(responseData, &productListApiResponse)



	rawHtml := ""

	productList := &productListApiResponse.Data.Products

	for i := 0; i < len(*productList); i++ {
		rawHtml += fmt.Sprintf(string(utils.ProductBodyHtml), (*productList)[i].Name, (*productList)[i].Description,(*productList)[i].SellingPrice, (*productList)[i].BaseCode)
	}

	productHtml := []byte(`
		<section id="seo_block" style="display: none">
		<h1>Product List For %s</h1>
			%s
		</section>
	`)


	s := fmt.Sprintf(string(productHtml), queryTerm, rawHtml)


	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    w.Write([]byte(utils.SeoHtml(s)))

}