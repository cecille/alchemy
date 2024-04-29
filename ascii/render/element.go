package render

import (
	"fmt"
	"strings"

	"github.com/hasty/adoc/elements"
	"github.com/hasty/alchemy/internal/parse"
)

type Section interface {
	GetASCIISection() *elements.Section
}

func Elements(cxt *Context, prefix string, elementList ...elements.Element) (err error) {
	var previous any
	for _, e := range elementList {
		if he, ok := e.(Section); ok {
			e = he.GetASCIISection()
		}
		if hb, ok := e.(parse.HasBase); ok {
			e = hb.GetBase()
		}
		switch el := e.(type) {
		case *elements.Section:
			err = renderSection(cxt, el)
			if err == nil {
				err = Elements(cxt, "", el.Elements()...)
			}
		case *elements.Paragraph:
			cxt.WriteString(prefix)
			err = renderParagraph(cxt, el, &previous)
			if err != nil {
				return
			}
		case *elements.Table:
			err = renderTable(cxt, el)
		case *elements.EmptyLine:
			cxt.WriteNewline()
			cxt.WriteRune('\n')
		case *elements.CrossReference:
			err = renderInternalCrossReference(cxt, el)
		case *elements.AttributeEntry:
			err = renderAttributeEntry(cxt, el)
		case *elements.String:
			text := el.Value
			if strings.HasPrefix(text, "ifdef::") || strings.HasPrefix(text, "ifndef::") || strings.HasPrefix(text, "endif::[]") {
				cxt.WriteNewline()
			}
			cxt.WriteString(text)
		case *elements.SingleLineComment:
			cxt.WriteNewline()
			cxt.WriteString("//")
			cxt.WriteString(el.Value)
			cxt.WriteNewline()
		case *elements.BlockImage:
			err = renderImageBlock(cxt, el)
		case *elements.Link:
			err = renderLink(cxt, el)
		case *elements.SpecialCharacter:
			err = renderSpecialCharacter(cxt, el)
		case *elements.Bold:
			err = renderFormattedText(cxt, el, "*")
		case *elements.DoubleBold:
			err = renderFormattedText(cxt, el, "**")
		case *elements.Monospace:
			err = renderFormattedText(cxt, el, "`")
		case *elements.DoubleMonospace:
			err = renderFormattedText(cxt, el, "``")
		case *elements.Superscript:
			err = renderFormattedText(cxt, el, "^")
		case *elements.Subscript:
			err = renderFormattedText(cxt, el, "~")
		case *elements.Italic:
			err = renderFormattedText(cxt, el, "_")
		case *elements.DoubleItalic:
			err = renderFormattedText(cxt, el, "__")
		case *elements.Marked:
			err = renderFormattedText(cxt, el, "#")
		case *elements.DoubleMarked:
			err = renderFormattedText(cxt, el, "##")
		case *elements.LineContinuation:
			cxt.WriteString(" +")
		case elements.AttributeReference:
			cxt.WriteString(fmt.Sprintf("{%s}", el.Name()))
		case *elements.InlineImage:
			err = renderInlineImage(cxt, el)
		//case *elements.FootnoteReference:
		//	err = renderFootnoteReference(cxt, el)
		/*case *elements.InlinePassthrough:
		switch el.Kind {
		case elements.SinglePlusPassthrough, elements.TriplePlusPassthrough:
			cxt.WriteString(string(el.Kind))
			err = Elements(cxt, "", el.Elements)
			cxt.WriteString(string(el.Kind))
		case elements.PassthroughMacro:
			cxt.WriteString("pass:[")
			err = Elements(cxt, "", el.Elements)
			cxt.WriteRune(']')
		}*/
		case *elements.AttributeReset:
			renderAttributeReset(cxt, el)
		case *elements.UnorderedListItem:
			err = renderUnorderedListElement(cxt, el)
		case *elements.OrderedListItem:
			err = renderOrderedListElement(cxt, el)
		case *elements.ListContinuation:
			cxt.WriteNewline()
			cxt.WriteString("+\n")
			err = Elements(cxt, "", el.Child())
		case nil:
		default:
			err = fmt.Errorf("unknown render element type: %T", el)
		}
		if err != nil {
			return
		}
		previous = e
	}
	return
}
