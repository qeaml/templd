package templd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	_ "embed"
)

//go:embed style.html
var style []byte

type views struct {
	debug     bool
	style     bool
	src       http.FileSystem
	templates map[string]http.File
	ext       string
}

func (v *views) Debug() *views {
	v.debug = !v.debug
	return v
}

func (v *views) EmbedStyle() *views {
	v.style = true
	return v
}

func (v *views) log(format string, args ...any) {
	if v.debug {
		log.Printf(format+"\n", args...)
	}
}

func (v *views) loadFrom(path string, root http.File) error {
	v.log("views: loading templates from `%s`", path)
	entries, err := root.Readdir(-1)
	if err != nil {
		return err
	}
	v.log("views: contains %d entries.", len(entries))
	for _, entry := range entries {
		v.log("views: entry: %s", entry.Name())
		entrypath := path + "/" + entry.Name()
		entryfd, err := v.src.Open(entrypath)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			v.log("       (directory)")
			err = v.loadFrom(entrypath, entryfd)
			if err != nil {
				return err
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), v.ext) {
			v.log("       (ignored)")
			continue
		}
		v.log("       (template)")
		templname := strings.ReplaceAll(strings.TrimSuffix(entrypath[1:], v.ext), "/", ":")
		v.log("views: registering template `%s`\n", templname)
		v.templates[templname] = entryfd
	}
	return nil
}

func (v *views) Load() error {
	root, err := v.src.Open("/")
	if err != nil {
		return err
	}
	return v.loadFrom("", root)
}

func (v *views) Render(out io.Writer, name string, bind interface{}, layout ...string) error {
	vars, ok := bind.(Vars)
	if !ok {
		return fmt.Errorf("vars must be template.Vars")
	}
	src, ok := v.templates[name]
	if !ok {
		return fmt.Errorf("unknown template `%s`", name)
	}
	_, err := out.Write([]byte("<!DOCTYPE html><html><body><main>"))
	if err != nil {
		return err
	}
	if v.style {
		_, err = out.Write(style)
		if err != nil {
			return err
		}
	}
	err = render(src, out, vars)
	if err != nil {
		return err
	}
	_, err = out.Write([]byte("</main></body></html>"))
	if err != nil {
		return err
	}
	return nil
}

func NewViews(fsys http.FileSystem, ext string) *views {
	return &views{
		debug:     false,
		src:       fsys,
		templates: make(map[string]http.File),
		ext:       ext,
	}
}
