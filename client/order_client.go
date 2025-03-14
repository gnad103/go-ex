package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/gnad103/go-ex/proto"
)

func main() {
	// Connect to user service (Go)
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	userClient := proto.NewUserServiceClient(userConn)

	// Connect to product service (Go)
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	productClient := proto.NewProductServiceClient(productConn)

	// Connect to order service (Python)
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	orderClient := proto.NewOrderServiceClient(orderConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create a user
	user, err := userClient.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "Bob Johnson",
		Email: "bob@example.com",
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	log.Printf("Created user: %v", user)

	// Create products
	product1, err := productClient.CreateProduct(ctx, &proto.CreateProductRequest{
		Name:        "Tablet",
		Description: "10-inch tablet",
		Price:       349.99,
		UserId:      user.Id,
	})
	if err != nil {
		log.Fatalf("Could not create product: %v", err)
	}
	log.Printf("Created product: %v", product1)

	product2, err := productClient.CreateProduct(ctx, &proto.CreateProductRequest{
		Name:        "Keyboard",
		Description: "Mechanical keyboard",
		Price:       129.99,
		UserId:      user.Id,
	})
	if err != nil {
		log.Fatalf("Could not create product: %v", err)
	}
	log.Printf("Created product: %v", product2)

	// Create an order using the Python service
	orderItems := []*proto.OrderItem{
		{
			ProductId: product1.Id,
			Quantity:  1,
		},
		{
			ProductId: product2.Id,
			Quantity:  1,
		},
	}

	order, err := orderClient.CreateOrder(ctx, &proto.CreateOrderRequest{
		UserId: user.Id,
		Items:  orderItems,
	})
	if err != nil {
		log.Fatalf("Could not create order: %v", err)
	}

	log.Printf("Created order: ID=%d, Total=$%.2f", order.Id, order.TotalAmount)
	log.Printf("Order items:")
	for _, item := range order.Items {
		log.Printf("  - %dx %s: $%.2f", item.Quantity, item.ProductName, item.Subtotal)
	}

	// Get user orders
	userOrders, err := orderClient.GetUserOrders(ctx, &proto.UserOrderRequest{
		UserId: user.Id,
	})
	if err != nil {
		log.Fatalf("Could not get user orders: %v", err)
	}

	log.Printf("User %d has %d orders", user.Id, len(userOrders.Orders))
}
