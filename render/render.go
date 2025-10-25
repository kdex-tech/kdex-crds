package render

import (
	"bytes"
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
		Language:     r.Language,
		Languages:    r.Languages,
		LastModified: date,
		Organization: r.Organization,
		PageMap:      pageMap,
		Title:        page.Title,
	}

	contentOutputs := make(map[string]template.HTML)
	for slot, content := range page.Contents {
		output, err := r.RenderOne(fmt.Sprintf("%s-content-%s", page.TemplateName, slot), content, templateData)
		if err != nil {
			return "", err
		}

		contentOutputs[slot] = template.HTML(output)
	}
	templateData.Content = contentOutputs

	footerOutput, err := r.RenderOne(fmt.Sprintf("%s-footer", page.TemplateName), page.Footer, templateData)
	if err != nil {
		return "", err
	}
	templateData.Footer = template.HTML(footerOutput)

	footerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-footscript", page.TemplateName), r.FootScript, templateData)
	if err != nil {
		return "", err
	}
	templateData.FootScript = template.HTML(footerScriptOutput)

	headerOutput, err := r.RenderOne(fmt.Sprintf("%s-header", page.TemplateName), page.Header, templateData)
	if err != nil {
		return "", err
	}
	templateData.Header = template.HTML(headerOutput)

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

	navigationOutputs := make(map[string]template.HTML)
	for name, content := range page.Navigations {
		output, err := r.RenderOne(fmt.Sprintf("%s-navigation-%s", page.TemplateName, name), content, templateData)
		if err != nil {
			return "", err
		}
		navigationOutputs[name] = template.HTML(output)
	}
	templateData.Navigation = navigationOutputs

	stylesheetOutput, err := r.RenderOne(fmt.Sprintf("%s-stylesheet", page.TemplateName), r.Stylesheet, templateData)
	if err != nil {
		return "", err
	}
	templateData.Stylesheet = template.HTML(stylesheetOutput)

	return r.RenderOne(page.TemplateName, page.TemplateContent, templateData)
}

func (r *Renderer) RenderOne(
	templateName string,
	templateContent string,
	data TemplateData,
) (string, error) {
	funcs := template.FuncMap{
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
