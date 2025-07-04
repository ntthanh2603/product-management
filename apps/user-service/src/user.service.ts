import { Injectable } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import {
  CreateUserRequest,
  DeleteUserRequest,
  DeleteUserResponse,
  GetUserRequest,
  GetUsersRequest,
  GetUsersResponse,
  UpdateUserRequest,
  User,
} from './interfaces/user.interface';

@Injectable()
export class UserService {
  private users: User[] = [];
  private idCounter = 1;

  @GrpcMethod('UserService', 'CreateUser')
  createUser(request: CreateUserRequest): User {
    const user: User = {
      id: this.idCounter++,
      name: request.name,
      email: request.email,
      phone: request.phone,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    this.users.push(user);
    return user;
  }

  @GrpcMethod('UserService', 'GetUser')
  getUser(request: GetUserRequest): User {
    const user = this.users.find((u) => u.id === request.id);
    if (!user) {
      throw new Error('User not found');
    }
    return user;
  }

  @GrpcMethod('UserService', 'UpdateUser')
  updateUser(request: UpdateUserRequest): User {
    const userIndex = this.users.findIndex((u) => u.id === request.id);
    if (userIndex === -1) {
      throw new Error('User not found');
    }

    this.users[userIndex] = {
      ...this.users[userIndex],
      name: request.name,
      email: request.email,
      phone: request.phone,
      updatedAt: new Date().toISOString(),
    };

    return this.users[userIndex];
  }

  @GrpcMethod('UserService', 'DeleteUser')
  deleteUser(request: DeleteUserRequest): DeleteUserResponse {
    const userIndex = this.users.findIndex((u) => u.id === request.id);
    if (userIndex === -1) {
      return {
        success: false,
        message: 'User not found',
      };
    }

    this.users.splice(userIndex, 1);
    return {
      success: true,
      message: 'User deleted successfully',
    };
  }

  @GrpcMethod('UserService', 'GetUsers')
  getUsers(request: GetUsersRequest): GetUsersResponse {
    const page = request.page || 1;
    const limit = request.limit || 10;
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;

    const users = this.users.slice(startIndex, endIndex);

    return {
      users,
      total: this.users.length,
    };
  }
}
