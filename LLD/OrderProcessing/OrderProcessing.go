package main

import (
	"fmt"
	"sync"
	"time"
)

// Models

type Order struct {
	OrderID     int       `json:"orderId"`
	UserID      int       `json:"userId"`
	ProductList []string  `json:"productList"`
	OrderStatus string    `json:"orderStatus"`
	TotalPrice  float64   `json:"totalPrice"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type OrderRequest struct {
	UserID      int      `json:"userId"`
	ProductList []string `json:"productList"`
	TotalPrice  float64  `json:"totalPrice"`
}

type PaymentRequest struct {
	OrderID int     `json:"orderId"`
	Amount  float64 `json:"amount"`
}

type InventoryRequest struct {
	ProductList []string `json:"productList"`
	Quantity    int      `json:"quantity"`
}

type NotificationRequest struct {
	UserID  int    `json:"userId"`
	Message string `json:"message"`
}

// Interfaces

type IOrderRepository interface {
	Save(order Order)
	FindByID(orderID int) (Order, bool)
	Update(order Order)
	Delete(orderID int)
	Count() int
}

type IOrderService interface {
	CreateOrder(request OrderRequest) Order
	GetOrder(orderID int) (Order, bool)
	UpdateOrder(orderID int, status string) (Order, bool)
	CancelOrder(orderID int) bool
}

type IPaymentService interface {
	ProcessPayment(request PaymentRequest) (bool, string)
}

type IInventoryService interface {
	CheckStock(request InventoryRequest) (bool, string)
	UpdateStock(request InventoryRequest) string
}

type INotificationService interface {
	SendNotification(request NotificationRequest) string
}

// Repositories

type OrderRepository struct {
	mu     sync.Mutex
	orders map[int]Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int]Order),
	}
}

func (repo *OrderRepository) Save(order Order) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.orders[order.OrderID] = order
}

func (repo *OrderRepository) FindByID(orderID int) (Order, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	order, exists := repo.orders[orderID]
	return order, exists
}

func (repo *OrderRepository) Update(order Order) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.orders[order.OrderID] = order
}

func (repo *OrderRepository) Delete(orderID int) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	delete(repo.orders, orderID)
}

func (repo *OrderRepository) Count() int {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return len(repo.orders)
}

// Services

type OrderService struct {
	repo            IOrderRepository
	paymentSvc      IPaymentService
	inventorySvc    IInventoryService
	notificationSvc INotificationService
}

func NewOrderService(repo IOrderRepository, paymentSvc IPaymentService, inventorySvc IInventoryService, notificationSvc INotificationService) *OrderService {
	return &OrderService{
		repo:            repo,
		paymentSvc:      paymentSvc,
		inventorySvc:    inventorySvc,
		notificationSvc: notificationSvc,
	}
}

func (service *OrderService) CreateOrder(request OrderRequest) Order {
	order := Order{
		OrderID:     service.repo.Count() + 1,
		UserID:      request.UserID,
		ProductList: request.ProductList,
		OrderStatus: "Created",
		TotalPrice:  request.TotalPrice,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Check inventory before creating order
	inventoryRequest := InventoryRequest{
		ProductList: request.ProductList,
		Quantity:    1, // Assuming 1 unit per product for simplicity
	}
	inStock, message := service.inventorySvc.CheckStock(inventoryRequest)
	if !inStock {
		fmt.Println("Inventory Check Failed:", message)
		return order
	}

	service.repo.Save(order)

	// Notify user about order creation
	service.notificationSvc.SendNotification(NotificationRequest{
		UserID:  request.UserID,
		Message: fmt.Sprintf("Your order %d has been created.", order.OrderID),
	})

	return order
}

func (service *OrderService) GetOrder(orderID int) (Order, bool) {
	return service.repo.FindByID(orderID)
}

func (service *OrderService) UpdateOrder(orderID int, status string) (Order, bool) {
	order, exists := service.repo.FindByID(orderID)
	if !exists {
		return Order{}, false
	}
	order.OrderStatus = status
	order.UpdatedAt = time.Now()
	service.repo.Update(order)

	// Notify user about order status update
	service.notificationSvc.SendNotification(NotificationRequest{
		UserID:  order.UserID,
		Message: fmt.Sprintf("Your order %d status has been updated to %s.", order.OrderID, status),
	})

	return order, true
}

func (service *OrderService) CancelOrder(orderID int) bool {
	order, exists := service.repo.FindByID(orderID)
	if exists {
		service.repo.Delete(orderID)

		// Notify user about order cancellation
		service.notificationSvc.SendNotification(NotificationRequest{
			UserID:  order.UserID,
			Message: fmt.Sprintf("Your order %d has been cancelled.", orderID),
		})
	}
	return exists
}

// PaymentService

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) ProcessPayment(request PaymentRequest) (bool, string) {
	// Simulate payment processing logic
	fmt.Printf("Processing payment for Order ID: %d, Amount: %.2f\n", request.OrderID, request.Amount)
	return true, "Payment Successful"
}

// InventoryService

type InventoryService struct{}

func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

func (i *InventoryService) CheckStock(request InventoryRequest) (bool, string) {
	// Simulate inventory check logic
	fmt.Printf("Checking stock for Products: %v\n", request.ProductList)
	// Assume all products are in stock for simplicity
	return true, "Stock Available"
}

func (i *InventoryService) UpdateStock(request InventoryRequest) string {
	// Simulate stock update logic
	fmt.Printf("Updating stock for Products: %v, Quantity: %d\n", request.ProductList, request.Quantity)
	return "Stock Updated Successfully"
}

// NotificationService

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (n *NotificationService) SendNotification(request NotificationRequest) string {
	// Simulate sending notification logic
	fmt.Printf("Sending notification to User ID: %d, Message: %s\n", request.UserID, request.Message)
	return "Notification Sent"
}

// Factory for Dependency Injection

type ServiceFactory struct{}

func (f *ServiceFactory) CreateOrderService() IOrderService {
	repo := NewOrderRepository()
	paymentSvc := NewPaymentService()
	inventorySvc := NewInventoryService()
	notificationSvc := NewNotificationService()
	return NewOrderService(repo, paymentSvc, inventorySvc, notificationSvc)
}

// Main function to demonstrate usage

func main() {
	factory := &ServiceFactory{}
	orderService := factory.CreateOrderService()

	// Create an order
	orderRequest := OrderRequest{
		UserID:      1,
		ProductList: []string{"item1", "item2"},
		TotalPrice:  100.0,
	}
	newOrder := orderService.CreateOrder(orderRequest)
	fmt.Printf("Created Order: %+v\n", newOrder)

	// Process Payment
	paymentRequest := PaymentRequest{
		OrderID: newOrder.OrderID,
		Amount:  newOrder.TotalPrice,
	}
	paymentSvc := NewPaymentService()
	paymentSuccess, message := paymentSvc.ProcessPayment(paymentRequest)
	fmt.Printf("Payment Status: %v, Message: %s\n", paymentSuccess, message)

	// Update Inventory
	inventoryRequest := InventoryRequest{
		ProductList: []string{"item1", "item2"},
		Quantity:    1,
	}
	inventorySvc := NewInventoryService()
	stockMessage := inventorySvc.UpdateStock(inventoryRequest)
	fmt.Println(stockMessage)

	// Get an order
	fetchedOrder, exists := orderService.GetOrder(newOrder.OrderID)
	if exists {
		fmt.Printf("Fetched Order: %+v\n", fetchedOrder)
	} else {
		fmt.Println("Order not found")
	}

	// Update an order
	updatedOrder, exists := orderService.UpdateOrder(newOrder.OrderID, "Shipped")
	if exists {
		fmt.Printf("Updated Order: %+v\n", updatedOrder)
	} else {
		fmt.Println("Order not found")
	}

	// Cancel an order
	if orderService.CancelOrder(newOrder.OrderID) {
		fmt.Println("Order cancelled successfully")
	} else {
		fmt.Println("Order not found")
	}
}
