import { NestFactory } from "@nestjs/core";
import { MicroserviceOptions } from "@nestjs/microservices";
import { AppModule } from "./app.module";
import { grpcConfig } from "@/config/grpc.config";
import { Logger } from "@nestjs/common";

async function bootstrap() {
  const logger = new Logger("UserService");

  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    grpcConfig as MicroserviceOptions
  );

  await app.listen();
  logger.log(
    `User Service is listening on port ${process.env.GRPC_PORT || 50051}`
  );
}

bootstrap();
