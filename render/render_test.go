package render

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRenderOne(t *testing.T) {
	data := TemplateData{
		Title: "World",
	}
	templateContent := "Hello, {{.Title}}!"

	r := &Renderer{}
	actual, err := r.RenderOne("test", templateContent, data)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", actual)
}

func TestRenderOne_InvalidTemplate(t *testing.T) {
	data := TemplateData{
		Title: "World",
	}
	templateContent := "Hello, {{.Invalid}}!"

	r := &Renderer{}
	_, err := r.RenderOne("test", templateContent, data)
	assert.Error(t, err)
}

func TestRenderAll(t *testing.T) {
	lastModified, _ := time.Parse("2006-01-02", "2025-09-20")

	page := Page{
		Contents: map[string]string{
			"main":    "<h1>Welcome</h1>",
			"sidebar": `<my-app-element id="sidebar" data-date="{{.LastModified.Format "2006-01-02"}}"></my-app-element>`,
		},
		Footer: "Page Footer",
		Header: "Page Header",
		Navigations: map[string]string{
			"main": "main-nav",
		},
		Title:        "Test Page",
		TemplateName: "main",
		TemplateContent: `<!DOCTYPE html>
<html lang="{{ .Language }}">
	<head>
	<title>{{.Title}}</title>
		{{.Meta}}
		{{.HeadScript}}
		{{.Stylesheet}}
	</head>
	<body>
		<header>{{.Header}}</header>
		<nav>{{range $key, $value := .Navigation}}
		{{$key}}: {{$value}}
		{{end}}</nav>
		<main>{{range $key, $value := .Content}}
		<div id="slot-{{$key}}">{{$value}}</div>
		{{end}}</main>
		<footer>{{.Footer}}</footer>
		{{.FootScript}}
	</body>
</html>`,
	}

	r := &Renderer{
		FootScript:   "<script>foot</script>",
		HeadScript:   "<script>head</script>",
		Language:     "en",
		LastModified: lastModified,
		PageMap:      &map[string]*PageEntry{"home": {Href: "/"}},
		Meta:         `<meta name="description" content="test">`,
		Organization: "Test Inc.",
	}
	actual, err := r.RenderPage(page)
	assert.NoError(t, err)

	assert.Contains(t, actual, "<title>Test Page</title>")
	assert.Contains(t, actual, r.Meta)
	assert.Contains(t, actual, r.HeadScript)
	assert.Contains(t, actual, "Page Header")
	assert.Contains(t, actual, "main: main-nav")
	assert.Contains(t, actual, "<h1>Welcome</h1>")
	assert.Contains(t, actual, "<my-app-element id=\"sidebar\"")
	assert.Contains(t, actual, "2025-09-20")
	assert.Contains(t, actual, "Page Footer")
	assert.Contains(t, actual, r.FootScript)
}

func TestRenderAll_InvalidHeaderTemplate(t *testing.T) {
	r := &Renderer{}
	page := Page{
		TemplateName: "main",
		Navigations:  nil,
		Header:       "{{.Invalid}}",
		Footer:       "",
	}
	_, err := r.RenderPage(page)
	assert.Error(t, err)
}

func TestRenderAll_InvalidFooterTemplate(t *testing.T) {
	r := &Renderer{}
	page := Page{
		TemplateName: "main",
		Navigations:  nil,
		Header:       "",
		Footer:       "{{.Invalid}}",
	}
	_, err := r.RenderPage(page)
	assert.Error(t, err)
}

func TestRenderAll_InvalidNavigationTemplate(t *testing.T) {
	r := &Renderer{}
	page := Page{
		TemplateName: "main",
		Navigations: map[string]string{
			"main": "{{.Invalid}}",
		},
		Header: "",
		Footer: "",
	}
	_, err := r.RenderPage(page)
	assert.Error(t, err)
}

func TestRenderAll_InvalidContentTemplate(t *testing.T) {
	r := &Renderer{}
	page := Page{
		TemplateName: "main",
		Contents: map[string]string{
			"main": "{{.Invalid}}",
		},
		Navigations: nil,
		Header:      "",
		Footer:      "",
	}
	_, err := r.RenderPage(page)
	assert.Error(t, err)
}

func TestRenderAll_InvalidMainTemplate(t *testing.T) {
	r := &Renderer{}
	page := Page{
		TemplateName:    "main",
		TemplateContent: "{{.Invalid}}",
		Navigations:     nil,
		Header:          "",
		Footer:          "",
	}
	_, err := r.RenderPage(page)
	assert.Error(t, err)
}
