package webpages

import (
	"html/template"
	"net/http"
)

const (
	EmailConfirmPage1 = "<!DOCTYPE html><html><body><h1>My First Heading</h1><p>My first paragraph.</p></body></html>"
)

func SendWebPage(w http.ResponseWriter, dir string, data map[string]string) error {
	// abs, err := filepath.Abs(dir)
	// if err != nil {
	// 	log.Print(abs, err)
	// 	return err
	// }
	tmp := template.Must(template.ParseFiles(dir))

	return tmp.Execute(w, data)
}
