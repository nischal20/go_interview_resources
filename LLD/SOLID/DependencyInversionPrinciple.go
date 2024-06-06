// Dependency Inversion Principle (DIP)
// High-level modules should not depend on low-level modules. Both should depend on abstractions (e.g., interfaces).
// Abstractions should not depend on details. Details (concrete implementations) should depend on abstractions.


// The below code follows the Dependency Inversion Principle by:
// Creating an abstraction (Database interface).
// Having high-level modules (Application) depend on this abstraction.
// Implementing low-level modules (MySQLDatabase) that conform to this abstraction.

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
