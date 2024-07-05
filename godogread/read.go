package godogread

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

func ReadData(arg1 *godog.Table) error {
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
			row[v] = r.Cells[i].Value
		}
		result = append(result, row)
	}

	var createSession = struct {
		ID           string    `json:"id"`
		TicketNumber int       `json:"ticket_number"`
		BusNumber    int       `json:"bus_number"`
		CreatedAt    time.Time `json:"created_at"`
	}{}

	byteData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("marshal:", err)
		return err
	}
	err = json.Unmarshal(byteData, &createSession)
	if err != nil {
		fmt.Println("unmarshal:", err)
		return err
	}
	return nil
}
