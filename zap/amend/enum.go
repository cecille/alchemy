package amend

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hasty/alchemy/conformance"
	"github.com/hasty/alchemy/matter"
	"github.com/hasty/alchemy/parse"
	"github.com/hasty/alchemy/zap"
)

func (r *renderer) amendEnum(d xmlDecoder, e xmlEncoder, el xml.StartElement, cluster *matter.Cluster, clusterIDs []string, enums map[*matter.Enum]struct{}) (err error) {
	name := getAttributeValue(el.Attr, "name")

	var matchingEnum *matter.Enum
	for en := range enums {
		if en.Name == name || strings.TrimSuffix(en.Name, "Enum") == name {
			matchingEnum = en
			delete(enums, en)
			break
		}
	}
	Ignore(d, "enum")

	if matchingEnum == nil {
		return nil
	}

	return r.writeEnum(e, el, matchingEnum, clusterIDs, false)
}

func (r *renderer) writeEnum(e xmlEncoder, el xml.StartElement, en *matter.Enum, clusterIDs []string, provisional bool) (err error) {
	xfb := el.Copy()

	enumType := en.Type
	if enumType != "" {
		enumType = zap.ConvertDataTypeNameToZap(en.Type)
	} else {
		enumType = "enum8"
	}
	var valFormat string
	switch enumType {
	case "enum16":
		valFormat = "0x%04X"
	default:
		valFormat = "0x%02X"
	}

	xfb.Attr = setAttributeValue(xfb.Attr, "name", en.Name)
	xfb.Attr = setAttributeValue(xfb.Attr, "type", enumType)

	err = e.EncodeToken(xfb)
	if err != nil {
		return
	}
	err = r.renderClusterCodes(e, clusterIDs)
	if err != nil {
		return
	}

	for _, v := range en.Values {
		if conformance.IsZigbee(v.Conformance) {
			continue
		}

		val := v.Value
		valNum, er := parse.HexOrDec(val)
		if er == nil {
			val = fmt.Sprintf(valFormat, valNum)
		}

		elName := xml.Name{Local: "item"}
		xfs := xml.StartElement{Name: elName}
		xfs.Attr = setAttributeValue(xfs.Attr, "name", v.Name)
		xfs.Attr = setAttributeValue(xfs.Attr, "value", val)
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
