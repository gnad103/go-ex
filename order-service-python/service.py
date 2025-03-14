import logging
import grpc
from datetime import datetime
from google.protobuf.timestamp_pb2 import Timestamp

# Import generated protobuf code
import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from proto import order_pb2, order_pb2_grpc
from proto import user_pb2, user_pb2_grpc
from proto import product_pb2, product_pb2_grpc

class OrderServicer(order_pb2_grpc.OrderServiceServicer):
    def __init__(self, user_stub, product_stub):
        self.user_stub = user_stub
        self.product_stub = product_stub
        self.orders = {}
        self.next_id = 1
    
    def CreateOrder(self, request, context):
        # Verify user exists
        try:
            user = self.user_stub.GetUser(user_pb2.UserRequest(id=request.user_id))
        except grpc.RpcError as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"User with ID {request.user_id} not found: {e.details()}")
            return order_pb2.OrderResponse()
        
        # Verify all products exist and calculate total
        items = []
        total_amount = 0.0
        
        for item in request.items:
            try:
                product = self.product_stub.GetProduct(product_pb2.ProductRequest(id=item.product_id))
                
                subtotal = product.price * item.quantity
                total_amount += subtotal
                
                item_detail = order_pb2.OrderItemDetail(
                    product_id=product.id,
                    product_name=product.name,
                    product_price=product.price,
                    quantity=item.quantity,
                    subtotal=subtotal
                )
                items.append(item_detail)
                
            except grpc.RpcError as e:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details(f"Product with ID {item.product_id} not found: {e.details()}")
                return order_pb2.OrderResponse()
        
        # Create order
        order_id = self.next_id
        self.next_id += 1
        
        created_at = datetime.now().isoformat()
        
        order = order_pb2.OrderResponse(
            id=order_id,
            user_id=request.user_id,
            items=items,
            total_amount=total_amount,
            status="CREATED",
            created_at=created_at
        )
        
        self.orders[order_id] = order
        
        logging.info(f"Created order {order_id} for user {request.user_id} with {len(items)} items")
        return order
    
    def GetOrder(self, request, context):
        order_id = request.id
        if order_id not in self.orders:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Order with ID {order_id} not found")
            return order_pb2.OrderResponse()
        
        return self.orders[order_id]
    
    def GetUserOrders(self, request, context):
        user_id = request.user_id
        
        # Verify user exists
        try:
            user = self.user_stub.GetUser(user_pb2.UserRequest(id=user_id))
        except grpc.RpcError as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"User with ID {user_id} not found: {e.details()}")
            return order_pb2.OrderListResponse()
        
        # Get all orders for this user
        user_orders = [order for order in self.orders.values() if order.user_id == user_id]
        
        return order_pb2.OrderListResponse(orders=user_orders)