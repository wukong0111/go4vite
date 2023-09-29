package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed dist/assets/*
var assets embed.FS

func main() {
	manifestData := GetViteMainScriptName()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal("can not parse template")
		}

		tmpl.Execute(w, manifestData)
	})

	subFS, err := fs.Sub(assets, "dist/assets")
	if err != nil {
		log.Fatal("can not access to assets subdirectory")
	}
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.FS(subFS))))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("can not run server:", err)
	}
}

func GetViteMainScriptName() map[string]map[string]any {
	manifestBytes, err := os.ReadFile("dist/manifest.json")

	if err != nil {
		log.Fatal("can not open vite manifest file", err)
	}
	var manifestMap map[string]map[string]any
	if err := json.Unmarshal(manifestBytes, &manifestMap); err != nil {
		log.Fatal("can not decode manifest file", err)
	}
	return manifestMap
}

