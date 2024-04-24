package disco

import (
	"github.com/hasty/adoc/elements"
	"github.com/hasty/alchemy/ascii"
	"github.com/hasty/alchemy/matter"
	"github.com/hasty/alchemy/matter/conformance"
)

func fixConformanceCells(doc *ascii.Doc, rows []*elements.TableRow, columnMap ascii.ColumnIndex) (err error) {
	if len(rows) < 2 {
		return
	}
	conformanceIndex, ok := columnMap[matter.TableColumnConformance]
	if !ok {
		return
	}
	for _, row := range rows[1:] {
		cell := row.Cells[conformanceIndex]
		vc, e := ascii.RenderTableCell(cell)
		if e != nil {
			continue
		}

		conf := conformance.ParseConformance(vc)

		cs := conf.ASCIIDocString()

		if cs != vc {
			err = setCellString(cell, cs)
			if err != nil {
				return
			}
		}

	}
	return
}
