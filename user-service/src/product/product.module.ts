import { Module } from "@nestjs/common";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { join } from "path";
import { ProductClientService } from "./product-client.service";

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'PRODUCT_PACKAGE',
        transport: Transport.GRPC,
        options: {
          package: 'product',
          protoPath: join(__dirname, "../../../protos/product.proto"),
          url: process.env.PRODUCT_SERVICE_URL || "localhost:50052",
        },
      },
    ]),
  ],
  providers: [ProductClientService],
  exports: [ProductClientService],
})
export class ProductModule {}
