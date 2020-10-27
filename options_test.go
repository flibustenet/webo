package webo

import (
	"html/template"
	"strings"
	"testing"
)

type O struct {
	v     int
	label string
	grp   string
	attr  string
}

func (o *O) OptionValue() int {
	return o.v
}
func (o *O) OptionLabel() string {
	return o.label
}
func (o *O) OptionAttrs() template.HTMLAttr {
	return template.HTMLAttr(o.attr)
}
func (o *O) OptionGroup() string {
	return o.grp
}

func equal(s, t string) bool {
	s = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
	t = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(t, " ", ""), "\n", ""), "\t", "")
	//	log.Println("S=", s)
	//	log.Println("T=", t)
	return s == t
}

func TestOptions(t *testing.T) {
	os := []*O{&O{1, "un", "", "class='cls' style='color:red'"}, &O{2, "deux", "", ""}}
	res := string(FmtOptions(os, 2))
	if !equal(res, "<option value='1' class='cls' style='color:red'>un</option><option value='2' selected>deux</option>") {
		t.Fatalf("Option : %s", res)
	}
}
func TestOptionsGroup(t *testing.T) {
	os := []*O{&O{1, "un", "g1", ""}, &O{2, "deux", "g1", "class='cls' style='color:red'"}, &O{3, "trois", "g2", ""}, &O{4, "quatre", "g2", ""}}
	res := string(FmtOptionsGroup(os, 2))
	if !equal(res, `
     	<optgroup label='g1'>
        <option value='1'  >un</option>
        <option value='2' class='cls' style='color:red' selected>deux</option>
     	 </optgroup> 
       	<optgroup label='g2'>
        <option value='3'  >trois</option>
        <option value='4'  >quatre</option>
        </optgroup>
`) {
		t.Fatalf("OptionGroup : %s", res)
	}
}
