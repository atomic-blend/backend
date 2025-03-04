import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { UserEntity } from './user.entity';
import { Model } from 'mongoose';
import { ObjectId } from 'mongodb';
import * as firebase from 'firebase-admin';

@Injectable()
export class UserService {
  private firebaseAdmin: any;
  constructor(
    @InjectModel(UserEntity.name) private userModel: Model<UserEntity>,
  ) {
    this.firebaseAdmin = firebase.apps[0];
  }

  getUserByFirebaseId(firebase_id: string) {
    return this.userModel
      .where({ firebase_id: firebase_id })
      .populate('roles')
      .exec();
  }

  getUsers() {
    return this.userModel.find().populate('roles').populate('purchases').exec();
  }

  getUser(id: string) {
    return this.userModel
      .where({ _id: new ObjectId(id) })
      .populate('roles')
      .exec();
  }

  getUserByEmail(email: string) {
    return this.userModel
      .where({ email })
      .populate('roles')
      .exec();
  }

  async createUser(user: UserEntity) {
    user._id = new ObjectId();
    user.last_updated = new Date();
    user.creation_date = new Date();
    await this.userModel.create(user);
    return user;
  }

  async edit(_id: ObjectId, user: UserEntity) {
    user.last_updated = new Date();
    return this.userModel.updateOne({ _id: _id }, user);
  }

  async deleteUserAndData(user: UserEntity) {
    await this.userModel.deleteOne({ _id: user._id });
    await firebase.apps[0].auth().deleteUser(user.firebase_id);
  }

  async createUserOnFirebase(email: string, password: string) {
    const newUser = await this.firebaseAdmin.auth().createUser({
      email: email,
      password: password,
      disabled: false,
    });
    return newUser.uid;
  }

  async validateFirstLogin(firebase_id: string) {
    const user = (await this.getUserByFirebaseId(firebase_id)).at(0);
    if (user == null) {
      throw new NotFoundException();
    }
    user.first_login = false;
    await this.edit(user._id, user);
  }

  async updateDeviceToken(firebase_id: string, device_token: string) {
    const user = (await this.getUserByFirebaseId(firebase_id)).at(0);
    if (user == null) {
      throw new NotFoundException();
    }
    if (user.device_token == null) {
      user.device_token = [];
    }
    user.device_token.push(device_token);
    await this.edit(user._id, user);
  } 
}
