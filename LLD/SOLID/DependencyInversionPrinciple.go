package main

import "fmt"

type Database interface {
    GetData() string
}

type MySQLDatabase struct{}

func (db MySQLDatabase) GetData() string {
    return "Data from MySQL"
}

type Application struct {
    db Database
}

func (app Application) FetchData() {
    fmt.Println(app.db.GetData())
}

func main() {
    mysqlDB := MySQLDatabase{}
    app := Application{db: mysqlDB}
    app.FetchData()
}
