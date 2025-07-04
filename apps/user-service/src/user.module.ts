import { Module } from '@nestjs/common';
import { UserService } from './user.service';
import { DatabaseModule } from './database/database.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: ['.env'],
    }),
  ],
  controllers: [],
  providers: [UserService, DatabaseModule],
  exports: [],
})
export class UserModule {}
