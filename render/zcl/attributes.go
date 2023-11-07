package zcl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/hasty/alchemy/matter"
	"github.com/hasty/alchemy/parse"
	"github.com/iancoleman/strcase"
)

func renderAttributes(cluster *matter.Cluster, cx *etree.Element, clusterPrefix string, errata *Errata) {
	if len(cluster.Attributes) > 0 {
		cx.CreateComment("Attributes")
	}
	for _, a := range cluster.Attributes {
		if a.Conformance == "Zigbee" {
			continue
		}
		mandatory := (a.Conformance == "M")

		if !a.ID.Valid() {
			continue
		}

		if !mandatory && len(a.Conformance) > 0 && a.Conformance != "O" {
			cx.CreateComment(fmt.Sprintf("Conformance feature %s - for now optional", a.Conformance))
		}
		attr := cx.CreateElement("attribute")
		attr.CreateAttr("code", a.ID.HexString())
		attr.CreateAttr("side", "server")
		writeDataType(attr, a.Type)
		define := GetDefine(a.Name, clusterPrefix, errata)
		attr.CreateAttr("define", define)
		if a.Quality.Has(matter.QualityNullable) {
			attr.CreateAttr("isNullable", "true")
		} else {
			attr.CreateAttr("isNullable", "false")
		}
		if a.Quality.Has(matter.QualityReportable) {
			attr.CreateAttr("reportable", "true")
		}
		if a.Default != "" {
			switch a.Default {
			case "null":
				switch a.Type.Name {
				case "uint8":
					attr.CreateAttr("default", "0xFF")
				case "uint16":
					attr.CreateAttr("default", "0xFFFF")
				case "uint32":
					attr.CreateAttr("default", "0xFFFFFFFF")
				case "uint64":
					attr.CreateAttr("default", "0xFFFFFFFFFFFFFFFF")
				}
			default:
				def, err := parse.HexOrDec(a.Default)
				if err == nil {
					attr.CreateAttr("default", strconv.Itoa(int(def)))
				}
			}
		}
		renderConstraint(a.Constraint, attr)
		renderAttributeAccess(a, errata, attr)
		if a.Conformance != "M" {
			attr.CreateAttr("optional", "true")
		} else {
			attr.CreateAttr("optional", "false")
		}

	}
}

func renderConstraint(a matter.Constraint, attr *etree.Element) {
	switch c := a.(type) {
	case *matter.RangeConstraint:
		attr.CreateAttr("min", c.Min.ZCLString())
		attr.CreateAttr("max", c.Max.ZCLString())
	case *matter.MinConstraint:
		attr.CreateAttr("min", c.Min.ZCLString())
	case *matter.MaxConstraint:
		attr.CreateAttr("max", c.Max.ZCLString())
	case *matter.MaxLengthConstraint:
		attr.CreateAttr("length", c.Length.ZCLString())
	case *matter.MinLengthConstraint:
		attr.CreateAttr("minLength", c.Length.ZCLString())
	case *matter.LengthRangeConstraint:
		attr.CreateAttr("length", c.Max.ZCLString())
		attr.CreateAttr("minLength", c.Min.ZCLString())
	}
}

func renderAttributeAccess(a *matter.Field, errata *Errata, attr *etree.Element) {
	if a.Quality.Has(matter.QualityFixed) || (a.Access.Read == matter.PrivilegeUnknown || a.Access.Read == matter.PrivilegeView) && a.Access.Write == matter.PrivilegeUnknown || errata.SuppressAttributePermissions {
		if a.Access.Write != matter.PrivilegeUnknown {
			attr.CreateAttr("writable", "true")
		} else {
			attr.CreateAttr("writable", "false")
		}
		attr.SetText(a.Name)
	} else {
		if a.Access.Read != matter.PrivilegeUnknown {
			ax := attr.CreateElement("access")
			ax.CreateAttr("op", "read")
			ax.CreateAttr("privilege", renderPrivilege(a.Access.Read))
		}
		if a.Access.Write != matter.PrivilegeUnknown {
			ax := attr.CreateElement("access")
			ax.CreateAttr("op", "write")
			ax.CreateAttr("privilege", renderPrivilege(a.Access.Write))
			attr.CreateAttr("writable", "true")
		} else {
			attr.CreateAttr("writable", "false")
		}
		attr.CreateElement("description").SetText(a.Name)
	}
}

func renderPrivilege(a matter.Privilege) string {
	switch a {
	case matter.PrivilegeView:
		return "view"
	case matter.PrivilegeManage:
		return "manage"
	case matter.PrivilegeAdminister:
		return "administer"
	case matter.PrivilegeOperate:
		return "operate"
	default:
		return ""
	}
}

func GetDefine(name string, prefix string, errata *Errata) string {
	define := strcase.ToScreamingDelimited(name, '_', "", true)
	if !strings.HasPrefix(define, prefix) {
		define = prefix + define
	}
	if errata.DefineOverrides != nil {
		if override, ok := errata.DefineOverrides[define]; ok {
			return override
		}
	}
	return define
}
