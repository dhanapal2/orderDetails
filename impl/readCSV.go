package impl

import (
	"orderDetails/model"
	"sync"

	"github.com/xuri/excelize/v2"
)

func ReadCSV() error {
	var wg sync.WaitGroup
	var finalResult model.Result
	file, err := excelize.OpenFile("order_details.xlsx")
	if err != nil {
		model.Log.Error("Error in opening file " + err.Error())
		return err
	}
	defer file.Close()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		model.Log.Error("Sheet not found in excel")
		return err
	}
	model.Log.Error("Readed xlsx successfully ...")
	doc := make([]interface{}, 0)
	for _, row := range rows[1:] {
		if len(row) < 14 {
			model.Log.Error("Incomplete row")
		}
		doc = append(doc, map[string]interface{}{
			"_id":           row[0],
			"orderID":       row[0],
			"productID":     row[1],
			"customerID":    row[2],
			"product":       row[3],
			"category":      row[4],
			"region":        row[5],
			"saleDate":      row[6],
			"quantitySold":  row[7],
			"price":         row[8],
			"discount":      row[9],
			"shippingCost":  row[10],
			"paymentMethod": row[11],
			"customerName":  row[12],
			"customerEmail": row[13],
		})

		if len(doc) == BatchValue {
			wg.Add(1)
			go InsertMany(doc, &wg, &finalResult)
			doc = make([]interface{}, 0)
		}
	}
	if len(doc) > 0 {
		wg.Add(1)
		go InsertMany(doc, &wg, &finalResult)
	}
	wg.Wait()
	model.Log.Debug("Total Duplicate ID ::: ", len(finalResult.DuplicatedKey), "List of IDs ::: ", finalResult.DuplicatedKey)
	model.Log.Debug("Total Inserted ID ::: ", len(finalResult.InsertedKey), "List of IDs ::: ", finalResult.InsertedKey)
	return nil
}
