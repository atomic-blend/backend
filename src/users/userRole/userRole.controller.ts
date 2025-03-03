import { Body, Controller, Get, Post, UseGuards } from '@nestjs/common';
import { UserRoleService } from './userRole.service';
import { UserRoleEntity } from './userRole.entity';
import { RoleName, UserRoleGuard } from '../../utils/userRole.guard';

@Controller('user/role')
export class UserRoleController {
  constructor(private readonly userRoleService: UserRoleService) {}

  @Get()
  @RoleName('admin')
  @UseGuards(UserRoleGuard)
  async getAllUserRoles() {
    return this.userRoleService.getAllUserRoles();
  }

  @Post()
  @RoleName('admin')
  @UseGuards(UserRoleGuard)
  async createUserRole(@Body() userRole: UserRoleEntity) {
    return this.userRoleService.createUserRole(userRole);
  }
}
