package utils


var HTML_SEG_1 string = "";
var HTML_SEG_2 string = "";


func GetDefaultHtml() string{
	return HTML_SEG_1 + HTML_SEG_2
}
func SeoHtml(html string) string{
	return HTML_SEG_1 + html + HTML_SEG_2
}