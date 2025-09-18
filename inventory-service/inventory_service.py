import grpc
from concurrent import futures
import logging
import uuid
from datetime import datetime, timedelta
from sqlalchemy.orm import Session
from sqlalchemy import and_

import inventory_pb2
import inventory_pb2_grpc
from models import InventoryItem, Order, OrderItem, StockReservation, get_db, SessionLocal
from kafka_producer import InventoryKafkaProducer
from kafka_consumer import InventoryKafkaConsumer

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class InventoryServiceImpl(inventory_pb2_grpc.InventoryServiceServicer):
    def __init__(self):
        self.kafka_producer = InventoryKafkaProducer()
        self.kafka_consumer = InventoryKafkaConsumer(self)
        self.kafka_consumer.start_consuming()
    
    def CreateInventoryItem(self, request, context):
        db = SessionLocal()
        try:
            # Check if inventory item already exists for this product
            existing_item = db.query(InventoryItem).filter(
                and_(InventoryItem.product_id == request.product_id,
                     InventoryItem.location == request.location)
            ).first()
            
            if existing_item:
                # Update existing item
                existing_item.quantity += request.quantity
                existing_item.updated_at = datetime.utcnow()
                db.commit()
                db.refresh(existing_item)
                item = existing_item
            else:
                # Create new item
                item = InventoryItem(
                    product_id=request.product_id,
                    quantity=request.quantity,
                    location=request.location
                )
                db.add(item)
                db.commit()
                db.refresh(item)
            
            # Send Kafka event
            self.kafka_producer.send_inventory_event("STOCK_UPDATED", {
                "product_id": item.product_id,
                "quantity": item.quantity,
                "location": item.location,
                "updated_at": item.updated_at.isoformat()
            })
            
            return inventory_pb2.InventoryItemResponse(
                item=inventory_pb2.InventoryItem(
                    id=item.id,
                    product_id=item.product_id,
                    quantity=item.quantity,
                    reserved_quantity=item.reserved_quantity,
                    location=item.location,
                    created_at=item.created_at.isoformat(),
                    updated_at=item.updated_at.isoformat()
                ),
                message="Inventory item created/updated successfully"
            )
        except Exception as e:
            logger.error(f"Error creating inventory item: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.InventoryItemResponse(message=f"Error: {e}")
        finally:
            db.close()
    
    def GetInventoryItem(self, request, context):
        db = SessionLocal()
        try:
            item = db.query(InventoryItem).filter(InventoryItem.id == request.id).first()
            if not item:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Inventory item not found")
                return inventory_pb2.InventoryItemResponse(message="Inventory item not found")
            
            return inventory_pb2.InventoryItemResponse(
                item=inventory_pb2.InventoryItem(
                    id=item.id,
                    product_id=item.product_id,
                    quantity=item.quantity,
                    reserved_quantity=item.reserved_quantity,
                    location=item.location,
                    created_at=item.created_at.isoformat(),
                    updated_at=item.updated_at.isoformat()
                ),
                message="Inventory item retrieved successfully"
            )
        except Exception as e:
            logger.error(f"Error getting inventory item: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.InventoryItemResponse(message=f"Error: {e}")
        finally:
            db.close()
    
    def CheckStock(self, request, context):
        db = SessionLocal()
        try:
            total_available = 0
            items = db.query(InventoryItem).filter(InventoryItem.product_id == request.product_id).all()
            
            for item in items:
                total_available += item.available_quantity()
            
            available = total_available >= request.required_quantity
            
            return inventory_pb2.CheckStockResponse(
                available=available,
                available_quantity=total_available,
                message=f"Available: {total_available}, Required: {request.required_quantity}"
            )
        except Exception as e:
            logger.error(f"Error checking stock: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.CheckStockResponse(message=f"Error: {e}")
        finally:
            db.close()
    
    def ReserveStock(self, request, context):
        db = SessionLocal()
        try:
            # Check if enough stock is available
            total_available = 0
            items = db.query(InventoryItem).filter(InventoryItem.product_id == request.product_id).all()
            
            for item in items:
                total_available += item.available_quantity()
            
            if total_available < request.quantity:
                return inventory_pb2.ReserveStockResponse(
                    success=False,
                    message=f"Insufficient stock. Available: {total_available}, Required: {request.quantity}"
                )
            
            # Reserve stock from available items
            remaining_to_reserve = request.quantity
            
            for item in items:
                if remaining_to_reserve <= 0:
                    break
                
                available = item.available_quantity()
                if available > 0:
                    reserve_amount = min(available, remaining_to_reserve)
                    item.reserved_quantity += reserve_amount
                    remaining_to_reserve -= reserve_amount
            
            # Create reservation record
            reservation_id = str(uuid.uuid4())
            reservation = StockReservation(
                id=reservation_id,
                product_id=request.product_id,
                quantity=request.quantity,
                order_id=request.order_id,
                expires_at=datetime.utcnow() + timedelta(minutes=30)  # 30 minute reservation
            )
            
            db.add(reservation)
            db.commit()
            
            # Send Kafka event
            self.kafka_producer.send_inventory_event("STOCK_RESERVED", {
                "product_id": request.product_id,
                "reserved_quantity": request.quantity,
                "order_id": request.order_id,
                "reservation_id": reservation_id,
                "updated_at": datetime.utcnow().isoformat()
            })
            
            return inventory_pb2.ReserveStockResponse(
                success=True,
                message="Stock reserved successfully",
                reservation_id=reservation_id
            )
        except Exception as e:
            logger.error(f"Error reserving stock: {e}")
            db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.ReserveStockResponse(success=False, message=f"Error: {e}")
        finally:
            db.close()
    
    def ReleaseStock(self, request, context):
        db = SessionLocal()
        try:
            reservation = db.query(StockReservation).filter(
                and_(StockReservation.id == request.reservation_id,
                     StockReservation.is_active == True)
            ).first()
            
            if not reservation:
                return inventory_pb2.ReleaseStockResponse(
                    success=False,
                    message="Reservation not found or already released"
                )
            
            # Release reserved stock
            items = db.query(InventoryItem).filter(InventoryItem.product_id == reservation.product_id).all()
            remaining_to_release = reservation.quantity
            
            for item in items:
                if remaining_to_release <= 0:
                    break
                
                release_amount = min(item.reserved_quantity, remaining_to_release)
                item.reserved_quantity -= release_amount
                remaining_to_release -= release_amount
            
            # Mark reservation as inactive
            reservation.is_active = False
            
            db.commit()
            
            # Send Kafka event
            self.kafka_producer.send_inventory_event("STOCK_RELEASED", {
                "product_id": reservation.product_id,
                "released_quantity": reservation.quantity,
                "order_id": reservation.order_id,
                "reservation_id": request.reservation_id,
                "updated_at": datetime.utcnow().isoformat()
            })
            
            return inventory_pb2.ReleaseStockResponse(
                success=True,
                message="Stock released successfully"
            )
        except Exception as e:
            logger.error(f"Error releasing stock: {e}")
            db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.ReleaseStockResponse(success=False, message=f"Error: {e}")
        finally:
            db.close()

class OrderServiceImpl(inventory_pb2_grpc.OrderServiceServicer):
    def __init__(self):
        self.kafka_producer = InventoryKafkaProducer()
    
    def CreateOrder(self, request, context):
        db = SessionLocal()
        try:
            # Generate order ID
            order_id = str(uuid.uuid4())
            
            # Calculate total amount
            total_amount = sum(item.price * item.quantity for item in request.items)
            
            # Create order
            order = Order(
                id=order_id,
                user_id=request.user_id,
                total_amount=total_amount,
                status="PENDING"
            )
            
            db.add(order)
            
            # Create order items
            order_items = []
            for item_req in request.items:
                order_item = OrderItem(
                    order_id=order_id,
                    product_id=item_req.product_id,
                    quantity=item_req.quantity,
                    price=item_req.price
                )
                db.add(order_item)
                order_items.append(order_item)
            
            db.commit()
            db.refresh(order)
            
            # Send Kafka event
            order_data = {
                "id": order.id,
                "user_id": order.user_id,
                "total_amount": order.total_amount,
                "status": order.status,
                "created_at": order.created_at.isoformat(),
                "items": [
                    {
                        "product_id": item.product_id,
                        "quantity": item.quantity,
                        "price": item.price
                    } for item in order_items
                ]
            }
            
            self.kafka_producer.send_order_event("ORDER_CREATED", order_data)
            
            # Convert order items for response
            response_items = [
                inventory_pb2.OrderItem(
                    product_id=item.product_id,
                    quantity=item.quantity,
                    price=item.price
                ) for item in order_items
            ]
            
            return inventory_pb2.OrderResponse(
                order=inventory_pb2.Order(
                    id=order.id,
                    user_id=order.user_id,
                    items=response_items,
                    total_amount=order.total_amount,
                    status=inventory_pb2.OrderStatus.PENDING,
                    created_at=order.created_at.isoformat(),
                    updated_at=order.updated_at.isoformat()
                ),
                message="Order created successfully"
            )
        except Exception as e:
            logger.error(f"Error creating order: {e}")
            db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return inventory_pb2.OrderResponse(message=f"Error: {e}")
        finally:
            db.close()

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    inventory_service = InventoryServiceImpl()
    order_service = OrderServiceImpl()
    
    inventory_pb2_grpc.add_InventoryServiceServicer_to_server(inventory_service, server)
    inventory_pb2_grpc.add_OrderServiceServicer_to_server(order_service, server)
    
    listen_addr = '[::]:50053'
    server.add_insecure_port(listen_addr)
    
    logger.info(f"Starting Inventory Service server on {listen_addr}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down Inventory Service")
        inventory_service.kafka_consumer.stop()
        inventory_service.kafka_producer.close()
        order_service.kafka_producer.close()
        server.stop(0)

if __name__ == '__main__':
    serve()