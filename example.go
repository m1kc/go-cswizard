package cswizard

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/m1kc/go-cswizard"
)

type Client struct {
	Name   string
	Age    uint64
	Height uint64
}

func main() {
	clients := make([]Client, 0, 5)
	clients = append(clients, Client{Name: "John Smith", Age: 34, Height: 181})
	clients = append(clients, Client{Name: "Peter Smith", Age: 44, Height: 179})
	clients = append(clients, Client{Name: "John Sebastine", Age: 33, Height: 159})
	clients = append(clients, Client{Name: "Markus Hallberg", Age: 18, Height: 169})
	clients = append(clients, Client{Name: "John Wonapaska", Age: 59, Height: 170})
	
	cw := csv.NewWriter(os.Stdout)
	w := cswizard.New(cw)
	
	// TODO:
	// 1. Add an empty column between Name and Age;
	// 2. Swap Age and Height columns;
	// 3. Remove the Age column completely.
	//
	// Try to play around with these tasks, and you will get a rough idea
	// why CSWizard was worth creating.
	
	colName := w.AddHeader("Client name")
	colAge := w.AddHeader("Client age")
	colHeight := w.AddHeader("Client height (predicted)")
	
	w.LockHeaders()
	
	for _, c := range clients {
		row := w.CreateRow()
		row[colName] = c.Name
		row[colAge] = strconv.FormatUint(c.Age, 10)
		row[colHeight] = strconv.FormatUint(c.Height, 10)
		
		err := w.CommitRow(row)
		if err != nil {
			return
		}
	}
	
	cw.Flush()

	return
}
