# NestJS Microservice

![NestJS](https://img.shields.io/badge/nestjs-%23E0234E.svg?style=for-the-badge&logo=nestjs&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white)
![Node.js](https://img.shields.io/badge/node.js-6DA55F?style=for-the-badge&logo=node.js&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

### Run project Tasks Manager use NestJS Microservice with transport TCP

```bash
docker compose up -d --build
```

### Functions in NestJS Microservice TCP

#### Send message and get response

```bash
this.client.send('get_user', { id: 1 }).subscribe(response => {
  console.log(response);
});
```

#### Send message but not get response

```bash
this.client.emit('user_created', { userId: 1, name: 'TuanThanh' });
```

#### Connect to microservice

```bash
await this.client.connect();
```

#### Unconnect to microservice

```bash
this.client.close();
```

### Decorators in NestJS Microservice

#### Handle request-response message

```bash
@MessagePattern('get_user')
getUser(data: { id: number }) {
  return this.userService.findById(data.id);
}
```

#### Handle event

```bash
@EventPattern('user_created')
handleUserCreated(data: { userId: number, name: string }) {
  // Handle event
}
```

#### Inject ClientProxy instance

```bash
@Client({
  transport: Transport.TCP,
  options: {
    host: 'localhost',
    port: 3001,
  },
})
client: ClientProxy;
```

#### Get data in message payload

```bash
@MessagePattern('get_user')
getUser(@Payload() data: { id: number }) {
  return data;
}
```

#### Get context of message

```bash
@MessagePattern('get_user')
getUser(@Payload() data: any, @Ctx() context: TcpContext) {
  console.log(context.getPattern());
}
```
