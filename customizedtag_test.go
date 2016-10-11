package blackthunder

import (
	"testing"
	"bytes"
	"sort"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func testCTag(t *testing.T, input string, expected string, cTag map[string]CustomizedTag) {
	result := string(MarkdownWithCustomizedTag([]byte(input), cTag))
	if result != expected {
		diff := diffmatchpatch.New().DiffMain(result, expected, false)
		var buff bytes.Buffer
		for _, line := range diff {
			switch line.Type {
			case diffmatchpatch.DiffDelete:
				buff.WriteString("+")
				buff.WriteString(line.Text)
				buff.WriteString("+")
			case diffmatchpatch.DiffEqual:
				buff.WriteString(line.Text)
			case diffmatchpatch.DiffInsert:
				buff.WriteString("-")
				buff.WriteString(line.Text)
				buff.WriteString("-")
			}
		}
		t.Log(buff.String())
		t.Log(result)
		t.Fail()
	}
}

func sortMap(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestCTagInline(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["inline"] = CustomizedTag{
		Parse:func(attr map[string]string, args []string) CTagNode {
			var buff bytes.Buffer
			buff.WriteString("Attributes:\n")
			keys := sortMap(attr)
			for _, k := range keys {
				buff.WriteString(k)
				buff.WriteString(" : ")
				buff.WriteString(attr[k])
				buff.WriteString("\n")
			}
			buff.WriteString("Arguments:\n")
			for _, v := range args {
				buff.WriteString(v)
				buff.WriteString("\n")
			}

			return CTagNode{
				Content: buff.Bytes(),
			}
		},
	}

	input := `
{inline a=a 0 b=b name="Tanaka Satoshi" 1/}
{inline " A B C"/}
{inline /}
\{inline key=value args/}
{inline/}
{inline a=a/}

191919{inline 49494949 /}323232

`
	output := `<p>Attributes:
a : a
b : b
name : Tanaka Satoshi
Arguments:
0
1

Attributes:
Arguments:
 A B C

Attributes:
Arguments:

{inline key=value args/}
Attributes:
Arguments:

Attributes:
a : a
Arguments:
</p>

<p>191919Attributes:
Arguments:
49494949
323232</p>
`

	testCTag(t, input, output, tag)
}

func TestCTagInlineChild(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["red"] = CustomizedTag{
		Parse:func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Before: []byte(`<span style="color:red;">`),
				After: []byte(`</span>`),
			}
		},
	}

	input := `
{red}This is **Red**{/red}.

**123{red}456{/red}789**
`
	output := `<p><span style="color:red;">This is <strong>Red</strong></span>.</p>

<p><strong>123<span style="color:red;">456</span>789</strong></p>
`

	testCTag(t, input, output, tag)
}

