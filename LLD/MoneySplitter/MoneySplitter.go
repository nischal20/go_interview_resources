package main

import (
	"fmt"
	"sync"
)

// User struct
type User struct {
	ID         string
	Name       string
	Email      string
	BankDetails string
}

// UserFactory using Factory pattern
type UserFactory struct {
	users map[string]*User
	mu    sync.Mutex
}

func NewUserFactory() *UserFactory {
	return &UserFactory{
		users: make(map[string]*User),
	}
}

func (f *UserFactory) CreateUser(id, name, email, bankDetails string) *User {
	f.mu.Lock()
	defer f.mu.Unlock()
	if user, exists := f.users[id]; exists {
		return user
	}
	user := &User{ID: id, Name: name, Email: email, BankDetails: bankDetails}
	f.users[id] = user
	return user
}

// ContributionStrategy using Strategy pattern
type ContributionStrategy interface {
	CalculateContributions(members map[string]float64, totalAmount float64) map[string]float64
}

type EqualContributionStrategy struct{}

func (e *EqualContributionStrategy) CalculateContributions(members map[string]float64, totalAmount float64) map[string]float64 {
	perMember := totalAmount / float64(len(members))
	contributions := make(map[string]float64)
	for member := range members {
		contributions[member] = perMember
	}
	return contributions
}

type CustomContributionStrategy struct{}

func (c *CustomContributionStrategy) CalculateContributions(members map[string]float64, totalAmount float64) map[string]float64 {
	return members
}

// ExpenseState using State pattern
type ExpenseState interface {
	AddMember(expense *Expense, userID string, amount float64)
	MoveToPending(expense *Expense)
	AddContribution(expense *Expense, userID string, amount float64)
}

type CreatedState struct{}

func (s *CreatedState) AddMember(expense *Expense, userID string, amount float64) {
	expense.Members[userID] = amount
}

func (s *CreatedState) MoveToPending(expense *Expense) {
	expense.SetState(&PendingState{})
}

func (s *CreatedState) AddContribution(expense *Expense, userID string, amount float64) {
	fmt.Println("Cannot add contribution in Created state")
}

type PendingState struct{}

func (s *PendingState) AddMember(expense *Expense, userID string, amount float64) {
	fmt.Println("Cannot add members in Pending state")
}

func (s *PendingState) MoveToPending(expense *Expense) {
	fmt.Println("Expense is already in Pending state")
}

func (s *PendingState) AddContribution(expense *Expense, userID string, amount float64) {
	if _, exists := expense.Paid[userID]; !exists {
		expense.Paid[userID] = 0
	}
	expense.Paid[userID] += amount
	if expense.isSettled() {
		expense.SetState(&SettledState{})
	}
}

type SettledState struct{}

func (s *SettledState) AddMember(expense *Expense, userID string, amount float64) {
	fmt.Println("Cannot add members in Settled state")
}

func (s *SettledState) MoveToPending(expense *Expense) {
	fmt.Println("Cannot move to Pending state from Settled state")
}

func (s *SettledState) AddContribution(expense *Expense, userID string, amount float64) {
	fmt.Println("Expense is already settled")
}

// Expense struct
type Expense struct {
	ID         string
	Title      string
	Amount     float64
	State      ExpenseState
	Members    map[string]float64
	Paid       map[string]float64
	mu         sync.Mutex
}

func NewExpense(id, title string, amount float64, strategy ContributionStrategy) *Expense {
	return &Expense{
		ID:      id,
		Title:   title,
		Amount:  amount,
		State:   &CreatedState{},
		Members: strategy.CalculateContributions(make(map[string]float64), amount),
		Paid:    make(map[string]float64),
	}
}

func (e *Expense) SetState(state ExpenseState) {
	e.State = state
}

func (e *Expense) AddMember(userID string, amount float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.State.AddMember(e, userID, amount)
}

func (e *Expense) MoveToPending() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.State.MoveToPending(e)
}

func (e *Expense) AddContribution(userID string, amount float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.State.AddContribution(e, userID, amount)
}

func (e *Expense) isSettled() bool {
	totalPaid := 0.0
	for _, amount := range e.Paid {
		totalPaid += amount
	}
	return totalPaid >= e.Amount
}

// Observer pattern for notifications
type Observer interface {
	Update(expenseID string, message string)
}

type Notifier struct {
	observers map[string]Observer
}

func NewNotifier() *Notifier {
	return &Notifier{
		observers: make(map[string]Observer),
	}
}

func (n *Notifier) Register(id string, observer Observer) {
	n.observers[id] = observer
}

func (n *Notifier) Deregister(id string) {
	delete(n.observers, id)
}

func (n *Notifier) Notify(expenseID string, message string) {
	for _, observer := range n.observers {
		observer.Update(expenseID, message)
	}
}

type UserNotifier struct {
	userID string
}

func (u *UserNotifier) Update(expenseID string, message string) {
	fmt.Printf("User %s notified for expense %s: %s\n", u.userID, expenseID, message)
}

// Main function to demonstrate the application
func main() {
	userFactory := NewUserFactory()
	notifier := NewNotifier()

	user1 := userFactory.CreateUser("1", "Alice", "alice@example.com", "Bank1")
	user2 := userFactory.CreateUser("2", "Bob", "bob@example.com", "Bank2")

	notifier.Register(user1.ID, &UserNotifier{userID: user1.ID})
	notifier.Register(user2.ID, &UserNotifier{userID: user2.ID})

	equalStrategy := &EqualContributionStrategy{}
	expense := NewExpense("e1", "Dinner", 100.0, equalStrategy)

	expense.AddMember(user1.ID, 50.0)
	expense.AddMember(user2.ID, 50.0)
	expense.MoveToPending()

	notifier.Notify(expense.ID, "Expense moved to Pending state")

	expense.AddContribution(user1.ID, 50.0)
	expense.AddContribution(user2.ID, 50.0)

	notifier.Notify(expense.ID, "Expense settled")

	fmt.Printf("Expense: %+v\n", expense)
}
