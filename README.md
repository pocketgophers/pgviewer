# Pocket Gophers' Viewer

- displays a post the same way as https://pocketgophers.com
- live reload when the source changes

## Installation

```
go get pocketgophers.com/pgviewer
```

## Run

```
cd directory-with-post
pgviewer
```

`pgviewer` will start a webserver that listens to 127.0.0.1:3001. Point your web browser there and the post will be displayed. When a file in the directory changes, the post will be automatically updated and the browser view refreshed.

## What is a post

- a directory and the files contained therein
- `index.html`: the content of the post
- `info.toml`: information about the post
- other files needed for your post

## `index.html`

- the html formatted content of the post
- does not include head or body, but is placed inside a template
- can do anything that can be done in html/js/css

## `info.toml`

[toml](https://github.com/toml-lang/toml) file with two fields: `Title` and `Versions`

- `Title`: a string that represents the title of the post. Will be used in an `h1` and in the page title. Can include `code` tags, others upon request.
- `Versions`: an array of strings used to indicate the timeliness of the post. Go versions (e.g., `go1.8.3`) and dates (formatted as `2006-01-02`) are linked to the [Go Release Timeline](https://pocketgophers.com/go-release-timeline) to provide temporal context. Other versions are displayed as text.

Example:

```
Title = "Your post title discussing <code>defer</code>"
Versions = ["go1.8.3", "2017-08-11"]
```

## Generated posts with recreated directories

Some of the posts are generated. In the process of generation, the directory with the post may be removed and recreated. When this happens, `pgviewer` will not see the updated post because it is watching the deleted directory (at least this is how it works on macOS). To get around this, I run `pgviewer` in another directory and copy the post contents there after the generation process completes. For example:

```
generate-post ; cp directory-with-post/* directory-where-pgviewer-is-running
```

## Stops after "Starting..." but before "Serving…"

This is a bug with the notifier package. If you know how to fix it, please tell me.

This is not a big problem because it only happens when `pgviewer` is starting. Just cancel the process (e.g., Ctrl-c) and try again. Unless you are working on `pgviewer` and therefore frequently restarting it, you probably won't encounter this problem.