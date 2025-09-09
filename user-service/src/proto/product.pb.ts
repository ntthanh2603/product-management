/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";

export const protobufPackage = "product";

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  userId: number;
  createdAt: string;
}

export interface CreateProductRequest {
  name: string;
  description: string;
  price: number;
  userId: number;
}

export interface CreateProductResponse {
  product: Product | undefined;
  success: boolean;
  message: string;
}

export interface GetProductRequest {
  productId: number;
}

export interface GetProductResponse {
  product: Product | undefined;
  found: boolean;
}

export interface UpdateProductRequest {
  productId: number;
  name: string;
  description: string;
  price: number;
}

export interface UpdateProductResponse {
  product: Product | undefined;
  success: boolean;
  message: string;
}

export interface DeleteProductRequest {
  productId: number;
}

export interface DeleteProductResponse {
  success: boolean;
  message: string;
}

export interface ListProductsRequest {
  page: number;
  limit: number;
}

export interface ListProductsResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
}

export interface GetProductsByUserRequest {
  userId: number;
}

export interface GetProductsByUserResponse {
  products: Product[];
  total: number;
}

export const PRODUCT_PACKAGE_NAME = "product";

export interface ProductServiceClient {
  createProduct(
    request: CreateProductRequest
  ): Observable<CreateProductResponse>;

  getProduct(request: GetProductRequest): Observable<GetProductResponse>;

  updateProduct(
    request: UpdateProductRequest
  ): Observable<UpdateProductResponse>;

  deleteProduct(
    request: DeleteProductRequest
  ): Observable<DeleteProductResponse>;

  listProducts(request: ListProductsRequest): Observable<ListProductsResponse>;

  getProductsByUser(
    request: GetProductsByUserRequest
  ): Observable<GetProductsByUserResponse>;
}
