import { Injectable, OnModuleInit, Logger, Inject } from "@nestjs/common";
import { ClientGrpc, ClientProxy } from "@nestjs/microservices";
import { Observable } from "rxjs";
import {
  ProductServiceClient,
  GetProductsByUserRequest,
  GetProductsByUserResponse,
} from "@/proto/product.pb";

@Injectable()
export class ProductClientService implements OnModuleInit {
  private readonly logger = new Logger(ProductClientService.name);
  private productService: ProductServiceClient;

  constructor(@Inject('PRODUCT_PACKAGE') private client: ClientGrpc) {}

  onModuleInit() {
    this.productService =
      this.client.getService<ProductServiceClient>("ProductService");
  }

  getProductsByUser(userId: number): Observable<GetProductsByUserResponse> {
    this.logger.log(`Getting products for user ${userId}`);
    return this.productService.getProductsByUser({ userId });
  }
}
