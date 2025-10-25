package render

import (
	"html/template"
	"time"

	"golang.org/x/text/message"
	"k8s.io/apimachinery/pkg/api/resource"
	"kdex.dev/crds/api/v1alpha1"
)

type Page struct {
	BasePath        string
	Contents        map[string]string
	Footer          string
	Header          string
	Title           string
	Navigations     map[string]string
	TemplateContent string
	TemplateName    string
}

type PageEntry struct {
	Children *map[string]*PageEntry `json:"children,omitempty" yaml:"children,omitempty"`
	Icon     string                 `json:"icon,omitempty" yaml:"icon,omitempty"`
	Label    string                 `json:"label" yaml:"label"`
	Name     string                 `json:"name" yaml:"name"`
	Href     string                 `json:"href,omitempty" yaml:"href,omitempty"`
	Weight   resource.Quantity      `json:"weight,omitempty" yaml:"weight,omitempty"`
}

type Renderer struct {
	DefaultLanguage string
	FootScript      string
	HeadScript      string
	Language        string
	Languages       []string
	LastModified    time.Time
	MessagePrinter  *message.Printer
	Meta            string
	Organization    string
	PageMap         *map[string]*PageEntry
	StyleItems      []v1alpha1.StyleItem
}

// Fields available when rendering templates.
type TemplateData struct {
	Content         map[string]template.HTML `json:"content" yaml:"content"`
	DefaultLanguage string
	Footer          template.HTML            `json:"footer,omitempty" yaml:"footer,omitempty"`
	FootScript      template.HTML            `json:"footScript,omitempty" yaml:"footScript,omitempty"`
	Header          template.HTML            `json:"header,omitempty" yaml:"header,omitempty"`
	HeadScript      template.HTML            `json:"headScript,omitempty" yaml:"headScript,omitempty"`
	Language        string                   `json:"language" yaml:"language"`
	Languages       []string                 `json:"languages" yaml:"languages"`
	LastModified    time.Time                `json:"lastModified" yaml:"lastModified"`
	LeftToRight     bool                     `json:"leftToRight" yaml:"leftToRight"`
	PageBasePath    string                   `json:"pageBasePath" yaml:"pageBasePath"`
	PageMap         map[string]*PageEntry    `json:"pageMap" yaml:"pageMap"`
	Meta            template.HTML            `json:"meta,omitempty" yaml:"meta,omitempty"`
	Navigation      map[string]template.HTML `json:"navigation" yaml:"navigation"`
	Organization    string                   `json:"organization" yaml:"organization"`
	Stylesheet      template.HTML            `json:"stylesheet,omitempty" yaml:"stylesheet,omitempty"`
	Title           string                   `json:"title" yaml:"title"`
}

// Functions available when rendering templates
type TemplateFuncs interface {
	// Return the key's translated form based on the current value of Language.
	// If no value is found the key is returned. If the translation contains
	// placeholders the arguments will be positionally interpolated into the
	// result.
	l10n(key string, args ...interface{}) string
}
