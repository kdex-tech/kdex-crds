package render

import (
	"bytes"
	"fmt"
	"html/template"
	"maps"
	"reflect"
	"sort"
	texttemplate "text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/text/language"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (r *Renderer) RenderPage() (string, error) {
	templateData, err := r.TemplateData()

	if err != nil {
		return "", err
	}

	return r.RenderOne(r.TemplateName, r.TemplateContent, templateData)
}

func (r *Renderer) TemplateData() (TemplateData, error) {
	date := r.LastModified
	if date.IsZero() {
		date = time.Now()
	}

	pageMap := map[string]any{}
	if r.PageMap != nil {
		maps.Copy(pageMap, r.PageMap)
	}

	templateData := TemplateData{
		BasePath:        r.BasePath,
		BrandName:       r.BrandName,
		DefaultLanguage: r.DefaultLanguage,
		Extra:           r.Extra,
		Host:            r.Host,
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
		return TemplateData{}, err
	}
	templateData.FootScript = template.HTML(footerScriptOutput)

	headerScriptOutput, err := r.RenderOne(fmt.Sprintf("%s-headscript", r.TemplateName), r.HeadScript, templateData)
	if err != nil {
		return TemplateData{}, err
	}
	templateData.HeadScript = template.HTML(headerScriptOutput)

	metaOutput, err := r.RenderOne(fmt.Sprintf("%s-meta", r.TemplateName), r.Meta, templateData)
	if err != nil {
		return TemplateData{}, err
	}
	templateData.Meta = template.HTML(metaOutput)

	themeOutput, err := r.RenderOne(fmt.Sprintf("%s-theme", r.TemplateName), r.Theme, templateData)
	if err != nil {
		return TemplateData{}, err
	}
	templateData.Theme = template.HTML(themeOutput)

	//
	// Content and Navigation
	//

	contentOutputs := make(map[string]template.HTML)
	for slot, content := range r.Contents {
		output, err := r.RenderOne(fmt.Sprintf("%s-content-%s", r.TemplateName, slot), content, templateData)
		if err != nil {
			return TemplateData{}, err
		}

		contentOutputs[slot] = template.HTML(output)
	}
	templateData.Content = contentOutputs

	navigationOutputs := make(map[string]template.HTML)
	for name, content := range r.Navigations {
		output, err := r.RenderOne(fmt.Sprintf("%s-navigation-%s", r.TemplateName, name), content, templateData)
		if err != nil {
			return TemplateData{}, err
		}
		navigationOutputs[name] = template.HTML(output)
	}
	templateData.Navigation = navigationOutputs

	//
	// Footer, Header may wish to access both content and navigation
	//

	footerOutput, err := r.RenderOne(fmt.Sprintf("%s-footer", r.TemplateName), r.Footer, templateData)
	if err != nil {
		return TemplateData{}, err
	}
	templateData.Footer = template.HTML(footerOutput)

	headerOutput, err := r.RenderOne(fmt.Sprintf("%s-header", r.TemplateName), r.Header, templateData)
	if err != nil {
		return TemplateData{}, err
	}
	templateData.Header = template.HTML(headerOutput)

	return templateData, nil
}

// funcMap builds the function map shared by both RenderOne (html/template)
// and RenderOneText (text/template). html/template.FuncMap is a type alias
// for text/template.FuncMap (both are map[string]any under the hood), so a
// single map value works unmodified with either engine's .Funcs() call.
func (r *Renderer) funcMap() texttemplate.FuncMap {
	funcs := sprig.FuncMap()
	funcs["extract"] = func(key string, v any) ([]any, error) {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return nil, fmt.Errorf("value is not a slice or array")
		}
		res := []any{}
		for i := 0; i < rv.Len(); i++ {
			itemValue := rv.Index(i)
			for itemValue.Kind() == reflect.Interface || itemValue.Kind() == reflect.Pointer {
				itemValue = itemValue.Elem()
			}
			if itemValue.Kind() != reflect.Struct {
				return nil, fmt.Errorf("item is not a struct")
			}
			if val := itemValue.FieldByName(key); val.IsValid() {
				res = append(res, val.Interface())
			} else {
				return nil, fmt.Errorf("field %s not found in struct", key)
			}
		}
		return res, nil
	}
	funcs["l10n"] = func(key string, args ...any) string {
		if r.MessagePrinter == nil {
			return key
		}
		return r.MessagePrinter.Sprintf(key, args...)
	}

	tag, err := language.Parse(r.Language)
	if err != nil {
		tag = language.English
	}

	funcs["number"] = func(v any) string {
		return formatNumber(tag, v)
	}
	funcs["currency"] = func(v float64, code string) string {
		return formatCurrency(tag, v, code)
	}
	funcs["percent"] = func(v float64) string {
		return formatPercent(tag, v)
	}
	funcs["bytes"] = func(v float64, baseUnit string) string {
		return formatBytes(tag, v, baseUnit)
	}
	funcs["date"] = func(t time.Time, style string) string {
		return formatDate(tag, t, style)
	}
	funcs["pop"] = func(v any, key string) any {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Map || rv.IsNil() {
			return ""
		}
		rk := reflect.ValueOf(key)
		val := rv.MapIndex(rk)
		if val.IsValid() {
			rv.SetMapIndex(rk, reflect.Value{})
			return val.Interface()
		}
		return ""
	}
	funcs["sortBy"] = func(field string, ascending bool, v any) ([]any, error) {
		if v == nil {
			return nil, nil
		}
		tp := reflect.TypeOf(v).Kind()
		switch tp {
		case reflect.Slice, reflect.Array:
			l2 := reflect.ValueOf(v)

			l := l2.Len()
			nl := make([]any, l)
			for i := range l {
				nl[i] = l2.Index(i).Interface()
			}

			sort.Slice(nl, func(i, j int) bool {
				val1 := reflect.ValueOf(nl[i])
				val2 := reflect.ValueOf(nl[j])

				if !ascending {
					val1, val2 = val2, val1
				}

				if val1.Kind() == reflect.Pointer {
					val1 = val1.Elem()
				}
				if val2.Kind() == reflect.Pointer {
					val2 = val2.Elem()
				}

				f1 := val1.FieldByName(field)
				f2 := val2.FieldByName(field)

				if !f1.IsValid() || !f2.IsValid() {
					return false
				}

				if f1.Type() == reflect.TypeFor[resource.Quantity]() {
					q1 := f1.Interface().(resource.Quantity)
					q2 := f2.Interface().(resource.Quantity)
					return q1.Cmp(q2) < 0
				}

				switch f1.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					return f1.Int() < f2.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return f1.Uint() < f2.Uint()
				case reflect.Float32, reflect.Float64:
					return f1.Float() < f2.Float()
				case reflect.String:
					return f1.String() < f2.String()
				}

				return false
			})

			return nl, nil
		default:
			return nil, fmt.Errorf("cannot sort on type %s by field %s", tp, field)
		}
	}

	return funcs
}

// RenderOne renders templateContent using html/template, so an action's
// output is HTML-escaped based on the surrounding markup context. Use this
// for HTML pages/fragments, which depend on that escaping.
func (r *Renderer) RenderOne(
	templateName string,
	templateContent string,
	data TemplateData,
) (string, error) {
	instance, err := template.New(templateName).Delims("[[", "]]").Funcs(r.funcMap()).Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := instance.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// RenderOneText renders templateContent using text/template, so an action's
// output is emitted verbatim with no HTML auto-escaping. Use this for
// non-HTML bodies (e.g. robots.txt, llms.txt, JSON, markdown), where
// html/template's escaping would corrupt the output (or, for a structured
// body, could trip an html/template context error at Execute).
func (r *Renderer) RenderOneText(
	templateName string,
	templateContent string,
	data TemplateData,
) (string, error) {
	instance, err := texttemplate.New(templateName).Delims("[[", "]]").Funcs(r.funcMap()).Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := instance.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
