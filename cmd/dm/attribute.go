package dm

import (
	"github.com/beevik/etree"
	"github.com/hasty/alchemy/matter"
)

func renderAttributes(cluster *matter.Cluster, c *etree.Element) (err error) {
	if len(cluster.Attributes) == 0 {
		return
	}
	attributes := c.CreateElement("attributes")
	for _, a := range cluster.Attributes {
		ax := attributes.CreateElement("attribute")
		ax.CreateAttr("id", a.ID.HexString())
		ax.CreateAttr("name", a.Name)
		renderDataType(a, ax)
		if len(a.Default) > 0 {
			ax.CreateAttr("default", a.Default)
		}
		renderAccess(ax, a)
		renderQuality(ax, a)
		err = renderConformanceString(a.Conformance, ax)
		if err != nil {
			return
		}

		err = renderConstraint(a.Constraint, a.Type, ax)
		if err != nil {
			return
		}

	}
	return
}
