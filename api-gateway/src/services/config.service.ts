import { Transport } from '@nestjs/microservices';

export class ConfigService {
  private readonly envConfig: { [key: string]: any } = {};

  constructor() {
    this.envConfig = {};
    this.envConfig.port = 3000;
    this.envConfig.workspacesService = {
      options: {
        port: 3001,
        host: 'localhost',
      },
      transport: Transport.TCP,
    };
    this.envConfig.tasksService = {
      options: {
        port: 3002,
        host: 'localhost',
      },
      transport: Transport.TCP,
    };
  }

  get(key: string): any {
    return this.envConfig[key];
  }
}
