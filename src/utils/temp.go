package utils

import (
	"bytes"
	"html/template"
	"log"
)

func RendTemp(src string, attr interface{}) string {
	temp, err := template.New("temp").Parse(src)
	if err != nil {
		log.Fatalf("%s New Temp Error", err)
	}

	var buf bytes.Buffer
	if err := temp.Execute(&buf, &attr); err != nil {
		log.Println("模板渲染失败", err)
	}

	return buf.String()
}
