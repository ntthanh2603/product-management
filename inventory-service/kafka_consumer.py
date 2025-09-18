import json
import logging
import threading
from kafka import KafkaConsumer
from kafka.errors import KafkaError
import os
from dotenv import load_dotenv

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class InventoryKafkaConsumer:
    def __init__(self, inventory_service=None):
        self.bootstrap_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
        self.inventory_service = inventory_service
        self.consumer = None
        self.running = False
        self.consumer_thread = None
    
    def connect(self):
        try:
            self.consumer = KafkaConsumer(
                'order-events',
                'inventory-events',
                bootstrap_servers=self.bootstrap_servers,
                value_deserializer=lambda x: json.loads(x.decode('utf-8')),
                key_deserializer=lambda x: x.decode('utf-8') if x else None,
                group_id='inventory-service-group',
                auto_offset_reset='latest'
            )
            logger.info("Connected to Kafka consumer")
            return True
        except Exception as e:
            logger.error(f"Failed to connect to Kafka consumer: {e}")
            return False
    
    def start_consuming(self):
        if not self.connect():
            return False
        
        self.running = True
        self.consumer_thread = threading.Thread(target=self._consume_messages)
        self.consumer_thread.daemon = True
        self.consumer_thread.start()
        logger.info("Started Kafka consumer thread")
        return True
    
    def _consume_messages(self):
        try:
            for message in self.consumer:
                if not self.running:
                    break
                
                try:
                    self._process_message(message)
                except Exception as e:
                    logger.error(f"Error processing message: {e}")
                    
        except Exception as e:
            logger.error(f"Error in consumer loop: {e}")
    
    def _process_message(self, message):
        topic = message.topic
        event_data = message.value
        event_type = event_data.get('event_type')
        data = event_data.get('data', {})
        
        logger.info(f"Processing {event_type} from topic {topic}")
        
        if topic == 'order-events':
            self._handle_order_event(event_type, data)
        elif topic == 'inventory-events':
            self._handle_inventory_event(event_type, data)
    
    def _handle_order_event(self, event_type, order_data):
        if not self.inventory_service:
            return
        
        try:
            if event_type == 'ORDER_CREATED':
                # Reserve stock for new order
                logger.info(f"Handling order creation: {order_data.get('id')}")
                # Implementation would be handled by the inventory service
                
            elif event_type == 'ORDER_CANCELLED':
                # Release reserved stock
                logger.info(f"Handling order cancellation: {order_data.get('id')}")
                # Implementation would be handled by the inventory service
                
            elif event_type == 'ORDER_CONFIRMED':
                # Confirm stock reservation
                logger.info(f"Handling order confirmation: {order_data.get('id')}")
                # Implementation would be handled by the inventory service
                
        except Exception as e:
            logger.error(f"Error handling order event {event_type}: {e}")
    
    def _handle_inventory_event(self, event_type, inventory_data):
        try:
            if event_type == 'STOCK_UPDATED':
                logger.info(f"Stock updated for product {inventory_data.get('product_id')}: {inventory_data.get('quantity')}")
                
            elif event_type == 'STOCK_RESERVED':
                logger.info(f"Stock reserved for product {inventory_data.get('product_id')}: {inventory_data.get('reserved_quantity')}")
                
            elif event_type == 'STOCK_RELEASED':
                logger.info(f"Stock released for product {inventory_data.get('product_id')}")
                
        except Exception as e:
            logger.error(f"Error handling inventory event {event_type}: {e}")
    
    def stop(self):
        self.running = False
        if self.consumer:
            self.consumer.close()
        if self.consumer_thread:
            self.consumer_thread.join(timeout=5)
        logger.info("Kafka consumer stopped")