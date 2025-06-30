import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ClientProxyFactory } from '@nestjs/microservices';
import { ConfigService } from './services/config.service';
import { TasksController } from './tasks.controler';
import { WorkspaceController } from './workspace.controler';

@Module({
  imports: [],
  controllers: [AppController, TasksController, WorkspaceController],
  providers: [
    AppService,
    ConfigService,
    {
      provide: 'TASKS_SERVICE',
      useFactory: (configService: ConfigService) => {
        const ordersServiceOptions = configService.get('tasksService');
        return ClientProxyFactory.create(ordersServiceOptions);
      },
      inject: [ConfigService],
    },
    {
      provide: 'WORKSPACES_SERVICE',
      useFactory: (configService: ConfigService) => {
        const workspacesServiceOptions = configService.get('workspacesService');
        return ClientProxyFactory.create(workspacesServiceOptions);
      },
      inject: [ConfigService],
    },
  ],
})
export class AppModule {}
