// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import (
	"bytes"
	"html/template"
)

type OptionAttrer interface {
	OptionAttrs() template.HTMLAttr
}

func hasAttr(v interface{}) bool {
	if _, ok := v.(OptionAttrer); ok {
		return true
	}
	return false
}

var tmpOptions = template.Must(template.New("options").Funcs(template.FuncMap{"hasAttr": hasAttr}).Parse(`
{{$sel := .Sel}}
{{range .Options}}
<option value='{{.OptionValue}}' {{if hasAttr . }}{{.OptionAttrs}}{{end}} {{if eq .OptionValue $sel}}selected{{end}}>{{.OptionLabel}}</option>
{{end}}
`)).Option("missingkey=error")

type OptionString struct {
	OptionValue string
	OptionLabel string
}
type OptionInt struct {
	OptionValue int
	OptionLabel string
}

func FmtOptions(slice interface{}, sel interface{}) template.HTML {
	var s bytes.Buffer
	err := tmpOptions.Execute(&s, map[string]interface{}{"Options": slice, "Sel": sel})
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(s.String())
}

var tmpOptionsGroup = template.Must(template.New("optionsGroup").Funcs(template.FuncMap{"hasAttr": hasAttr}).Parse(`
{{$sel := .Sel}}
{{$group := ""}}
{{range .Options}}
	{{if ne $group .OptionGroup}}
	{{if ne $group ""}} </optgroup> {{end}}
	<optgroup label='{{.OptionGroup}}'>
	{{$group = .OptionGroup}}
	{{end}}
<option value='{{.OptionValue}}' {{if hasAttr . }}{{.OptionAttrs}}{{end}} {{if eq .OptionValue $sel}}selected{{end}}>{{.OptionLabel}}</option>
{{end}}
</optgroup>
`)).Option("missingkey=error")

type OptionGroupString struct {
	OptionGroup string
	OptionValue string
	OptionLabel string
}
type OptionGroupInt struct {
	OptionGroup string
	OptionValue int
	OptionLabel string
}

func FmtOptionsGroup(slice interface{}, sel interface{}) template.HTML {
	var s bytes.Buffer
	err := tmpOptionsGroup.Execute(&s, map[string]interface{}{"Options": slice, "Sel": sel})
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(s.String())
}
