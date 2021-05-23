# How to Embed Assets in Go 1.16

The original article is [here](https://www.josephspurrier.com/how-to-embed-assets-in-go-1-16).

In Go 1.16, there is now a way to [natively embed files and templates](https://golang.org/pkg/embed/) which allows us to now perform single binary deployments easily. A single binary deployment is beneficial because you can easily distribute and replace a single file that contains both the compiled Go code as well as any static assets like HTML templates and CSS files. It's less effort to replace a single file than it is to replace an entire folder of assets. Prior, most people used third-party offerings like the [jteeuwen/gobindata](https://github.com/jteeuwen/go-bindata) library that provided similar functionality. 

It doesn't require a lot of changes if you are already using static assets in your application so I'll walk how to migrate it over and what to watch out for. You can read the original proposal for embedding static assets [here](https://github.com/golang/go/issues/35950).

## Initial Setup
You'll need to be using Go 1.16 or newer. You'll want to read the [Go docs on how the embed package](https://golang.org/pkg/embed/) works. I'll be referencing the article, [How to use Template Blocks in Go 1.6](https://www.josephspurrier.com/how-to-use-template-blocks-in-go-1-6), which shows how to parse Go templates.

The file structure we'll be referencing throughout the article is here:

```
/static/pages/about.tmpl - About template
/static/pages/home.tmpl  - Home template
/static/base.tmpl        - Base template
/static/static.go        - Go file with //go:embed directive
/main.go                 - Go main and render functions
```

The repository for this article is available on GitHub: [josephspurrier/embedded-assets](https://github.com/josephspurrier/embedded-assets).

## //go:embed Directive
The `//go:embed` directive (which is a comment in Go) allows you to specify which files you want Go to embed in your application when you compile. You can use a single line directive or multi-line directives, whichever you prefer. You will need a Go file at the root of the directory where your static files will be on disk and this will determine the asset paths.

For this example, we'll create a Go file (/static/static.go) that shows which files to include. You can use wildcards for both the file name and the folder names if you want. It's best to use an extension at the end of the wildcard so you don't accidentally include extra files. If the directive contains just a folder, then containing files starting with '.' and '_' in the name are excluded.

```go
package static

import  "embed"

// Assets represents the embedded files.
//go:embed *.tmpl pages/*.tmpl
var Assets embed.FS
```

**Referencing:** If you have a template that exists at /static/pages/about.tmpl and your embed directive is in /static/static.go, then you will access the about.tmpl file using this path from the /main.go file: `pages/about.tmpl`.

## Assets as Templates
The assets are available as a file system so when you reference them in your code, you can use them like you're used to. One use case for embedded assets is with the Go template library.

If you want to access the files from disk, you may have a render function like this in your main.go file. Notice the paths start with the folder name `static`. The template locations are relative to the main.go file.

```go
func renderFiles(tmpl string) {
    t, err := template.ParseFiles("static/base.tmpl", fmt.Sprintf("static/pages/%v.tmpl", tmpl))
    if err != nil {
        log.Fatal(err)
    }

    if err := t.Execute(os.Stdout, nil); err != nil {
        log.Fatal(err)
    }
}
```

To use the embedded assets, you would the add the embed directive to a Go file inside of the static folder, update `template.ParseFiles` to `template.ParseFS`, pass the `static.Assets` variable as the first argument, and then remove the `static/` folder from the paths.

```go
func renderEmbeddedAssets(tmpl string) {
    t, err := template.ParseFS(static.Assets, "base.tmpl", fmt.Sprintf("pages/%v.tmpl", tmpl))
    if err != nil {
        log.Fatal(err)
    }

    if err := t.Execute(os.Stdout, nil); err != nil {
        log.Fatal(err)
    }
}
```

## Assets as Files
You can also access the embedded assets using the typical `io` commands because it implements the fs.File interface. This is an example of how to read the file contents and access the metadata.

```go
func fileContents(tmpl string) (fs.File, error) {
    return static.Assets.Open(fmt.Sprintf("pages/%v.tmpl", tmpl))
}

func main() {
    f, _ := fileContents("about") // Get the embedded asset
    b, _ := io.ReadAll(f)         // Get the contents in bytes
    fmt.Println(string(b))        // Will output: {{define "title"}}About{{end}}...
    fi, _ := f.Stat()             // Get the asset information
    fmt.Println(fi.ModTime())     // Will output: 0001-01-01 00:00:00 +0000 UTC
    fmt.Printf(`%x`, md5.Sum(b))  // Output the MD5 checksum
}
```

One caveat is the `ModTime()` of an embedded asset will always return: `0001-01-01 00:00:00 +0000 UTC`. It was decided that all embedded files would have a modification time of zero to help with [reproducibility](https://github.com/golang/go/issues/35950#issuecomment-666997865) just like modules. If you do need to check for changes, you could always generate a MD5 checksum to compare.