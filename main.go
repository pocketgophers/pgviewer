//go:generate packr

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/rjeczalik/notify"
	"github.com/spf13/hugo/livereload"
)

// / for packed assets
// /post for the post
// /livejs or whatever for the live view

func main() {
	log.SetFlags(log.Lshortfile)
	fmt.Println("starting...")

	addr := flag.String("addr", "127.0.0.1:3001", "Address string for http server")
	flag.Parse()

	assets := packr.NewBox("./pocketgophers")

	livereloader()

	tmpl, err := template.New("post.html").Parse(assets.String("templates/post.html"))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = tmpl.New("partials.html").Parse(assets.String("templates/partials.html"))
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/site.css", func(w http.ResponseWriter, r *http.Request) {
		// combines normalize.css and style.css
		w.Header().Add("Content-Type", "text/css")

		_, err := w.Write(assets.Bytes("normalize.css"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(assets.Bytes("style.css"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/post/" {
			http.StripPrefix("/post/", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
			return
		}

		post, err := loadPost(".", "/post/")
		if err != nil {
			httpError(w, err)
			return
		}

		buf := &bytes.Buffer{}
		err = tmpl.ExecuteTemplate(buf, "post.html", post)
		if err != nil {
			httpError(w, err)
			return
		}
		w.Write(buf.Bytes())
	})

	http.Handle("/", http.FileServer(assets))

	// addr := "127.0.0.1:3000"
	fmt.Printf("Serving on http://%s\n", *addr)
	log.Fatalln(http.ListenAndServe(*addr, nil))
}

func livereloader() {
	livereload.Initialize()
	http.HandleFunc("/livereload.js", livereload.ServeJS)
	http.HandleFunc("/livereload", livereload.Handler)

	root, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln(err)
	}

	events := make(chan notify.EventInfo, 1)
	err = notify.Watch(root+"...", events, notify.All)
	if err != nil {
		log.Fatalln(err)
	}
	// defer func() { notify.Stop(events) }()
	go func() {
		for event := range events {
			livereload.RefreshPath(event.Path())
		}
	}()
}

func httpError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, `<script>document.write('<script src="/livereload.js?port=' 
+ (location.port || 80)
+'"></'
+ 'script>')</script><h1>Error</h1><p>%v</p>`, err)
}
