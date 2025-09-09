import grpc
from concurrent import futures
import product_pb2
import product_pb2_grpc
from models import Product, SessionLocal
from datetime import datetime
import logging
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class ProductService(product_pb2_grpc.ProductServiceServicer):
    
    def CreateProduct(self, request, context):
        db = SessionLocal()
        try:
            # TODO: Validate user exists via gRPC call to user service
            # For now, we'll assume the user_id is valid
            
            # Create new product
            product = Product(
                name=request.name,
                description=request.description,
                price=request.price,
                user_id=request.user_id
            )
            db.add(product)
            db.commit()
            db.refresh(product)
            
            logger.info(f"Created product: {product.id}")
            
            return product_pb2.CreateProductResponse(
                product=product_pb2.Product(
                    id=product.id,
                    name=product.name,
                    description=product.description,
                    price=float(product.price),
                    user_id=product.user_id,
                    created_at=product.created_at.isoformat()
                ),
                success=True,
                message="Product created successfully"
            )
            
        except Exception as e:
            logger.error(f"Error creating product: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.CreateProductResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def GetProduct(self, request, context):
        db = SessionLocal()
        try:
            product = db.query(Product).filter(Product.id == request.product_id).first()
            
            if not product:
                return product_pb2.GetProductResponse(found=False)
            
            return product_pb2.GetProductResponse(
                product=product_pb2.Product(
                    id=product.id,
                    name=product.name,
                    description=product.description,
                    price=float(product.price),
                    user_id=product.user_id,
                    created_at=product.created_at.isoformat()
                ),
                found=True
            )
            
        except Exception as e:
            logger.error(f"Error getting product: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.GetProductResponse(found=False)
        finally:
            db.close()
    
    def UpdateProduct(self, request, context):
        db = SessionLocal()
        try:
            product = db.query(Product).filter(Product.id == request.product_id).first()
            
            if not product:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Product not found")
                return product_pb2.UpdateProductResponse(
                    success=False,
                    message="Product not found"
                )
            
            # Update product fields
            if request.name:
                product.name = request.name
            if request.description:
                product.description = request.description
            if request.price > 0:
                product.price = request.price
            
            db.commit()
            db.refresh(product)
            
            logger.info(f"Updated product: {product.id}")
            
            return product_pb2.UpdateProductResponse(
                product=product_pb2.Product(
                    id=product.id,
                    name=product.name,
                    description=product.description,
                    price=float(product.price),
                    user_id=product.user_id,
                    created_at=product.created_at.isoformat()
                ),
                success=True,
                message="Product updated successfully"
            )
            
        except Exception as e:
            logger.error(f"Error updating product: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.UpdateProductResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def DeleteProduct(self, request, context):
        db = SessionLocal()
        try:
            product = db.query(Product).filter(Product.id == request.product_id).first()
            
            if not product:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Product not found")
                return product_pb2.DeleteProductResponse(
                    success=False,
                    message="Product not found"
                )
            
            db.delete(product)
            db.commit()
            
            logger.info(f"Deleted product: {request.product_id}")
            
            return product_pb2.DeleteProductResponse(
                success=True,
                message="Product deleted successfully"
            )
            
        except Exception as e:
            logger.error(f"Error deleting product: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.DeleteProductResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def ListProducts(self, request, context):
        db = SessionLocal()
        try:
            page = max(1, request.page or 1)
            limit = min(100, max(1, request.limit or 10))
            offset = (page - 1) * limit
            
            products = db.query(Product).offset(offset).limit(limit).all()
            total = db.query(Product).count()
            
            product_list = [
                product_pb2.Product(
                    id=product.id,
                    name=product.name,
                    description=product.description,
                    price=float(product.price),
                    user_id=product.user_id,
                    created_at=product.created_at.isoformat()
                )
                for product in products
            ]
            
            return product_pb2.ListProductsResponse(
                products=product_list,
                total=total,
                page=page,
                limit=limit
            )
            
        except Exception as e:
            logger.error(f"Error listing products: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.ListProductsResponse(products=[], total=0, page=1, limit=10)
        finally:
            db.close()
    
    def GetProductsByUser(self, request, context):
        db = SessionLocal()
        try:
            products = db.query(Product).filter(Product.user_id == request.user_id).all()
            
            product_list = [
                product_pb2.Product(
                    id=product.id,
                    name=product.name,
                    description=product.description,
                    price=float(product.price),
                    user_id=product.user_id,
                    created_at=product.created_at.isoformat()
                )
                for product in products
            ]
            
            return product_pb2.GetProductsByUserResponse(
                products=product_list,
                total=len(product_list)
            )
            
        except Exception as e:
            logger.error(f"Error getting products by user: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return product_pb2.GetProductsByUserResponse(products=[], total=0)
        finally:
            db.close()

def serve():
    port = os.getenv("GRPC_PORT", "50052")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    product_pb2_grpc.add_ProductServiceServicer_to_server(ProductService(), server)
    
    listen_addr = f"[::]:{port}"
    server.add_insecure_port(listen_addr)
    
    logger.info(f"Starting Product Service on {listen_addr}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down Product Service")
        server.stop(0)

if __name__ == "__main__":
    serve()