package amend

import (
	"encoding/xml"
	"strings"

	"github.com/hasty/alchemy/conformance"
	"github.com/hasty/alchemy/matter"
)

func (r *renderer) amendStruct(d xmlDecoder, e xmlEncoder, el xml.StartElement, cluster *matter.Cluster, clusterIDs []string, structs map[*matter.Struct]bool) (err error) {
	name := getAttributeValue(el.Attr, "name")

	var skip bool
	var matchingStruct *matter.Struct
	for s, handled := range structs {
		if s.Name == name || strings.TrimSuffix(s.Name, "Struct") == name {
			matchingStruct = s
			skip = handled
			structs[s] = true
			break
		}
	}

	Ignore(d, "struct")

	if matchingStruct == nil || skip {
		return
	}

	if r.errata.SeparateStructs != nil {
		if _, ok := r.errata.SeparateStructs[name]; ok {
			for _, clusterID := range clusterIDs {
				err = r.writeStruct(e, el, matchingStruct, []string{clusterID}, false)
				if err != nil {
					return
				}
			}
			return
		}
	}

	return r.writeStruct(e, el, matchingStruct, clusterIDs, false)
}

func (r *renderer) writeStruct(e xmlEncoder, el xml.StartElement, s *matter.Struct, clusterIDs []string, provisional bool) (err error) {
	xfb := el.Copy()
	xfb.Name = xml.Name{Local: "struct"}
	xfb.Attr = setAttributeValue(xfb.Attr, "name", s.Name)
	if provisional {
		xfb.Attr = setAttributeValue(xfb.Attr, "apiMaturity", "provisional")
	}
	if s.FabricScoped {
		xfb.Attr = setAttributeValue(xfb.Attr, "isFabricScoped", "true")
	} else {
		xfb.Attr = removeAttribute(xfb.Attr, "isFabricScoped")
	}
	err = e.EncodeToken(xfb)
	if err != nil {
		return
	}

	err = r.renderClusterCodes(e, clusterIDs)
	if err != nil {
		return
	}

	for _, v := range s.Fields {
		if conformance.IsZigbee(v.Conformance) {
			continue
		}

		elName := xml.Name{Local: "item"}
		xfs := xml.StartElement{Name: elName}
		xfs.Attr = setAttributeValue(xfs.Attr, "fieldId", v.ID.IntString())
		xfs.Attr = setAttributeValue(xfs.Attr, "name", v.Name)
		xfs.Attr = writeDataType(s.Fields, v, xfs.Attr)
		xfs.Attr = r.renderConstraint(s.Fields, v, xfs.Attr)
		if v.Quality.Has(matter.QualityNullable) {
			xfs.Attr = setAttributeValue(xfs.Attr, "isNullable", "true")
		} else {
			xfs.Attr = removeAttribute(xfs.Attr, "isNullable")
		}
		if !conformance.IsMandatory(v.Conformance) {
			xfs.Attr = setAttributeValue(xfs.Attr, "optional", "true")
		} else {
			xfs.Attr = removeAttribute(xfs.Attr, "optional")
		}
		if v.Access.FabricSensitive {
			xfs.Attr = setAttributeValue(xfs.Attr, "isFabricSensitive", "true")
		} else {
			xfs.Attr = removeAttribute(xfs.Attr, "isFabricSensitive")
		}
		err = e.EncodeToken(xfs)
		if err != nil {
			return
		}
		xfe := xml.EndElement{Name: elName}
		err = e.EncodeToken(xfe)
		if err != nil {
			return
		}

	}
	return e.EncodeToken(xml.EndElement{Name: xfb.Name})
}
