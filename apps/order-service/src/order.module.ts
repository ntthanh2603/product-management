import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { OrderService } from './order.service';
import { DatabaseModule } from './database/database.module';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'USER_SERVICE',
        transport: Transport.GRPC,
        options: {
          package: 'user',
          protoPath: join(__dirname, '../../../proto/user.proto'),
          url: 'localhost:5001',
        },
      },
    ]),
  ],
  controllers: [],
  providers: [OrderService, DatabaseModule],
})
export class OrderModule {}
