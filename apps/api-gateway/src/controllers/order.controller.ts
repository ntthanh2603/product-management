import { Controller, Get, Post, Body, Inject, Param } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';

interface OrderService {
  CreateOrder(data: any): Observable<any>;
  GetOrder(data: any): Observable<any>;
  GetOrdersByUser(data: any): Observable<any>;
}

@Controller('orders')
export class OrderController {
  private orderService: OrderService;

  constructor(@Inject('ORDER_SERVICE') private client: ClientGrpc) {}

  onModuleInit() {
    this.orderService = this.client.getService<OrderService>('OrderService');
  }

  @Post()
  createOrder(@Body() body: any) {
    return this.orderService.CreateOrder(body);
  }

  @Get(':id')
  getOrder(@Param('id') id: number) {
    return this.orderService.GetOrder({ id: Number(id) });
  }

  @Get('user/:userId')
  getOrdersByUser(@Param('userId') userId: number) {
    return this.orderService.GetOrdersByUser({ user_id: Number(userId) });
  }
}
