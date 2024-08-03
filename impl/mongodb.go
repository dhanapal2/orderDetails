package impl

import (
	"context"
	"log"
	"orderDetails/model"
	"regexp"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB model.MongoDB

func DbConnect() error {
	var err error
	clientOptions := options.Client().ApplyURI(MongoURL)
	MongoDB.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Please check the Mongo url ")
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = MongoDB.Client.Ping(ctx, nil)
	if err != nil {
		log.Println("Mongodb is not connecting. Please check the connection ...")
		return err
	}

	log.Println("Connected to MongoDB!")

	MongoDB.Database = MongoDB.Client.Database("Lumel")
	MongoDB.Collection = MongoDB.Database.Collection("orders")
	return nil
}

func InsertMany(docs []interface{}, wg *sync.WaitGroup, response *model.Result) {
	defer wg.Done()
	re := regexp.MustCompile(Regex)
	opts := options.InsertMany().SetOrdered(false)
	result, err := MongoDB.Collection.InsertMany(context.TODO(), docs, opts)
	if err != nil {
		if writeException, ok := err.(mongo.BulkWriteException); ok {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == 11000 {
					matches := re.FindStringSubmatch(writeError.Message) //append the duplicate key
					response.DuplicatedKey = append(response.DuplicatedKey, matches[1])
				}
			}
		}
	}
	if len(result.InsertedIDs) > 0 {
		response.InsertedKey = append(response.InsertedKey, result.InsertedIDs...)
	}
}

func ExecuteQuery(query mongo.Pipeline) ([]bson.M, error) {
	cursor, err := MongoDB.Collection.Aggregate(context.TODO(), query)
	if err != nil {
		log.Println("Error in exeting query")
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("Error in formting query Result ::: " + err.Error())
		return nil, err
	}

	return results, nil

}
