package impl

import (
	"log"
	"orderDetails/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetRevenue(reqBody model.Request) (interface{}, error) {
	var query mongo.Pipeline
	validFields := []string{"category", "productID", "region"}
	if reqBody.Filter == "All" {
		query = mongo.Pipeline{
			{{"$match", bson.D{
				{"saleDate", bson.D{
					{"$gte", reqBody.StartAt},
					{"$lte", reqBody.EndAt},
				}},
			}}},
			{{"$addFields", bson.D{
				{"price", bson.D{{"$toDouble", "$price"}}},
			}}},
			{{"$group", bson.D{
				{"_id", nil},
				{"totalPrice", bson.D{
					{"$sum", "$price"},
				}},
			}}},
			{{"$project", bson.D{
				{"_id", 0},
				{"totalPrice", 1},
			}}},
		}
	} else {
		isvalidFiled := false
		for _, filter := range validFields {
			if filter == reqBody.Filter {
				isvalidFiled = true
				break
			}
		}
		if !isvalidFiled {
			return map[string]interface{}{"success": false, "error": "Invalid Filter field"}, nil
		}
		query = mongo.Pipeline{
			{{"$match", bson.D{
				{"saleDate", bson.D{
					{"$gte", reqBody.StartAt},
					{"$lte", reqBody.EndAt},
				}},
			}}},
			{{"$addFields", bson.D{
				{"price", bson.D{{"$toDouble", "$price"}}},
			}}},
			{{"$group", bson.D{
				{"_id", "$" + reqBody.Filter},
				{"totalPrice", bson.D{
					{"$sum", "$price"},
				}},
			}}},
			{{"$project", bson.D{
				{"_id", 0},
				{reqBody.Filter, "$_id"},
				{"totalPrice", 1},
			}}},
		}

	}
	resp, err := ExecuteQuery(query)
	if err != nil {
		log.Println("Error in Executing query ::: ", err)
		return map[string]interface{}{"success": false, "error": err}, err
	}
	response := make(map[string]interface{})
	if reqBody.Filter == "All" {
		for _, doc := range resp {
			if val, ok := doc["totalPrice"].(float64); ok {
				response = map[string]interface{}{
					"success": true,
					"filter":  reqBody.Filter,
					"response": map[string]interface{}{
						"revenue": val,
					},
				}
			}
		}
	} else {
		respList := make([]map[string]interface{}, 0)
		for _, doc := range resp {
			respList = append(respList, map[string]interface{}{
				reqBody.Filter: doc[reqBody.Filter],
				"revenue":      doc["totalPrice"],
			})
		}
		response = map[string]interface{}{
			"success":      true,
			"filter":       reqBody.Filter,
			"responselist": respList,
		}
	}
	return response, err
}
