import { Controller, Get, Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('workspaces')
@Controller('workspaces')
export class WorkspaceController {
  constructor(
    @Inject('WORKSPACES_SERVICE')
    private readonly workspacesService: ClientProxy,
  ) {}

  @Get()
  async getWorkspaces() {
    return this.workspacesService.send('get_workspaces', {});
  }
}
