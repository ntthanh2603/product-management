import { NestFactory } from "@nestjs/core";
import { MicroserviceOptions, Transport } from "@nestjs/microservices";
import { AppModule } from "./app.module";
import { Logger } from "@nestjs/common";
import { join } from "path";

const grpcConfig = {
  transport: Transport.GRPC,
  options: {
    package: "user",
    protoPath: join(__dirname, "../../protos/user.proto"),
    url: `0.0.0.0:${process.env.GRPC_PORT || 50051}`,
  },
};

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
