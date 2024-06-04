package main

import "fmt"

// Logger struct responsible only for logging
type Logger struct{}

func (l Logger) Log(message string) {
    fmt.Println("Log:", message)
}

// UserManager struct responsible only for user management
type UserManager struct {
    logger Logger
}

func (um UserManager) CreateUser(name string) {
    // Create user logic
    um.logger.Log("User created: " + name)
}

func main() {
    logger := Logger{}
    userManager := UserManager{logger: logger}
    userManager.CreateUser("John Doe")
}
