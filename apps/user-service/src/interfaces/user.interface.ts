export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  phone: string;
}

export interface GetUserRequest {
  id: number;
}

export interface UpdateUserRequest {
  id: number;
  name: string;
  email: string;
  phone: string;
}

export interface DeleteUserRequest {
  id: number;
}

export interface DeleteUserResponse {
  success: boolean;
  message: string;
}

export interface GetUsersRequest {
  page?: number;
  limit?: number;
}

export interface GetUsersResponse {
  users: User[];
  total: number;
}
