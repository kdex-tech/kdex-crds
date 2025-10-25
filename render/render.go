package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"
)

func (r *Renderer) RenderPage(page Page) (string, error) {
	date := r.LastModified
	if date.IsZero() {
		date = time.Now()
	}

	pageMap := map[string]*PageEntry{}
	if r.PageMap != nil {
		pageMap = *r.PageMap
	}

	templateData := TemplateData{
		DefaultLanguage: r.DefaultLanguage,
		Language:        r.Language,
		Languages:       r.Languages,
		LastModified:    date,
		Organization:    r.Organization,
		PageBasePath:    page.BasePath,
		PageMap:         pageMap,
		Title:           page.Title,
	}

	//
	// These don't need access to content or navigation
	//

	footerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-footscript", page.TemplateName), r.FootScript, templateData)
	if err != nil {
		return "", err
	}
	templateData.FootScript = template.HTML(footerScriptOutput)

	headerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-headscript", page.TemplateName), r.HeadScript, templateData)
	if err != nil {
		return "", err
	}
	templateData.HeadScript = template.HTML(headerScriptOutput)

	metaOutput, err := r.RenderOne(fmt.Sprintf("%s-meta", page.TemplateName), r.Meta, templateData)
	if err != nil {
		return "", err
	}
	templateData.Meta = template.HTML(metaOutput)

	stylesheetOutput, err := r.RenderOne(fmt.Sprintf("%s-stylesheet", page.TemplateName), r.StyleItemsToString(), templateData)
	if err != nil {
		return "", err
	}
	templateData.Stylesheet = template.HTML(stylesheetOutput)

	//
	// Content and Navigation
	//

	contentOutputs := make(map[string]template.HTML)
	for slot, content := range page.Contents {
		output, err := r.RenderOne(fmt.Sprintf("%s-content-%s", page.TemplateName, slot), content, templateData)
		if err != nil {
			return "", err
		}

		contentOutputs[slot] = template.HTML(output)
	}
	templateData.Content = contentOutputs

	navigationOutputs := make(map[string]template.HTML)
	for name, content := range page.Navigations {
		output, err := r.RenderOne(fmt.Sprintf("%s-navigation-%s", page.TemplateName, name), content, templateData)
		if err != nil {
			return "", err
		}
		navigationOutputs[name] = template.HTML(output)
	}
	templateData.Navigation = navigationOutputs

	//
	// Footer, Header may wish to access both content and navigation
	//

	footerOutput, err := r.RenderOne(fmt.Sprintf("%s-footer", page.TemplateName), page.Footer, templateData)
	if err != nil {
		return "", err
	}
	templateData.Footer = template.HTML(footerOutput)

	headerOutput, err := r.RenderOne(fmt.Sprintf("%s-header", page.TemplateName), page.Header, templateData)
	if err != nil {
		return "", err
	}
	templateData.Header = template.HTML(headerOutput)

	return r.RenderOne(page.TemplateName, page.TemplateContent, templateData)
}

func (r *Renderer) RenderOne(
	templateName string,
	templateContent string,
	data TemplateData,
) (string, error) {
	funcs := template.FuncMap{
		"json": func(payload interface{}, args ...string) string {
			b, err := json.Marshal(payload)
			if err != nil {
				return err.Error()
			}
			return string(b)
		},
		"l10n": func(key string, args ...string) string {
			if r.MessagePrinter == nil {
				return key
			}
			return r.MessagePrinter.Sprintf(key, args)
		},
	}

	instance, err := template.New(templateName).Funcs(funcs).Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := instance.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (h *Renderer) StyleItemsToString() string {
	var styleBuffer bytes.Buffer

	for _, item := range h.StyleItems {
		if item.LinkHref != "" {
			styleBuffer.WriteString(`<link`)
			for key, value := range item.Attributes {
				if key == "href" || key == "src" {
					continue
				}
				styleBuffer.WriteRune(' ')
				styleBuffer.WriteString(key)
				styleBuffer.WriteString(`="`)
				styleBuffer.WriteString(value)
				styleBuffer.WriteRune('"')
			}
			styleBuffer.WriteString(` href="`)
			styleBuffer.WriteString(item.LinkHref)
			styleBuffer.WriteString(`"/>`)
			styleBuffer.WriteRune('\n')
		} else if item.Style != "" {
			styleBuffer.WriteString(`<style`)
			for key, value := range item.Attributes {
				if key == "href" || key == "src" {
					continue
				}
				styleBuffer.WriteRune(' ')
				styleBuffer.WriteString(key)
				styleBuffer.WriteString(`="`)
				styleBuffer.WriteString(value)
				styleBuffer.WriteRune('"')
			}
			styleBuffer.WriteRune('>')
			styleBuffer.WriteRune('\n')
			styleBuffer.WriteString(item.Style)
			styleBuffer.WriteString("</style>")
			styleBuffer.WriteRune('\n')
		}
	}

	return styleBuffer.String()
}
