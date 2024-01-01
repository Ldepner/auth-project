package render

import (
	"fmt"
	"github.com/Ldepner/auth-project/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

const TemplatesPath = "./templates"

func Template(w http.ResponseWriter, fileName string, td *models.TemplateData) error {
	name, _ := filepath.Glob(fmt.Sprintf("%s/%s", TemplatesPath, fileName))
	tmpl, err := template.ParseFiles(name[0])
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, td); err != nil {
		return err
	}

	return nil
}
