import grpc
import user_pb2
import user_pb2_grpc

def test_user_service():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)
        
        # Create user
        create_response = stub.CreateUser(user_pb2.CreateUserRequest(
            name="John Doe",
            email="john@example.com",
            age=30
        ))
        print(f"Create User: {create_response.success}, {create_response.message}")
        
        if create_response.success:
            user_id = create_response.user.id
            
            # Get user
            get_response = stub.GetUser(user_pb2.GetUserRequest(user_id=user_id))
            print(f"Get User: {get_response.found}, Name: {get_response.user.name}")
            
            # Update user
            update_response = stub.UpdateUser(user_pb2.UpdateUserRequest(
                user_id=user_id,
                name="John Updated",
                email="john.updated@example.com",
                age=31
            ))
            print(f"Update User: {update_response.success}")
            
            # List users
            list_response = stub.ListUsers(user_pb2.ListUsersRequest(page=1, limit=10))
            print(f"List Users: Total {list_response.total}")
            
            # Delete user
            delete_response = stub.DeleteUser(user_pb2.DeleteUserRequest(user_id=user_id))
            print(f"Delete User: {delete_response.success}")

if __name__ == "__main__":
    test_user_service()
