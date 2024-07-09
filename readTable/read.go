package readtable

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

const (
	String = iota
	Int
	Float
)

type Column struct {
	ColimnName string
	ColumnType int
}

func ReadTableData(arg1 *godog.Table, fields []Column) ([]map[string]interface{}, error) {

	var result = []map[string]interface{}{}
	columnNames := arg1.Rows[0]
	// modify column names to json compatible format like "ticket number" -->ticket_number
	modifiedColumnNames := make([]string, 0)
	for _, v := range columnNames.Cells {
		small := strings.ToLower(v.Value)
		split := strings.Split(small, " ")
		if len(split) > 1 {
			small = strings.Join(split, "_")
			modifiedColumnNames = append(modifiedColumnNames, small)
		} else {
			modifiedColumnNames = append(modifiedColumnNames, split[0])
		}
	}
	fmt.Println("modified column Names:", modifiedColumnNames)

	// convert tables of data into []map[string]interface{}{}
	for i, r := range arg1.Rows {
		if i == 0 {
			continue
		}
		row := make(map[string]interface{})
		for i, v := range modifiedColumnNames {
			for _, field := range fields {
				if v == field.ColimnName {
					switch field.ColumnType {
					case Int:
						value, err := strconv.ParseInt(r.Cells[i].Value, 10, 64)
						if err != nil {
							return nil, err
						}
						row[v] = value
					case Float:
						value, err := strconv.ParseFloat(r.Cells[i].Value, 64)
						if err != nil {
							return nil, err
						}
						row[v] = value
					case String:
						row[v] = r.Cells[i].Value
					default:
						return nil, fmt.Errorf("column type is unknown type")
					}
				}

			}

		}
		result = append(result, row)
	}
	// var people = []struct {
	// 	Name string `mapstructure:"name"`
	// 	Age  int64  `mapstructure:"age"`
	// }{{}}
	// err := mapstructure.Decode(result, &people)
	// if err != nil {
	// 	fmt.Println("decode error:----->", err)
	// 	return err
	// }
	// fmt.Println("people:--->", people)

	return result, nil
}
