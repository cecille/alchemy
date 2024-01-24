package matter

import (
	"log/slog"

	"github.com/hasty/alchemy/matter/types"
)

type Cluster struct {
	ID          *Number     `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Revisions   []*Revision `json:"revisions,omitempty"`
	Base        bool        `json:"baseCluster,omitempty"`

	Hierarchy string `json:"hierarchy,omitempty"`
	Role      string `json:"role,omitempty"`
	Scope     string `json:"scope,omitempty"`
	PICS      string `json:"pics,omitempty"`

	Features   *Bitmap    `json:"features,omitempty"`
	Bitmaps    []*Bitmap  `json:"bitmaps,omitempty"`
	Enums      []*Enum    `json:"enums,omitempty"`
	Structs    []*Struct  `json:"structs,omitempty"`
	Attributes FieldSet   `json:"attributes,omitempty"`
	Events     EventSet   `json:"events,omitempty"`
	Commands   CommandSet `json:"commands,omitempty"`
}

func (c *Cluster) EntityType() types.EntityType {
	return types.EntityTypeCluster
}

func (c *Cluster) Inherit(parent *Cluster) (err error) {
	slog.Info("Inheriting cluster", "parent", parent.Name, "child", c.Name)
	if parent.Features != nil {
		if c.Features == nil || len(c.Features.Bits) == 0 {
			c.Features = parent.Features.Clone()
		} else {
			err = c.Features.Inherit(parent.Features)
			if err != nil {
				return
			}
		}
	}

	if len(c.Description) == 0 {
		c.Description = parent.Description
	}

	c.Attributes = c.Attributes.Inherit(parent.Attributes)

	for _, pbm := range parent.Bitmaps {
		var matching *Bitmap
		for _, b := range c.Bitmaps {
			if b.Name == pbm.Name {
				matching = b
				break
			}
		}
		if matching == nil {
			c.Bitmaps = append(c.Bitmaps, pbm.Clone())
			continue
		}
		err = matching.Inherit(pbm)
		if err != nil {
			return
		}
	}

	for _, pe := range parent.Enums {
		var matching *Enum
		for _, en := range c.Enums {
			if en.Name == pe.Name {
				matching = en
				break
			}
		}
		if matching == nil {
			c.Enums = append(c.Enums, pe.Clone())
			continue
		}
		err = matching.Inherit(pe)
		if err != nil {
			return
		}
	}

	for _, ps := range parent.Structs {
		var matching *Struct
		for _, s := range c.Structs {
			if s.Name == ps.Name {
				matching = s
				break
			}
		}
		if matching == nil {
			c.Structs = append(c.Structs, ps.Clone())
			continue
		}
		matching.Inherit(ps)
	}

	for _, pe := range parent.Events {
		var matching *Event
		for _, e := range c.Events {
			if e.ID.Equals(pe.ID) {
				matching = e
				break
			}
		}
		if matching == nil {
			c.Events = append(c.Events, pe.Clone())
			continue
		}
		matching.Inherit(pe)
	}

	for _, pc := range parent.Commands {
		var matching *Command
		for _, c := range c.Commands {
			if c.ID.Equals(pc.ID) {
				matching = c
				break
			}
		}
		if matching == nil {
			c.Commands = append(c.Commands, pc.Clone())
			continue
		}
		matching.Inherit(pc)
	}

	return nil
}

func (c *Cluster) Reference(name string) types.Entity {
	if c == nil {
		return nil
	}
	var cr types.Entity
	if c.Features != nil {
		cr = c.Features.Reference(name)
		if cr, ok := cr.(*Bit); ok && cr != nil {
			return cr
		}

	}
	cr = c.Attributes.Reference(name)
	if cr, ok := cr.(*Field); ok && cr != nil {
		return cr
	}
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return cmd
		}
	}
	for _, e := range c.Events {
		if e.Name == name {
			return e
		}
	}
	for _, e := range c.Enums {
		if e.Name == name {
			return e
		}
	}
	for _, e := range c.Bitmaps {
		if e.Name == name {
			return e
		}
	}
	for _, e := range c.Structs {
		if e.Name == name {
			return e
		}
	}
	return nil
}
