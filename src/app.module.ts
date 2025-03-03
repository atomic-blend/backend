import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { HealthController } from './health.controller';
import { TerminusModule } from '@nestjs/terminus';
import { ConfigModule } from '@nestjs/config';
import { APP_GUARD } from '@nestjs/core';
import { FirebaseMiddleware } from './middleware/firebase.middleware';
import { AppController } from './app.controller';
import { HttpModule } from '@nestjs/axios';
import { UserController } from './users/user/user.controller';
import { UserService } from './users/user/user.service';
import { UserRoleService } from './users/userRole/userRole.service';
import { UserRoleController } from './users/userRole/userRole.controller';
import { UserRoleEntity, UserRoleEntitySchema } from './users/userRole/userRole.entity';
import { UserEntity, UserEntitySchema } from './users/user/user.entity';
import { UuidModule } from 'nestjs-uuid';
import { TaskEntity, TaskEntitySchema } from './tasks/task.entity';
import { TaskController } from './tasks/task.controller';
import { TaskService } from './tasks/tasks.service';

const models = [
  {
    name: UserEntity.name,
    schema: UserEntitySchema,
  },
  {
    name: UserRoleEntity.name,
    schema: UserRoleEntitySchema,
  },
  {
    name: TaskEntity.name,
    schema: TaskEntitySchema,
  }
];

const controllers = [
  HealthController,
  AppController,
  TaskController,
  UserController,
  UserRoleController,
]

const services = [
  {
    provide: APP_GUARD,
    useClass: FirebaseMiddleware,
  },
  UserService,
  TaskService,
  UserRoleService,
];

@Module({
  imports: [
    TerminusModule,
    HttpModule,
    UuidModule,
    ConfigModule.forRoot(),
    MongooseModule.forRoot(process.env.MONGODB_URI),
    MongooseModule.forFeature(models)
  ],
  controllers: [...controllers],
  providers: [...services],
})
export class AppModule { }
