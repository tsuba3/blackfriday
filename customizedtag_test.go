package blackthunder

import (
	"bytes"
	"github.com/sergi/go-diff/diffmatchpatch"
	"runtime/debug"
	"sort"
	"strconv"
	"testing"
	"time"
)

func testCTag(t *testing.T, input string, expected string, cTag map[string]CustomizedTag, msg string) {
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
		t.Log(msg)
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
		Parse: func(attr map[string]string, args []string) CTagNode {
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

{inline 4989/}

{inline A B C/}

{inline A   b/}

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

<p>Attributes:
Arguments:
4989
</p>
<p>Attributes:
Arguments:
A
B
C
</p>
<p>Attributes:
Arguments:
A
b
</p>
`

	testCTag(t, input, output, tag, "")
}

func TestCTagInlineChild(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["red"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Before: []byte(`<span style="color:red;">`),
				After:  []byte(`</span>`),
			}
		},
	}

	input := `
{red}This is **Red**{/red}.

**123{red}456{/red}789**

{red}{/red}

Hello {red}Tanaka{/red}!

A{red}B{red}CD{/red}E{red}{/red}{red}F{red}G{/red}{/red}H{red}I{/red}J{/red}

`
	output := `<p><span style="color:red;">This is <strong>Red</strong></span>.</p>

<p><strong>123<span style="color:red;">456</span>789</strong></p>

<p><span style="color:red;"></span></p>
<p>Hello <span style="color:red;">Tanaka</span>!</p>

<p>A<span style="color:red;">B<span style="color:red;">CD</span>E<span style="color:red;"></span><span style="color:red;">F<span style="color:red;">G</span></span>H<span style="color:red;">I</span>J</span></p>
`

	testCTag(t, input, output, tag, "")
}

func TestCTagBlock(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["div"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Before:  []byte(`<div>`),
				After:   []byte(`</div>`),
				IsBlock: true,
			}
		},
	}
	tag["span"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Before: []byte(`<span>`),
				After:  []byte(`</span>`),
			}
		},
	}
	tag["br"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Content: []byte(`<br>`),
			}
		},
	}

	input1 := `
{div}

A

**B**

{span}SPAN{br/}{/span}

{br/}

{/div}

`

	output1 := `<div>

A

<strong>B</strong>
<span>SPAN<br></span>
<br>
</div>`

	testCTag(t, input1, output1, tag, "div")

	input2 := `
{div}{/div}

{span}{/span}

{br/}
`

	output2 := `<div></div><p><span></span></p>
<p><br></p>
`
	testCTag(t, input2, output2, tag, "block")

}

func TestCTagChild(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["block"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{
				Before:  []byte("<div>"),
				After:   []byte("</div>"),
				IsBlock: true,
				Child: map[string]CustomizedTag{
					"name": {
						Parse: func(attr map[string]string, args []string) CTagNode {
							return CTagNode{Content: []byte("Takeshi")}
						},
					},
					"age": {
						Parse: func(attr map[string]string, args []string) CTagNode {
							return CTagNode{Content: []byte("42")}
						},
					},
				},
			}
		},
	}
	tag["name"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{Content: []byte("NAME")}
		},
	}

	input := `
{block}
	Name {name/}
	Age **{age/}**
{/block}

{name /}
{age /}
**{name /}**

`

	output := `<div>
    Name Takeshi
    Age **42**</div><p>NAME

<strong>NAME</strong></p>
`

	testCTag(t, input, output, tag, "child")

}

func TestCTagAsync(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["async"] = CustomizedTag{
		Async: true,
		Parse: func(attr map[string]string, args []string) CTagNode {
			if len(args) < 1 {
				return CTagNode{}
			}
			d, _ := strconv.Atoi(args[0])
			time.Sleep(time.Duration(d) * time.Millisecond)
			return CTagNode{
				Before:  []byte{'<'},
				After:   []byte{'>'},
				Content: []byte(args[0]),
			}
		},
	}

	input := `
{async 500}
**Hello**
{async 500/}
{/async}

{async 100/}

{async 200}
**Hello**
{async 800/}
{async/}

{async 1000/}
`
	output := `<p><
<strong>Hello</strong>500></p>
<p>100</p>
<p><
<strong>Hello</strong>800></p>

<p>1000</p>
`

	testCTag(t, input, output, tag, "child")
}

// should not panic.
func TestCTagError(t *testing.T) {
	tag := map[string]CustomizedTag{}
	tag["a"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{}
		},
	}
	tag["b"] = CustomizedTag{
		Parse: func(attr map[string]string, args []string) CTagNode {
			return CTagNode{IsBlock: true}
		},
	}

	inputs := []string{
		"{",
		"{/",
		"{\n",
		"{/\n",
		"a{\n",
		"{a",
		"{a\n",
		"{a}AAA\n\n",
		"*AA{*\n",
		"{}{}\n",
		"{{a}}\n",
		"{L}{/L}",
		"{a}{/a}",
		"{a}{/}{/}\n\n",
		"{/a}\n",
		"A{/a}\n",
		"{a/}",
		"{a=",
		"{a b/}",
		"{a v=",
		"{a/",
		"{b /}",
		"a\n{b}{/}",
		"{b}{a}{/b}{/a}",
		"**{b}{/b}**",
		"{/b}",
		"A{b}/{",
		"{b \"",
		"{a \"}",
		"{b k=\"",
	}

	i := 0
	defer func() {
		e := recover()
		if e != nil {
			t.Log(e)
			t.Log(inputs[i])
			debug.PrintStack()
			t.Fail()
		}
	}()
	for ; i < len(inputs); i++ {
		MarkdownWithCustomizedTag([]byte(inputs[i]), tag)
	}

}
