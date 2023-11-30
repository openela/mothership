package main

import (
	"embed"
	"fmt"
	"github.com/openela/mothership/base"
	"google.golang.org/protobuf/types/known/timestamppb"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed templates/*.html.tpl templates/partials/*.html.tpl
var templates embed.FS

func pbTimeToNaturalDate(t *timestamppb.Timestamp, args ...bool) (ret string) {
	formatStr := "Jan 2, 2006 15:04 UTC"
	defer func() {
		if len(args) > 0 && args[0] {
			if t != nil {
				tObj := t.AsTime()
				if tObj.Format(formatStr) == ret {
					ret = ""
					return
				}
			}

			ret = fmt.Sprintf(" (%s)", ret)
		}
	}()
	if t == nil {
		ret = "--"
		return
	}
	tObj := t.AsTime()

	// If less than 10 seconds old, return "Just now"
	if time.Since(tObj) < 10*time.Second {
		ret = "Just now"
		return
	}

	// If less than a minute, show seconds
	if time.Since(tObj) < time.Minute {
		seconds := time.Since(tObj).Seconds()
		ret = fmt.Sprintf("%d seconds ago", int(seconds))
		return
	}

	// If less than an hour, show minutes
	if time.Since(tObj) < time.Hour {
		minutes := time.Since(tObj).Minutes()
		if int(minutes) == 1 {
			ret = "a minute ago"
			return
		}
		ret = fmt.Sprintf("%d minutes ago", int(minutes))
		return
	}

	// If less than a day, show hours
	if time.Since(tObj) < 24*time.Hour {
		hours := time.Since(tObj).Hours()
		if int(hours) == 1 {
			ret = "an hour ago"
			return
		}
		ret = fmt.Sprintf("%d hours ago", int(hours))
		return
	}

	// If tObj is less than 10 days old, return a relative date
	if time.Since(tObj) < 10*24*time.Hour {
		days := time.Since(tObj).Hours() / 24
		if int(days) == 1 {
			ret = "a day ago"
			return
		}
		ret = fmt.Sprintf("%d days ago", int(days))
		return
	}

	ret = tObj.Format(formatStr)
	return
}

func ctx(key string) interface{} {
	ver := version
	//goland:noinspection GoBoolExpressions
	if ver == "DEV" {
		ver = strconv.FormatInt(time.Now().Unix(), 10)
	}
	return map[string]interface{}{
		"version":      ver,
		"instanceName": instanceName,
		"links": map[string]string{
			"/entries": "Entries",
		},
		"authLinks": map[string]string{
			"/workers": "Workers",
		},
	}[key]
}

func newTmpl(tmplFs fs.FS, name string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(template.FuncMap{
			"ctx":           ctx,
			"pbNaturalTime": pbTimeToNaturalDate,
			"hasPrefix":     strings.HasPrefix,
		}).
		ParseFS(tmplFs, "templates/"+name, "templates/partials/*.html.tpl")
}

func (s *server) loadTemplates() error {
	if s.templateBundle == nil {
		s.templateBundle = make(map[string]*template.Template)
	}

	tmplFiles, err := templates.ReadDir("templates")
	if err != nil {
		return err
	}

	for _, tmplFile := range tmplFiles {
		if tmplFile.IsDir() {
			continue
		}

		tmpl, err := newTmpl(templates, tmplFile.Name())
		if err != nil {
			return err
		}

		s.templateBundle[tmplFile.Name()] = tmpl
	}

	return nil
}

func (s *server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// In dev just reload
	if version == "DEV" {
		tmpl, err := newTmpl(os.DirFS("."), name+".html.tpl")
		if err != nil {
			base.LogErrorf("html render error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			base.LogErrorf("html render error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return
	}

	err := s.templateBundle[name+".html.tpl"].Execute(w, data)
	if err != nil {
		base.LogErrorf("html render error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
