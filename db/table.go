package db

import (
	"strings"

	"github.com/dolthub/go-mysql-server/memory"
	mms "github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/hasty/matterfmt/ascii"
	"github.com/hasty/matterfmt/matter"
	"github.com/iancoleman/strcase"
)

func (h *Host) createTable(cxt *mms.Context, tableName string, parentTable string, sections []*sectionInfo, columns []matter.TableColumn) error {
	schema, extra := buildTableSchema(sections, tableName, parentTable, columns)
	t := memory.NewTable(tableName, mms.NewPrimaryKeySchema(schema), h.db.GetForeignKeyCollection())
	h.tables[tableName] = t
	h.db.AddTable(tableName, t)
	err := populateTable(cxt, t, tableName, parentTable, sections, schema, columns, extra)
	if err != nil {
		return err
	}
	return nil
}

func buildTableSchema(sections []*sectionInfo, tableName string, parentName string, columns []matter.TableColumn) (mms.Schema, []ascii.ExtraColumn) {
	extraColumns := make(map[string]*ascii.ExtraColumn)
	for _, si := range sections {
		for e := range si.values.extras {
			extraColumns[strings.ToLower(e)] = &ascii.ExtraColumn{Name: e}
		}
	}
	schema := mms.Schema{
		{Name: tableName + "_id", Type: types.Int32, Nullable: false, Source: tableName, PrimaryKey: true},
	}
	if len(parentName) > 0 {
		schema = append(schema, &mms.Column{Name: parentName + "_id", Type: types.Int32, Nullable: false, Source: tableName, PrimaryKey: true})
	}

	for _, col := range columns {
		name, ok := matter.TableColumnNames[col]
		if !ok {
			continue
		}
		var columnName = strcase.ToSnake(name)
		var colType mms.Type = types.Text
		switch col {
		case matter.TableColumnAccess:
			schema = append(schema, getAccessSchemaColumns(tableName)...)
			continue
		case matter.TableColumnQuality:
			schema = append(schema, getQualitySchemaColumns(tableName)...)
			continue
		case matter.TableColumnID, matter.TableColumnBit, matter.TableColumnValue:
			colType = types.Int64
		case matter.TableColumnType:
			columnName = "data_type"
		}
		schema = append(schema, &mms.Column{Name: columnName, Type: colType, Nullable: true, Source: tableName, PrimaryKey: false})
	}

	offset := len(schema)
	var extra []ascii.ExtraColumn
	for _, e := range extraColumns {
		schema = append(schema, &mms.Column{Name: strcase.ToSnake(e.Name), Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: true})
		extra = append(extra, ascii.ExtraColumn{Name: e.Name})
		e.Offset = offset
		offset++
	}
	return schema, extra
}

func populateTable(cxt *mms.Context, t *memory.Table, tableName string, parentTable string, sections []*sectionInfo, schema mms.Schema, columns []matter.TableColumn, extra []ascii.ExtraColumn) error {
	for _, si := range sections {
		row := mms.NewRow(si.id)
		if len(parentTable) > 0 {
			row = append(row, si.parent.id)
		}
		for _, col := range columns {
			v, ok := si.values.values[col]
			switch col {
			case matter.TableColumnAccess:
				accessRows := getAccessSchemaColumnValues(tableName, v)
				row = append(row, accessRows...)
			case matter.TableColumnQuality:
				qualityRows := getQualitySchemaColumnValues(v)
				row = append(row, qualityRows...)
			default:
				if !ok {
					row = append(row, nil)
				} else {
					sr := schema[len(row)]
					var err error
					switch sr.Type {
					case types.Int64:
						switch val := v.(type) {
						case matter.DocType:
							v = int64(val)
						case string:
							v, err = parseNumber(val)
						}
					}
					if err != nil {
						v = nil
					}
					row = append(row, v)

				}
			}
		}
		for _, e := range extra {
			v, ok := si.values.extras[e.Name]
			if ok {
				row = append(row, v)
			} else {
				//slog.Info("missing val", "name", e.Name)
				row = append(row, nil)
			}
		}
		err := t.Insert(cxt, row)
		if err != nil {
			return err
		}

	}
	return nil
}
