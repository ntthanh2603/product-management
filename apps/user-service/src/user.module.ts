import { Module } from '@nestjs/common';
import { UserService } from './user.service';

@Module({
  imports: [],
  controllers: [],
  providers: [UserService],
  exports: [],
})
export class UserModule {}
