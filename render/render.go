package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"
)

func (r *Renderer) RenderPage() (string, error) {
	date := r.LastModified
	if date.IsZero() {
		date = time.Now()
	}

	pageMap := map[string]*PageEntry{}
	if r.PageMap != nil {
		pageMap = *r.PageMap
	}

	templateData := TemplateData{
		BasePath:        r.BasePath,
		BrandName:       r.BrandName,
		DefaultLanguage: r.DefaultLanguage,
		Language:        r.Language,
		Languages:       r.Languages,
		LastModified:    date,
		Organization:    r.Organization,
		PageMap:         pageMap,
		PatternPath:     r.PatternPath,
		Title:           r.Title,
	}

	//
	// These don't need access to content or navigation
	//

	footerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-footscript", r.TemplateName), r.FootScript, templateData)
	if err != nil {
		return "", err
	}
	templateData.FootScript = template.HTML(footerScriptOutput)

	headerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-headscript", r.TemplateName), r.HeadScript, templateData)
	if err != nil {
		return "", err
	}
	templateData.HeadScript = template.HTML(headerScriptOutput)

	metaOutput, err := r.RenderOne(fmt.Sprintf("%s-meta", r.TemplateName), r.Meta, templateData)
	if err != nil {
		return "", err
	}
	templateData.Meta = template.HTML(metaOutput)

	themeOutput, err := r.RenderOne(fmt.Sprintf("%s-theme", r.TemplateName), r.Theme, templateData)
	if err != nil {
		return "", err
	}
	templateData.Theme = template.HTML(themeOutput)

	//
	// Content and Navigation
	//

	contentOutputs := make(map[string]template.HTML)
	for slot, content := range r.Contents {
		output, err := r.RenderOne(fmt.Sprintf("%s-content-%s", r.TemplateName, slot), content, templateData)
		if err != nil {
			return "", err
		}

		contentOutputs[slot] = template.HTML(output)
	}
	templateData.Content = contentOutputs

	navigationOutputs := make(map[string]template.HTML)
	for name, content := range r.Navigations {
		output, err := r.RenderOne(fmt.Sprintf("%s-navigation-%s", r.TemplateName, name), content, templateData)
		if err != nil {
			return "", err
		}
		navigationOutputs[name] = template.HTML(output)
	}
	templateData.Navigation = navigationOutputs

	//
	// Footer, Header may wish to access both content and navigation
	//

	footerOutput, err := r.RenderOne(fmt.Sprintf("%s-footer", r.TemplateName), r.Footer, templateData)
	if err != nil {
		return "", err
	}
	templateData.Footer = template.HTML(footerOutput)

	headerOutput, err := r.RenderOne(fmt.Sprintf("%s-header", r.TemplateName), r.Header, templateData)
	if err != nil {
		return "", err
	}
	templateData.Header = template.HTML(headerOutput)

	return r.RenderOne(r.TemplateName, r.TemplateContent, templateData)
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
