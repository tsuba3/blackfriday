package blackthunder

import (
	"testing"
	"bytes"
)

func TestCTagInline(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["inline"] = CustomizedTag{
		Parse:func(attr map[string]string, args []string, child []byte) CTagNode {
			var buff bytes.Buffer
			buff.WriteString("Attributes:\n")
			for k, v := range attr {
				buff.WriteString(k)
				buff.WriteString(" : ")
				buff.WriteString(v)
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
\{inline key=value args/}`
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

{inline key=value args/}</p>
`

	result := string(MarkdownWithCustomizedTag([]byte(input), tag))
	if result != output {
		t.Log(result)
		t.Fail()
	}
}

func TestCTagInlineChild(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["red"] = CustomizedTag{
		HasChild: true,
		Parse:func(attr map[string]string, args []string, child []byte) CTagNode {
			return CTagNode{
				Before: []byte(`<span style="color:red;">`),
				After: []byte(`</span>`),
			}
		},
	}

	input := `
{red}This is **Red**.{/red}
`
	output := `<p><span style="color:red;">This is <em>Red</em>.</span></p>
`

	result := string(MarkdownWithCustomizedTag([]byte(input), tag))
	if result != output {
		t.Log(result)
		t.Fail()
	}
}

