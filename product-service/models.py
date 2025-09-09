from sqlalchemy import create_engine, Column, Integer, String, Text, DECIMAL, DateTime, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from datetime import datetime, timezone
import os

# Database configuration
DATABASE_URL = os.getenv("DATABASE_URL", "sqlite:///./products.db")

engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

class Product(Base):
    __tablename__ = "products"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(255), nullable=False)
    description = Column(Text)
    price = Column(DECIMAL(10, 2), nullable=False)
    user_id = Column(Integer, nullable=False)  # Reference to user in user-service
    created_at = Column(DateTime, default=lambda: datetime.now(timezone.utc))

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "description": self.description,
            "price": float(self.price),
            "user_id": self.user_id,
            "created_at": self.created_at.isoformat() if self.created_at else None
        }

def get_db() -> Session:
    db = SessionLocal()
    try:
        return db
    finally:
        pass

# Create tables
Base.metadata.create_all(bind=engine)