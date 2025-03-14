package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gnad103/go-ex/proto"
)

// APIHandler handles the REST API requests and communicates with the gRPC services
type APIHandler struct {
	userClient    pb.UserServiceClient
	productClient pb.ProductServiceClient
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(userClient pb.UserServiceClient, productClient pb.ProductServiceClient) *APIHandler {
	return &APIHandler{
		userClient:    userClient,
		productClient: productClient,
	}
}

// ========== User Handlers ==========

// GetUser handles GET /api/users/{id}
func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	user, err := h.userClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			http.Error(w, st.Message(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
}

// CreateUser handles POST /api/users
func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int32  `json:"age"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := h.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  requestBody.Name,
		Email: requestBody.Email,
		Age:   requestBody.Age,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
}

// ListUsers handles GET /api/users
func (h *APIHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Extract pagination parameters if provided
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := int32(1)      // Default page
	pageSize := int32(10) // Default page size

	if pageStr != "" {
		if pageInt, err := strconv.Atoi(pageStr); err == nil && pageInt > 0 {
			page = int32(pageInt)
		}
	}

	if pageSizeStr != "" {
		if pageSizeInt, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeInt > 0 {
			pageSize = int32(pageSizeInt)
		}
	}

	ctx := context.Background()
	response, err := h.userClient.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	// Convert users to a format suitable for JSON response
	users := make([]map[string]interface{}, 0, len(response.Users))
	for _, user := range response.Users {
		users = append(users, map[string]interface{}{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}

// ========== Product Handlers ==========

// GetProduct handles GET /api/products/{id}
func (h *APIHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	product, err := h.productClient.GetProduct(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			http.Error(w, st.Message(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          product.Id,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
	})
}

// CreateProduct handles POST /api/products
func (h *APIHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	product, err := h.productClient.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		Price:       requestBody.Price,
	})
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          product.Id,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
	})
}

// ListProducts handles GET /api/products
func (h *APIHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Extract pagination parameters if provided
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := int32(1)      // Default page
	pageSize := int32(10) // Default page size

	if pageStr != "" {
		if pageInt, err := strconv.Atoi(pageStr); err == nil && pageInt > 0 {
			page = int32(pageInt)
		}
	}

	if pageSizeStr != "" {
		if pageSizeInt, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeInt > 0 {
			pageSize = int32(pageSizeInt)
		}
	}

	ctx := context.Background()
	response, err := h.productClient.ListProducts(ctx, &pb.ListProductsRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	// Convert products to a format suitable for JSON response
	products := make([]map[string]interface{}, 0, len(response.Products))
	for _, product := range response.Products {
		products = append(products, map[string]interface{}{
			"id":          product.Id,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
	})
}
