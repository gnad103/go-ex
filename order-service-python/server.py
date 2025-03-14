import time
import grpc
import logging
from concurrent import futures
from datetime import datetime

# Import generated protobuf code
import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from proto import order_pb2, order_pb2_grpc
from proto import user_pb2, user_pb2_grpc
from proto import product_pb2, product_pb2_grpc

from service import OrderServicer

def serve():
    # Create gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    # Create stubs for user and product services
    user_channel = grpc.insecure_channel('localhost:50051')
    user_stub = user_pb2_grpc.UserServiceStub(user_channel)
    
    product_channel = grpc.insecure_channel('localhost:50052')
    product_stub = product_pb2_grpc.ProductServiceStub(product_channel)
    
    # Add order servicer to server
    order_servicer = OrderServicer(user_stub, product_stub)
    order_pb2_grpc.add_OrderServiceServicer_to_server(order_servicer, server)
    
    # Start server
    server.add_insecure_port('0.0.0.0:50053')
    server.start()
    
    logging.info("Order service (Python) is running on port 50053...")
    
    try:
        while True:
            time.sleep(86400)  # One day in seconds
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("Server stopped")

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()