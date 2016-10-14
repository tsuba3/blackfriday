Blackthunder
===========

Blackthunder added customized tags to
[Blackfriday](https://github.com/russross/blackfriday).

Blackfriday is a [Markdown][1] processor implemented in [Go][2]. It
is paranoid about its input (so you can safely feed it user-supplied
data), it is fast, it supports common extensions (tables, smart
punctuation substitutions, etc.), and it is safe for all utf-8
(unicode) input.

HTML output is currently supported, along with Smartypants
extensions.

It started as a translation from C of [Sundown][3].


Installation
------------

    go get github.com/tsuba3/blackthunder

Usage
-----

``` go
    tags = map[string]blackthunder.CustomizedTag{}
    tags["red"] = blackthunder.CustomizedTag{
        Parse: func(attr map[string]string, args []string) blackthunder.CTagNode {
            return blackthunder.CTagNode{
                Before: []byte(`<span style="color:red;">`),
                After:  []byte(`</span>`),
            }
        }
    }

    input := []byte(`{red}This is *RED* text.{/red}`)
    output := blackthunder.MarkdownWithCustomizedTag(input, tags)
    // output:
    // <p><span style="color:red;">This is <em>RED</em> text.</span></p>
```

### Sanitize untrusted content

Blackthunder itself does nothing to protect against malicious content. If you are
dealing with user-supplied markdown, we recommend running blackthunder's output
through HTML sanitizer such as
[Bluemonday](https://github.com/microcosm-cc/bluemonday).

Here's an example of simple usage of blackthunder together with bluemonday:

``` go
import (
    "github.com/tsuba3/blackthunder"
    "github.com/russross/blackfriday"
)

// ...
unsafe := blackfriday.MarkdownWithCustomizedTag(input, tag)
html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
```
Features
--------

[All features and syntax of blackfriday](https://github.com/russross/blackfriday#features)
are supported.

Customized tags
--------

A customized tag is like XHTML tags e.g.  
`{name key=value 0}...{/name}`  
`key=value` is an attribute. `0` is an argument.  
A tag without children like `{name /}` requires to end with `/}`.

``` go
tags = map[string]blackthunder.CustomizedTag{}
tags["red"] = blackthunder.CustomizedTag{
    Async: false, // if true, Parse will run async
    Parse: func(attr map[string]string, args []string) blackthunder.CTagNode {
        return blackthunder.CTagNode{
            Before: []byte(`<span style="color:red;">`),
            After:  []byte(`</span>`),
        }
    }
}
```

Parse function recieve attributes and arguments and return CTagNode.
CTagNode is
``` go
type CTagNode struct {
	IsBlock bool // Block will not be wrapped by <p>

	Child map[string]CustomizedTag // Tags available in children

	// For tags with children
	Before []byte // before children
	After  []byte // after children

	// For tags without children
	Content []byte
}
```



License
-------

[Blackthunder is distributed under the Simplified BSD License](LICENSE.txt)


   [1]: http://daringfireball.net/projects/markdown/ "Markdown"
   [2]: http://golang.org/ "Go Language"
   [3]: https://github.com/vmg/sundown "Sundown"
