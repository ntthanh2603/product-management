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
  async getTasks() {
    return this.tasksService.send({ cmd: 'get_tasks' }, {});
  }
}
