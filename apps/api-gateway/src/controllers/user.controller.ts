import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Body,
  Param,
  Query,
  Inject,
  OnModuleInit,
} from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { ApiTags } from '@nestjs/swagger';
import { Observable } from 'rxjs';
import {
  CreateUserDto,
  GetUsersQueryDto,
  UpdateUserDto,
} from 'src/dto/user.dto';

interface UserService {
  createUser(request: any): Observable<any>;
  getUser(request: any): Observable<any>;
  updateUser(request: any): Observable<any>;
  deleteUser(request: any): Observable<any>;
  getUsers(request: any): Observable<any>;
}

@ApiTags('Users')
@Controller('users')
export class UserController implements OnModuleInit {
  private userService: UserService;

  constructor(@Inject('USER_SERVICE') private client: ClientGrpc) {}

  onModuleInit() {
    this.userService = this.client.getService<UserService>('UserService');
  }

  @Post()
  createUser(@Body() createUserDto: CreateUserDto) {
    return this.userService.createUser(createUserDto);
  }

  @Get(':id')
  getUser(@Param('id') id: string) {
    return this.userService.getUser({ id: parseInt(id) });
  }

  @Put(':id')
  updateUser(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.userService.updateUser({
      id: parseInt(id),
      ...updateUserDto,
    });
  }

  @Delete(':id')
  deleteUser(@Param('id') id: string) {
    return this.userService.deleteUser({ id: parseInt(id) });
  }

  @Get()
  getUsers(@Query() query: GetUsersQueryDto) {
    return this.userService.getUsers({
      page: query.page ? parseInt(query.page) : 1,
      limit: query.limit ? parseInt(query.limit) : 10,
    });
  }
}
