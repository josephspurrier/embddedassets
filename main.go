package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/josephspurrier/embedded-assets/static"
)

func main() {
	// Render files from disk
	renderFiles("home")
	fmt.Print("\n\n")
	renderFiles("about")
	fmt.Print("\n\n")

	// Render embedded files
	renderEmbeddedAssets("home")
	fmt.Print("\n\n")
	renderEmbeddedAssets("about")
	fmt.Print("\n\n")

	// Access embedded file.
	f, _ := fileContents("about") // Get the embedded asset
	b, _ := io.ReadAll(f)         // Get the contents in bytes
	fmt.Println(string(b))        // Will output: {{define "title"}}About{{end}}...
	fi, _ := f.Stat()             // Get the asset information
	fmt.Println(fi.ModTime())     // Will output: 0001-01-01 00:00:00 +0000 UTC
	fmt.Printf(`%x`, md5.Sum(b))  // Output the MD5 checksum
}

// renderFiles takes the name of a template on disk and renders it to the screen.
func renderFiles(tmpl string) {
	t, err := template.ParseFS(static.Assets, "base.tmpl", fmt.Sprintf("pages/%v.tmpl", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(os.Stdout, nil); err != nil {
		log.Fatal(err)
	}
}

// renderEmbeddedAssets takes the name of an embedded template and renders it to the screen.
func renderEmbeddedAssets(tmpl string) {
	t, err := template.ParseFS(static.Assets, "base.tmpl", fmt.Sprintf("pages/%v.tmpl", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(os.Stdout, nil); err != nil {
		log.Fatal(err)
	}
}

// fileContents returns an embedded template and an error if one occurs.
func fileContents(tmpl string) (fs.File, error) {
	return static.Assets.Open(fmt.Sprintf("pages/%v.tmpl", tmpl))
}
