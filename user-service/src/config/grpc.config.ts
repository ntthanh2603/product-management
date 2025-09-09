import { Transport } from '@nestjs/microservices';
import { join } from 'path';

export const grpcConfig = {
  transport: Transport.GRPC,
  options: {
    package: 'user',
    protoPath: join(__dirname, '../../../protos/user.proto'),
    url: `0.0.0.0:${process.env.GRPC_PORT || 50051}`,
  },
};