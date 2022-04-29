package main

import "net/http"

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["stripe-key"] = app.config.stripe.key

	var data = &tempalteData{
		StringMap: stringMap,
	}

	if err := app.renderTemplate(w, r, "terminal", data); err != nil {
		app.errorLog.Println(err)
	}
}
