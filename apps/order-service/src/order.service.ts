import { Injectable, Inject, OnModuleInit } from '@nestjs/common';
import { ClientGrpc, GrpcMethod } from '@nestjs/microservices';
import { Observable, firstValueFrom } from 'rxjs';

interface UserService {
  getUser(request: any): Observable<any>;
}

@Injectable()
export class OrderService implements OnModuleInit {
  private userService: UserService;
  private orders: Order[] = [];
  private idCounter = 1;

  constructor(@Inject('USER_SERVICE') private userClient: ClientGrpc) {}

  onModuleInit() {
    this.userService = this.userClient.getService<UserService>('UserService');
  }

  @GrpcMethod('OrderService', 'CreateOrder')
  async createOrder(request: CreateOrderRequest): Promise<Order> {
    // Verify user exists
    try {
      await firstValueFrom(this.userService.getUser({ id: request.userId }));
    } catch (error) {
      throw new Error('User not found');
    }

    const totalAmount = request.items.reduce(
      (total, item) => total + item.price * item.quantity,
      0,
    );

    const order: Order = {
      id: this.idCounter++,
      userId: request.userId,
      items: request.items,
      totalAmount,
      status: 'pending',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    this.orders.push(order);
    return order;
  }

  // ... other methods
}
