import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { UserRoleEntity } from './userRole.entity';
import { Model } from 'mongoose';

@Injectable()
export class UserRoleService {
  constructor(
    @InjectModel(UserRoleEntity.name)
    private userRoleModel: Model<UserRoleEntity>,
  ) {}

  getAllUserRoles() {
    return this.userRoleModel.find();
  }

  createUserRole(userRole: UserRoleEntity) {
    userRole.creation_date = new Date();
    userRole.last_updated = new Date();
    return this.userRoleModel.create(userRole);
  }

  getRoleByName(name: string) {
    return this.userRoleModel.where({ name: name });
  }
}
