package render

import (
	"html/template"
	"time"

	"golang.org/x/text/message"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PageEntry struct {
	BasePath string            `json:"basePath" yaml:"basePath"`
	Children *map[string]any   `json:"children,omitempty" yaml:"children,omitempty"`
	Href     string            `json:"href,omitempty" yaml:"href,omitempty"`
	Icon     string            `json:"icon,omitempty" yaml:"icon,omitempty"`
	Label    string            `json:"label" yaml:"label"`
	Name     string            `json:"name" yaml:"name"`
	Weight   resource.Quantity `json:"weight" yaml:"weight,omitempty"`
}

type Renderer struct {
	BasePath        string
	BrandName       string
	Contents        map[string]string
	DefaultLanguage string
	Footer          string
	FootScript      string
	Header          string
	HeadScript      string
	Language        string
	Languages       []string
	LastModified    time.Time
	MessagePrinter  *message.Printer
	Navigations     map[string]string
	Meta            string
	Organization    string
	PageMap         map[string]any
	PatternPath     string
	TemplateContent string
	TemplateName    string
	Title           string
	Theme           string
}

// Fields available when rendering templates.
type TemplateData struct {
	BasePath        string                   `json:"basePath" yaml:"basePath"`
	BrandName       string                   `json:"brandName" yaml:"brandName"`
	Content         map[string]template.HTML `json:"content" yaml:"content"`
	DefaultLanguage string                   `json:"defaultLanguage" yaml:"defaultLanguage"`
	Footer          template.HTML            `json:"footer,omitempty" yaml:"footer,omitempty"`
	FootScript      template.HTML            `json:"footScript,omitempty" yaml:"footScript,omitempty"`
	Header          template.HTML            `json:"header,omitempty" yaml:"header,omitempty"`
	HeadScript      template.HTML            `json:"headScript,omitempty" yaml:"headScript,omitempty"`
	Language        string                   `json:"language" yaml:"language"`
	Languages       []string                 `json:"languages" yaml:"languages"`
	LastModified    time.Time                `json:"lastModified" yaml:"lastModified"`
	LeftToRight     bool                     `json:"leftToRight" yaml:"leftToRight"`
	Navigation      map[string]template.HTML `json:"navigation" yaml:"navigation"`
	Meta            template.HTML            `json:"meta,omitempty" yaml:"meta,omitempty"`
	Organization    string                   `json:"organization" yaml:"organization"`
	PageMap         map[string]any           `json:"pageMap" yaml:"pageMap"`
	PatternPath     string                   `json:"patternPath" yaml:"patternPath"`
	Theme           template.HTML            `json:"theme,omitempty" yaml:"theme,omitempty"`
	Title           string                   `json:"title" yaml:"title"`
}

// Functions available when rendering templates
type TemplateFuncs interface {
	// Return the key's translated form based on the current value of Language.
	// If no value is found the key is returned. If the translation contains
	// placeholders the arguments will be positionally interpolated into the
	// result.
	l10n(key string, args ...any) string
}
