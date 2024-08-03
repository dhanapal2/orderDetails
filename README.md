# OrderDetails
* Language: Golang
* Verion: 1.22
* DataBase: MongoDB
* MockData: order_details.xlsx

# Install Packages
    go get go.mongodb.org/mongo-driver
    go get github.com/xuri/excelize/v2
    go get github.com/sirupsen/logrus

# Run Code
    go mod tidy
    go run main.go
* Server Port: 8080

# API Route

| Method | Path         | Description                                    | Request Body                                             | Response                                                                                              |
|--------|--------------|------------------------------------------------|----------------------------------------------------------|-------------------------------------------------------------------------------------------------------|
| GET    | /healthCheck | To check if the server is running              | -                                                        | OK                                                                                                    |
| GET    | /reloadAll   | Read the CSV and load it to the database       | -                                                        | Order data is being read. It will take some time to load into the database. Check logs for issues.  |
| POST   | /revenue     | Calculate revenue based on provided filters    | `{ "action": "revenue", "filter": "category", "startAt": "YYYY-MM-DD", "endAt": "YYYY-MM-DD" }` | JSON object containing the revenue data or an error message. For example: `{ "success": true, filter: "All/category/product", "response":[{ "category": "Shoes","revenue": 1080},{"category": "Clothing","revenue": 179.97}}]` |

# Scheduler
Note: Scheduler will run for every one hour to read the data in csv and store it in mongo. If dupicate found in db it won't upsert and also we can see list of orers stores in mongo and list of items found duplicates in executor-logs file.


