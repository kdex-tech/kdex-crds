package render

import (
	htmltemplate "html/template"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"k8s.io/apimachinery/pkg/api/resource"
)

var renderer Renderer
var templateData TemplateData

func init() {
	translations := catalog.NewBuilder()
	translations.SetString(language.English, "name", "Name")
	translations.SetString(language.French, "name", "Nom")

	messagePrinter := message.NewPrinter(
		language.English,
		message.Catalog(translations),
	)

	templateData = TemplateData{
		DefaultLanguage: "en",
		Footer:          `<p>footer</p>`,
		FootScript:      `<script type="text/javascript"></script>`,
		Header:          `<p>header</p>`,
		HeadScript:      `<script type="text/javascript"></script>`,
		Language:        "en",
		Languages:       []string{"en", "fr"},
		LastModified:    time.Now(),
		LeftToRight:     true,
		Meta:            `<meta charset="UTF-8">`,
		Organization:    "KDex Tech Inc.",
		PageBasePath:    "/one",
		PageMap: map[string]*PageEntry{
			"One": {
				Href:   "/one",
				Icon:   "one",
				Label:  "One",
				Name:   "one",
				Weight: resource.MustParse("0"),
			},
			"Two": {
				Href:   "/two",
				Icon:   "two",
				Label:  "Two",
				Name:   "two",
				Weight: resource.MustParse("1"),
			},
			"Three": {
				Children: &map[string]*PageEntry{
					"Four": {
						Href:   "/four",
						Icon:   "four",
						Label:  "Four",
						Name:   "four",
						Weight: resource.MustParse("0"),
					},
				},
				Href:   "/three",
				Icon:   "three",
				Label:  "Three",
				Name:   "three",
				Weight: resource.MustParse("3"),
			},
		},
		Theme: `<style>body {color: red;}</style>`,
		Title: "name",
	}

	contents := map[string]htmltemplate.HTML{}
	contents["main"] = htmltemplate.HTML("<p>content</p>")
	templateData.Content = contents

	navigations := map[string]htmltemplate.HTML{}
	navigations["main"] = htmltemplate.HTML("<p>navigation</p>")
	templateData.Navigation = navigations

	renderer = Renderer{
		MessagePrinter: messagePrinter,
	}
}

func DefaultTemplateData() TemplateData {
	return templateData
}

func ValidateContent(
	name string, content string,
) error {
	_, err := renderer.RenderOne(name, content, DefaultTemplateData())

	return err
}
