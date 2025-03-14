import threading
from concurrent import futures
import grpc
import service_pb2
import service_pb2_grpc

class ProductServiceServicer(service_pb2_grpc.ProductServiceServicer):
    """Implementation of the ProductService service."""

    def __init__(self):
        # In-memory database for simplicity
        self.products = {}
        self.lock = threading.Lock()

    def GetProduct(self, request, context):
        """Get a product by ID."""
        with self.lock:
            if request.id not in self.products:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"Product with ID {request.id} not found")
                return service_pb2.Product()
            
            return self.products[request.id]

    def CreateProduct(self, request, context):
        """Create a new product."""
        with self.lock:
            # Simple ID generation
            product_id = f"product-{len(self.products) + 1}"
            
            # Create the product
            product = service_pb2.Product(
                id=product_id,
                name=request.name,
                description=request.description,
                price=request.price
            )
            
            # Store the product
            self.products[product_id] = product
            
            return product

    def ListProducts(self, request, context):
        """List all products."""
        with self.lock:
            response = service_pb2.ListProductsResponse()
            for product in self.products.values():
                response.products.append(product)
            
            return response

def serve():
    """Start the gRPC server."""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_pb2_grpc.add_ProductServiceServicer_to_server(
        ProductServiceServicer(), server
    )
    server.add_insecure_port('[::]:50052')
    server.start()
    
    print("Python Product Service is running on port 50052...")
    server.wait_for_termination()