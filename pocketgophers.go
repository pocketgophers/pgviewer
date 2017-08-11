package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/naoina/toml"
	"github.com/pkg/errors"
)

type Post struct {
	Info    PostInfo
	Content template.HTML
	Path    string
}

func (p Post) Version() template.HTML {
	if len(p.Info.Versions) == 0 || len(p.Info.Versions[0]) == 0 {
		// template renderer errors when rendering ""
		// so, avoid ""
		return " "
	}

	buf := &bytes.Buffer{}
	for i, version := range p.Info.Versions {
		if i != 0 {
			buf.WriteString(", ")
		}

		// go release?
		if strings.HasPrefix(version, "go") {
			fmt.Fprintf(buf, "<a href=\"https://pocketgophers.com/go-release-timeline/#%s\">%s</a>", version, version)
			continue
		}

		// date?
		_, err := time.Parse("2006-01-02", version)
		if err == nil {
			fmt.Fprintf(buf, "<a href=\"https://pocketgophers.com/go-release-timeline/#%s\">%s</a>", version, version)
			continue
		}

		buf.WriteString(template.HTMLEscapeString(version))
	}

	return template.HTML(buf.String())
}

type PostInfo struct {
	TitleStr string `toml:"Title"`
	Versions []string
}

func (pi *PostInfo) Title() template.HTML {
	return template.HTML(pi.TitleStr)
}

func (pi *PostInfo) TitleText() string {
	noCode := strings.Replace(pi.TitleStr, "<code>", "", -1)
	return strings.Replace(noCode, "</code>", "", -1)
}

func loadPost(dirname string, path string) (*Post, error) {
	info, err := loadPostInfo(filepath.Join(dirname, "info.toml"))
	if err != nil {
		return nil, err
	}

	post := &Post{
		Info: *info,
		Path: path,
	}

	if post.Info.TitleStr == "" {
		post.Info.TitleStr = "Set your post title in <code>info.toml</code>, e.g., <code>Title = \"Your Post Title\"</code>"
	}

	b, err := ioutil.ReadFile("index.html")
	if err != nil {
		return nil, errors.Wrap(err, "index.html")
	}

	post.Content = template.HTML(b)

	return post, nil
}

func loadPostInfo(filename string) (*PostInfo, error) {
	info := &PostInfo{}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, filename)
	}

	err = toml.Unmarshal(b, &info)
	if err != nil {
		return nil, errors.Wrap(err, filename)
	}

	return info, nil
}
