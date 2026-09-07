package main

import (
	"bytes"
	"log"
	"os"
	"text/template"
)

func main() {
	const source = `package routes

import "github.com/gin-gonic/gin"
import (
	{{- range .}}
	"radio_stream/routes/{{.}}"
	{{- end}}
)

func Register_all_route(router *gin.Engine) {
	{{- range .}}
	{{.}}.Register(router)
	{{- end}}
}
`

	tmpl := template.Must(template.New("import").Parse(source))

	folder, err := os.ReadDir("routes")

	if err != nil {
		log.Fatal(err)
	}

	route := []string{}
	for _, file := range folder {
		if !file.IsDir() {
			continue
		}

		route = append(route, file.Name())
	}

	var finish bytes.Buffer
	tmpl.Execute(&finish, route)

	err = os.WriteFile("./routes/routers.go", []byte(finish.String()), 0664)
	if err != nil {
		log.Fatal(err)
	}
}
