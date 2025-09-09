import grpc
from concurrent import futures
import user_pb2
import user_pb2_grpc
from models import User, SessionLocal
from datetime import datetime
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class UserService(user_pb2_grpc.UserServiceServicer):
    
    def CreateUser(self, request, context):
        db = SessionLocal()
        try:
            # Check if email already exists
            existing_user = db.query(User).filter(User.email == request.email).first()
            if existing_user:
                return user_pb2.CreateUserResponse(
                    success=False,
                    message="Email already exists"
                )
            
            # Create new user
            user = User(
                name=request.name,
                email=request.email,
                age=request.age
            )
            db.add(user)
            db.commit()
            db.refresh(user)
            
            logger.info(f"Created user: {user.id}")
            
            return user_pb2.CreateUserResponse(
                user=user_pb2.User(
                    id=user.id,
                    name=user.name,
                    email=user.email,
                    age=user.age,
                    created_at=user.created_at.isoformat()
                ),
                success=True,
                message="User created successfully"
            )
            
        except Exception as e:
            logger.error(f"Error creating user: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return user_pb2.CreateUserResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def GetUser(self, request, context):
        db = SessionLocal()
        try:
            user = db.query(User).filter(User.id == request.user_id).first()
            
            if not user:
                return user_pb2.GetUserResponse(found=False)
            
            return user_pb2.GetUserResponse(
                user=user_pb2.User(
                    id=user.id,
                    name=user.name,
                    email=user.email,
                    age=user.age,
                    created_at=user.created_at.isoformat()
                ),
                found=True
            )
            
        except Exception as e:
            logger.error(f"Error getting user: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return user_pb2.GetUserResponse(found=False)
        finally:
            db.close()
    
    def UpdateUser(self, request, context):
        db = SessionLocal()
        try:
            user = db.query(User).filter(User.id == request.user_id).first()
            
            if not user:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("User not found")
                return user_pb2.UpdateUserResponse(
                    success=False,
                    message="User not found"
                )
            
            # Update user fields
            user.name = request.name
            user.email = request.email
            user.age = request.age
            
            db.commit()
            db.refresh(user)
            
            logger.info(f"Updated user: {user.id}")
            
            return user_pb2.UpdateUserResponse(
                user=user_pb2.User(
                    id=user.id,
                    name=user.name,
                    email=user.email,
                    age=user.age,
                    created_at=user.created_at.isoformat()
                ),
                success=True,
                message="User updated successfully"
            )
            
        except Exception as e:
            logger.error(f"Error updating user: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return user_pb2.UpdateUserResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def DeleteUser(self, request, context):
        db = SessionLocal()
        try:
            user = db.query(User).filter(User.id == request.user_id).first()
            
            if not user:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("User not found")
                return user_pb2.DeleteUserResponse(
                    success=False,
                    message="User not found"
                )
            
            db.delete(user)
            db.commit()
            
            logger.info(f"Deleted user: {request.user_id}")
            
            return user_pb2.DeleteUserResponse(
                success=True,
                message="User deleted successfully"
            )
            
        except Exception as e:
            logger.error(f"Error deleting user: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return user_pb2.DeleteUserResponse(success=False, message="Internal error")
        finally:
            db.close()
    
    def ListUsers(self, request, context):
        db = SessionLocal()
        try:
            page = max(1, request.page or 1)
            limit = min(100, max(1, request.limit or 10))
            offset = (page - 1) * limit
            
            users = db.query(User).offset(offset).limit(limit).all()
            total = db.query(User).count()
            
            user_list = [
                user_pb2.User(
                    id=user.id,
                    name=user.name,
                    email=user.email,
                    age=user.age,
                    created_at=user.created_at.isoformat()
                )
                for user in users
            ]
            
            return user_pb2.ListUsersResponse(
                users=user_list,
                total=total,
                page=page,
                limit=limit
            )
            
        except Exception as e:
            logger.error(f"Error listing users: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return user_pb2.ListUsersResponse(users=[], total=0, page=1, limit=10)
        finally:
            db.close()

def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServiceServicer_to_server(UserService(), server)
    
    listen_addr = f"[::]:{port}"
    server.add_insecure_port(listen_addr)
    
    logger.info(f"Starting User Service on {listen_addr}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down User Service")
        server.stop(0)

if __name__ == "__main__":
    serve()