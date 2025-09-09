import { Controller } from "@nestjs/common";
import { UserService } from "./user.service";
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
  UserServiceControllerMethods,
} from "@/proto/user.pb";

@Controller()
@UserServiceControllerMethods()
export class UserController {
  constructor(private readonly userService: UserService) {}

  async createUser(request: CreateUserRequest): Promise<CreateUserResponse> {
    return this.userService.createUser(request);
  }

  async getUser(request: GetUserRequest): Promise<GetUserResponse> {
    return this.userService.getUser(request);
  }

  async updateUser(request: UpdateUserRequest): Promise<UpdateUserResponse> {
    return this.userService.updateUser(request);
  }

  async deleteUser(request: DeleteUserRequest): Promise<DeleteUserResponse> {
    return this.userService.deleteUser(request);
  }

  async listUsers(request: ListUsersRequest): Promise<ListUsersResponse> {
    return this.userService.listUsers(request);
  }
}
