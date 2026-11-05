package render

import (
	"html/template"
	"time"

	"golang.org/x/text/message"
	"k8s.io/apimachinery/pkg/api/resource"
)

// PageType is a page's operating mode, as a PageEntry reports it. A KDexPage
// runs in exactly one mode, chosen by which fields its spec sets.
type PageType string

const (
	// PageTypeHTML is a page composed from an archetype (spec.contentEntries).
	PageTypeHTML PageType = "html"
	// PageTypeText is a raw text document (spec.mimeType + spec.body), such as
	// robots.txt, llms.txt or sitemap.xml.
	PageTypeText PageType = "text"
	// PageTypeRewrite is an internal alias of another KDexPage or KDexFunction
	// (spec.rewrite).
	PageTypeRewrite PageType = "rewrite"
)

// PageEntry is one page in a page tree (.PageMap), as navigation and text-page
// templates see it. Children holds the pages that name this one as their
// parent, keyed by label.
type PageEntry struct {
	BasePath string          `json:"basePath" yaml:"basePath"`
	Children *map[string]any `json:"children,omitempty" yaml:"children,omitempty"`
	Href     string          `json:"href,omitempty" yaml:"href,omitempty"`
	Icon     string          `json:"icon,omitempty" yaml:"icon,omitempty"`
	Label    string          `json:"label" yaml:"label"`
	// MimeType is a text page's spec.mimeType, one of the values that
	// field's enum allows. Empty unless Type is PageTypeText.
	MimeType string `json:"mimeType,omitempty" yaml:"mimeType,omitempty"`
	Name     string `json:"name" yaml:"name"`
	// Rewrite describes an alias's target. Nil unless Type is PageTypeRewrite.
	Rewrite *PageRewrite `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	Tags    []PageTag    `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Type is the page's operating mode: html, text or rewrite. A template
	// compares it as a string: [[ if eq $p.Type "html" ]].
	Type   PageType          `json:"type,omitempty" yaml:"type,omitempty"`
	Weight resource.Quantity `json:"weight" yaml:"weight,omitempty"`
}

// PageRewrite is a rewrite page's target, as a PageEntry reports it.
type PageRewrite struct {
	// Canonical is the alias's spec.rewrite.canonical: the alias sends
	// Link: <target URL>; rel="canonical".
	Canonical bool `json:"canonical,omitempty" yaml:"canonical,omitempty"`
	// Href is the path a request at the alias's own Href is served from,
	// after any redirect the server follows internally: the URL to list in a
	// sitemap instead of the alias. Empty when the target is not servable
	// (missing, not Ready, internal, or itself a rewrite), when the path
	// cannot be built (rewrite.path needs a wildcard value the alias's
	// Href does not supply), or when the target is a page the tree leaves
	// out for its caller.
	Href string `json:"href,omitempty" yaml:"href,omitempty"`
	// Kind is the target's kind: KDexPage or KDexFunction.
	Kind string `json:"kind" yaml:"kind"`
	// Name is the target's name.
	Name string `json:"name" yaml:"name"`
}

// PageTag is a navigation-facing projection of a page's spec.tags entry. It
// mirrors api/v1alpha1.Tag's fields deliberately rather than importing that
// type, so the render package keeps its zero-dependency boundary on the API
// types. Populated by the navigation builder; consumed by nav templates as
// [[ range $t := $value.Tags ]].
type PageTag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Renderer struct {
	BasePath        string
	BrandName       string
	Contents        map[string]string
	DefaultLanguage string
	Extra           map[string]any
	Footer          string
	FootScript      string
	Header          string
	HeadScript      string
	Host            Host
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

type Host struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// Fields available when rendering templates.
type TemplateData struct {
	BasePath        string                   `json:"basePath" yaml:"basePath"`
	BrandName       string                   `json:"brandName" yaml:"brandName"`
	Content         map[string]template.HTML `json:"content" yaml:"content"`
	DefaultLanguage string                   `json:"defaultLanguage" yaml:"defaultLanguage"`
	Extra           map[string]any
	Footer          template.HTML            `json:"footer,omitempty" yaml:"footer,omitempty"`
	FootScript      template.HTML            `json:"footScript,omitempty" yaml:"footScript,omitempty"`
	Header          template.HTML            `json:"header,omitempty" yaml:"header,omitempty"`
	HeadScript      template.HTML            `json:"headScript,omitempty" yaml:"headScript,omitempty"`
	Host            Host                     `json:"host" yaml:"host"`
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
