import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { DatabaseModule } from "@/database/database.module";
import { UserModule } from "@/user/user.module";
import { ProductModule } from "@/product/product.module";

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: ".env",
    }),
    DatabaseModule,
    UserModule,
    ProductModule,
  ],
})
export class AppModule {}
