import {
  CanActivate,
  ExecutionContext,
  Injectable,
  SetMetadata,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { UserService } from '../users/user/user.service';
import { Reflector } from '@nestjs/core';
import { decode } from 'jsonwebtoken';

@Injectable()
export class UserRoleGuard implements CanActivate {
  constructor(
    private readonly userService: UserService,
    private reflector: Reflector,
  ) {}

  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();
    return this.validateRequest(context, request);
  }

  async validateRequest(
    context: ExecutionContext,
    request: Request,
  ): Promise<boolean> {
    const roleName = this.reflector.get<string>(
      'roleName',
      context.getHandler(),
    );
    const jwtPayload = decode(
      request.headers['authorization'].replace('Bearer ', ''),
    );
    const firebase_id = jwtPayload['user_id'];
    const users = await this.userService.getUserByFirebaseId(firebase_id);
    return users.at(0).roles.some((role) => role.name === roleName);
  }
}

export const RoleName = (roleName: string) => SetMetadata('roleName', roleName);
