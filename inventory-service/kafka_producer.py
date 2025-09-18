import json
import logging
from kafka import KafkaProducer
from kafka.errors import KafkaError
import os
from dotenv import load_dotenv

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class InventoryKafkaProducer:
    def __init__(self):
        self.bootstrap_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
        self.producer = None
        self.connect()
    
    def connect(self):
        try:
            self.producer = KafkaProducer(
                bootstrap_servers=self.bootstrap_servers,
                value_serializer=lambda x: json.dumps(x).encode('utf-8'),
                key_serializer=lambda x: x.encode('utf-8') if x else None
            )
            logger.info("Connected to Kafka producer")
        except Exception as e:
            logger.error(f"Failed to connect to Kafka: {e}")
            self.producer = None
    
    def send_order_event(self, event_type, order_data):
        if not self.producer:
            logger.warning("Kafka producer not available")
            return False
        
        try:
            topic = "order-events"
            message = {
                "event_type": event_type,
                "timestamp": order_data.get("created_at", ""),
                "data": order_data
            }
            
            future = self.producer.send(topic, key=order_data.get("id"), value=message)
            self.producer.flush()
            logger.info(f"Sent {event_type} event for order {order_data.get('id')}")
            return True
        except KafkaError as e:
            logger.error(f"Failed to send order event: {e}")
            return False
    
    def send_inventory_event(self, event_type, inventory_data):
        if not self.producer:
            logger.warning("Kafka producer not available")
            return False
        
        try:
            topic = "inventory-events"
            message = {
                "event_type": event_type,
                "timestamp": inventory_data.get("updated_at", ""),
                "data": inventory_data
            }
            
            future = self.producer.send(topic, key=str(inventory_data.get("product_id")), value=message)
            self.producer.flush()
            logger.info(f"Sent {event_type} event for product {inventory_data.get('product_id')}")
            return True
        except KafkaError as e:
            logger.error(f"Failed to send inventory event: {e}")
            return False
    
    def close(self):
        if self.producer:
            self.producer.close()
            logger.info("Kafka producer closed")