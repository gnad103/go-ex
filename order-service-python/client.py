import grpc
import logging
import sys
import os

# Import generated protobuf code
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from proto import order_pb2, order_pb2_grpc
from proto import user_pb2, user_pb2_grpc
from proto import product_pb2, product_pb2_grpc

def run():
    # Connect to services
    with grpc.insecure_channel('localhost:50051') as user_channel, \
         grpc.insecure_channel('localhost:50052') as product_channel, \
         grpc.insecure_channel('localhost:50053') as order_channel:
        
        # Create stubs
        user_stub = user_pb2_grpc.UserServiceStub(user_channel)
        product_stub = product_pb2_grpc.ProductServiceStub(product_channel)
        order_stub = order_pb2_grpc.OrderServiceStub(order_channel)
        
        # Create a user
        user = user_stub.CreateUser(user_pb2.CreateUserRequest(
            name="Alice Smith",
            email="alice@example.com"
        ))
        logging.info(f"Created user: {user.id}, {user.name}, {user.email}")
        
        # Create products
        product1 = product_stub.CreateProduct(product_pb2.CreateProductRequest(
            name="Smartphone",
            description="Latest smartphone model",
            price=799.99,
            user_id=user.id
        ))
        logging.info(f"Created product: {product1.id}, {product1.name}, ${product1.price}")
        
        product2 = product_stub.CreateProduct(product_pb2.CreateProductRequest(
            name="Headphones",
            description="Wireless noise-cancelling headphones",
            price=199.99,
            user_id=user.id
        ))
        logging.info(f"Created product: {product2.id}, {product2.name}, ${product2.price}")
        
        # Create an order
        order_items = [
            order_pb2.OrderItem(product_id=product1.id, quantity=1),
            order_pb2.OrderItem(product_id=product2.id, quantity=2)
        ]
        
        order = order_stub.CreateOrder(order_pb2.CreateOrderRequest(
            user_id=user.id,
            items=order_items
        ))
        
        logging.info(f"Created order: {order.id}")
        logging.info(f"Order total: ${order.total_amount}")
        logging.info(f"Order items:")
        for item in order.items:
            logging.info(f"  - {item.quantity}x {item.product_name}: ${item.subtotal}")
        
        # Get user orders
        user_orders = order_stub.GetUserOrders(order_pb2.UserOrderRequest(user_id=user.id))
        logging.info(f"User {user.id} has {len(user_orders.orders)} orders")

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    run()