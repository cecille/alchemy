package ascii

import (
	"log/slog"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/hasty/alchemy/parse"
)

type Anchor struct {
	ID      string
	Label   string
	Element types.WithAttributes
	Parent  parse.HasElements
	Name    string
}

func (doc *Doc) Anchors() (map[string]*Anchor, error) {
	if doc.anchors != nil {
		return doc.anchors, nil
	}
	anchors := make(map[string]*Anchor)
	crossReferences := doc.CrossReferences()
	parse.Traverse(doc, doc.Elements, func(el interface{}, parent parse.HasElements, index int) bool {
		var wa types.WithAttributes
		e, ok := el.(*Element)
		if ok {
			if wa, ok = e.Base.(types.WithAttributes); !ok {
				return false
			}
		} else if s, ok := el.(*Section); ok {
			wa = s.Base
		} else {
			return false
		}
		attrs := wa.GetAttributes()
		idAttr, ok := attrs[types.AttrID]
		if !ok {
			return false
		}
		id := strings.TrimSpace(idAttr.(string))
		var label string
		if parts := strings.Split(id, ","); len(parts) > 1 {
			id = strings.TrimSpace(parts[0])
			label = strings.TrimSpace(parts[1])
		}
		reftext, ok := attrs.GetAsString("reftext")
		if ok {
			label = reftext
		}
		info := &Anchor{
			ID:      id,
			Label:   label,
			Element: wa,
			Parent:  parent,
		}
		name := ReferenceName(wa)
		if name != "" {
			info.Name = name
		} else if len(label) > 0 {
			info.Name = label
		}
		if _, ok := anchors[id]; ok {
			slog.Warn("duplicate anchor; can't fix", "id", id)
			return false
		}

		if !strings.HasPrefix(id, "_") {
			anchors[id] = info
		} else { // Anchors prefaced with "_" may have been created by the parser
			if _, ok := crossReferences[id]; ok { // If there's a cross-reference for it, then we'll render it
				anchors[id] = info
			} else { // If there isn't a cross reference to the id, there might be one to its original version
				unescaped := strings.TrimSpace(strings.ReplaceAll(id, "_", " "))
				if _, ok = crossReferences[unescaped]; ok {
					if _, ok := anchors[unescaped]; ok {
						slog.Warn("duplicate anchor; can't fix", "id", unescaped)
						return false
					}
					anchors[unescaped] = info
				}
			}
		}
		return false
	})
	doc.anchors = anchors
	return anchors, nil
}
