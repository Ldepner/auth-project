package render

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

const TemplatesPath = "./templates"

func Template(w http.ResponseWriter, fileName string) error {
	name, _ := filepath.Glob(fmt.Sprintf("%s/%s", TemplatesPath, fileName))
	tmpl, err := template.ParseFiles(name[0])
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, nil); err != nil {
		return err
	}

	return nil
}
