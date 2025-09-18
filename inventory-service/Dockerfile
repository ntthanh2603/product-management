FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Install grpcio-tools for proto generation
RUN pip install grpcio-tools

# Copy proto files and generate Python files
COPY ../protos/inventory.proto ./protos/
RUN python -m grpc_tools.protoc --proto_path=./protos --python_out=. --grpc_python_out=. ./protos/inventory.proto

# Copy application code
COPY . .

# Create directory for SQLite database
RUN mkdir -p /app/data

# Expose gRPC port
EXPOSE 50053

# Run the service
CMD ["python", "inventory_service.py"]