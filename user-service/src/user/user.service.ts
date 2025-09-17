import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { User } from "@/database/user.entity";
import {
  CreateUserRequest,
  CreateUserResponse,
  GetUserRequest,
  GetUserResponse,
  UpdateUserRequest,
  UpdateUserResponse,
  DeleteUserRequest,
  DeleteUserResponse,
  ListUsersRequest,
  ListUsersResponse,
} from "@/proto/user.pb";

@Injectable()
export class UserService {
  private readonly logger = new Logger(UserService.name);

  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>
  ) {}

  async createUser(request: CreateUserRequest): Promise<CreateUserResponse> {
    try {
      // Check if email already exists
      const existingUser = await this.userRepository.findOne({
        where: { email: request.email },
      });

      if (existingUser) {
        return {
          user: undefined,
          success: false,
          message: "Email already exists",
        };
      }

      // Create new user
      const user = this.userRepository.create({
        name: request.name,
        email: request.email,
        age: request.age,
      });

      const savedUser = await this.userRepository.save(user);
      this.logger.log(`Created user: ${savedUser.id}`);

      return {
        user: {
          id: savedUser.id,
          name: savedUser.name,
          email: savedUser.email,
          age: savedUser.age,
          createdAt: savedUser.createdAt.toISOString(),
        },
        success: true,
        message: "User created successfully",
      };
    } catch (error) {
      this.logger.error(`Error creating user: ${error.message}`);
      return {
        user: undefined,
        success: false,
        message: "Internal error",
      };
    }
  }

  async getUser(request: GetUserRequest): Promise<GetUserResponse> {
    try {
      const user = await this.userRepository.findOne({
        where: { id: request.userId },
      });

      if (!user) {
        return {
          user: undefined,
          found: false,
        };
      }

      return {
        user: {
          id: user.id,
          name: user.name,
          email: user.email,
          age: user.age,
          createdAt: user.createdAt.toISOString(),
        },
        found: true,
      };
    } catch (error) {
      this.logger.error(`Error getting user: ${error.message}`);
      return {
        user: undefined,
        found: false,
      };
    }
  }

  async updateUser(request: UpdateUserRequest): Promise<UpdateUserResponse> {
    try {
      const user = await this.userRepository.findOne({
        where: { id: request.userId },
      });

      if (!user) {
        return {
          user: undefined,
          success: false,
          message: "User not found",
        };
      }

      // Update user fields
      user.name = request.name;
      user.email = request.email;
      user.age = request.age;

      const savedUser = await this.userRepository.save(user);
      this.logger.log(`Updated user: ${savedUser.id}`);

      return {
        user: {
          id: savedUser.id,
          name: savedUser.name,
          email: savedUser.email,
          age: savedUser.age,
          createdAt: savedUser.createdAt.toISOString(),
        },
        success: true,
        message: "User updated successfully",
      };
    } catch (error) {
      this.logger.error(`Error updating user: ${error.message}`);
      return {
        user: undefined,
        success: false,
        message: "Internal error",
      };
    }
  }

  async deleteUser(request: DeleteUserRequest): Promise<DeleteUserResponse> {
    try {
      const user = await this.userRepository.findOne({
        where: { id: request.userId },
      });

      if (!user) {
        return {
          success: false,
          message: "User not found",
        };
      }

      await this.userRepository.remove(user);
      this.logger.log(`Deleted user: ${request.userId}`);

      return {
        success: true,
        message: "User deleted successfully",
      };
    } catch (error) {
      this.logger.error(`Error deleting user: ${error.message}`);
      return {
        success: false,
        message: "Internal error",
      };
    }
  }

  async listUsers(request: ListUsersRequest): Promise<ListUsersResponse> {
    try {
      const page = Math.max(1, request.page || 1);
      const limit = Math.min(100, Math.max(1, request.limit || 10));
      const skip = (page - 1) * limit;

      const [users, total] = await this.userRepository.findAndCount({
        skip,
        take: limit,
        order: { createdAt: "DESC" },
      });

      const userList = users.map((user) => ({
        id: user.id,
        name: user.name,
        email: user.email,
        age: user.age,
        createdAt: user.createdAt.toISOString(),
      }));

      return {
        users: userList,
        total,
        page,
        limit,
      };
    } catch (error) {
      this.logger.error(`Error listing users: ${error.message}`);
      return {
        users: [],
        total: 0,
        page: 1,
        limit: 10,
      };
    }
  }
}
