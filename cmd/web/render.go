package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type tempalteData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	Float           map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Info            string
	Warning         string
	Error           string
	IsAuthenticated bool
	API             string
	CSSVersion      string
}

const TEMPLATE_PATH = "templates/"
const TEMPLATE_EXTENSION = ".gohtml"
const ENV_PRODUCTION = "prod"
const ENV_DEVELOPMENT = "dev"

var functions = template.FuncMap{}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *tempalteData, r *http.Request) *tempalteData {
	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *tempalteData, partials ...string) error {
	var t *template.Template
	var err error

	templateFilePath := fmt.Sprintf("%s%s.page%s", TEMPLATE_PATH, page, TEMPLATE_EXTENSION)

	_, templateInMap := app.templateCache[templateFilePath]

	if app.config.env == ENV_PRODUCTION && templateInMap {
		t = app.templateCache[templateFilePath]
	} else {
		t, err = app.parseTemplate(partials, page, templateFilePath)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &tempalteData{}
	}

	td = app.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateFilePath string) (*template.Template, error) {
	var t *template.Template
	var err error

	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("%s%s.partial%s", TEMPLATE_PATH, x, TEMPLATE_EXTENSION)
		}
	}

	templateName := filepath.Base(templateFilePath)
	baseLayoutFilePath := fmt.Sprintf("%sbase.layout%s", TEMPLATE_PATH, TEMPLATE_EXTENSION)
	if len(partials) > 0 {
		t, err = template.New(templateName).Funcs(functions).ParseFS(templateFS, baseLayoutFilePath, strings.Join(partials, ","), templateFilePath)
	} else {
		t, err = template.New(templateName).Funcs(functions).ParseFS(templateFS, baseLayoutFilePath, templateFilePath)
	}

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateFilePath] = t

	return t, nil
}
