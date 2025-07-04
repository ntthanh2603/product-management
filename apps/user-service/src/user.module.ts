import { Module } from '@nestjs/common';
import { UserService } from './user.service';
import { DatabaseModule } from './database/database.module';

@Module({
  imports: [],
  controllers: [],
  providers: [UserService, DatabaseModule],
  exports: [],
})
export class UserModule {}
