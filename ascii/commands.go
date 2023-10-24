package ascii

import (
	"log/slog"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/hasty/matterfmt/matter"
	"github.com/hasty/matterfmt/parse"
)

func (s *Section) toCommands() (commands []*matter.Command, err error) {
	var rows []*types.TableRow
	var headerRowIndex int
	var columnMap map[matter.TableColumn]int
	rows, headerRowIndex, columnMap, _, err = parseFirstTable(s)
	if err != nil {
		if err == NoTableFound {
			err = nil
		}
		return
	}
	commandMap := make(map[string]*matter.Command)
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		cmd := &matter.Command{}
		cmd.ID, err = readRowValue(row, columnMap, matter.TableColumnID)
		if err != nil {
			return
		}
		cmd.Name, err = readRowValue(row, columnMap, matter.TableColumnName)
		if err != nil {
			return
		}
		cmd.Direction, err = readRowValue(row, columnMap, matter.TableColumnDirection)
		if err != nil {
			return
		}
		cmd.Response, err = readRowValue(row, columnMap, matter.TableColumnResponse)
		if err != nil {
			return
		}
		cmd.Conformance, err = readRowValue(row, columnMap, matter.TableColumnConformance)
		if err != nil {
			return
		}
		var a string
		a, err = readRowValue(row, columnMap, matter.TableColumnAccess)
		if err != nil {
			return
		}
		cmd.Access = matter.ParseAccess(a)
		commands = append(commands, cmd)
		slog.Info("registering event", "event", cmd.Name)
		commandMap[cmd.Name] = cmd
	}

	for _, s := range parse.Skim[*Section](s.Elements) {
		switch s.SecType {
		case matter.SectionCommand:

			name := strings.TrimSuffix(s.Name, " Command")
			e, ok := commandMap[name]
			if !ok {
				slog.Info("unknown command", "command", name)
				continue
			}
			var rows []*types.TableRow
			var headerRowIndex int
			var columnMap map[matter.TableColumn]int
			rows, headerRowIndex, columnMap, _, err = parseFirstTable(s)
			if err != nil {
				if err == NoTableFound {
					err = nil
					continue
				}
				return
			}
			e.Fields, err = readFields(headerRowIndex, rows, columnMap)
		}
	}
	return
}
