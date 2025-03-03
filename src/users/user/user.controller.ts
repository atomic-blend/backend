import { UserService } from './user.service';
import {
  BadRequestException,
  Body,
  Controller,
  Delete,
  Get,
  NotFoundException,
  Post,
  Put,
  Req,
  UseGuards,
} from '@nestjs/common';
import { UserEntity } from './user.entity';
import { UserRoleService } from '../userRole/userRole.service';
import { RoleName, UserRoleGuard } from '../../utils/userRole.guard';
import { UuidService } from 'nestjs-uuid';

@Controller('user')
export class UserController {
  constructor(
    private readonly userService: UserService,
    private readonly userRoleService: UserRoleService,
    private readonly uuidService: UuidService,
  ) { }

  @Get('/all')
  @RoleName('admin')
  @UseGuards(UserRoleGuard)
  async getUsers() {
    return this.userService.getUsers();
  }

  @Get()
  async getUser(@Req() request: any) {
    const user = await this.userService.getUserByFirebaseId(request.user.uid);
    if (user.length != 1) {
      throw new NotFoundException('User not found');
    }
    return user.at(0);
  }

  @Post('/setup')
  async setupUser(@Req() request: any, @Body() user: UserEntity) {
    const existingUser = await this.userService.getUserByFirebaseId(
      request.user.uid,
    );
    if (existingUser.length != 0) {
      throw new BadRequestException('User already exists');
    }
    user.firebase_id = request.user.uid;
    const role = await this.userRoleService.getRoleByName('user');
    user.roles = [role.at(0)];
    user.salt = this.uuidService.generate({ version: 4 });
    return await this.userService.createUser(user);
  }

  @Get('/storage/status')
  async getStorageStatus(@Req() req: any) {
    const user = await this.userService.getUserByFirebaseId(req.user.uid);
    if (user.length != 1) {
      throw new NotFoundException('User not found');
    }
  }

  @Delete('/delete')
  async deleteAccount(@Req() req: any) {
    const user = await this.userService.getUserByFirebaseId(req.user.uid);
    if (user.length != 1) {
      throw new NotFoundException('User not found');
    }
    return this.userService.deleteUserAndData(user.at(0));
  }

  @Post('/firstLogin')
  async firstLogin(@Req() req: any) {
    const firebaseUserId = req.user.uid;
    await this.userService.validateFirstLogin(firebaseUserId);
  }

  @Put('/updateDeviceToken')
  async updateDeviceToken(@Req() req: any, @Body() body: any) {
    const firebaseUserId = req.user.uid;
    return this.userService.updateDeviceToken(firebaseUserId, body.device_token);
  }
}
