import { Controller, Get, Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('tasks')
@Controller('tasks')
export class TasksController {
  constructor(
    @Inject('TASKS_SERVICE') private readonly tasksService: ClientProxy,
  ) {}

  @Get()
  async getTasks(data: string = 'Get tasks from API Gateway!') {
    return this.tasksService.send('get_tasks', { data });
  }
}
