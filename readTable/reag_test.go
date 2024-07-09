package readtable

import (
	"fmt"
	"testing"

	messages "github.com/cucumber/messages/go/v21"
)

func TestTableToStruct(t *testing.T) {
	tbl := messages.PickleTable{
		Rows: []*messages.PickleTableRow{
			{
				Cells: []*messages.PickleTableCell{
					{
						Value: "name",
					},
					{
						Value: "age",
					},
				},
			},
			{
				Cells: []*messages.PickleTableCell{
					{
						Value: "jack",
					},
					{
						Value: "45",
					},
				},
			},
			{
				Cells: []*messages.PickleTableCell{
					{
						Value: "john",
					},
					{
						Value: "40",
					},
				},
			},
		},
	}
	fmt.Println(ReadTableData(&tbl, []Column{
		{
			ColimnName: "name",
			ColumnType: String,
		},
		{
			ColimnName: "age",
			ColumnType: Int,
		},
	}))
}
