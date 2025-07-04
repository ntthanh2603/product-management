import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { OrderModule } from './order.module';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    OrderModule,
    {
      transport: Transport.GRPC,
      options: {
        package: 'order',
        protoPath: join(__dirname, '../../../proto/order.proto'),
        url: 'localhost:5002',
      },
    },
  );

  await app.listen();
  console.log('Order Service is running on port 5002');
}
bootstrap();
