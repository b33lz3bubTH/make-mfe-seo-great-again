package utils

// Thinking For Something, not sure will work or not.
// https://neilpatel.com/blog/html-tags-for-seo/ 


var FULL_HTML string = "";

var HTML_SEG_1 string = "";
var HTML_SEG_2 string = "";


func GetDefaultHtml() string{
	return HTML_SEG_1 + HTML_SEG_2
}
func SeoHtml(html string) string{
	return HTML_SEG_1 + html + HTML_SEG_2
}


var ProductBodyHtml = `
        <article id="seo_block" style="display: none">
            <h2>Product Name: <span> %s </span></h2>
            <p>Product Description: <span> %s </span></p>
            <p>Product Price: <span> Rs. %d </span></p>
            <p>Product ID: <span> %s </span></p>
        </article>
`