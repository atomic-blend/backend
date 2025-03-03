import { ObjectId } from 'mongodb';
import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { UserRoleEntity } from '../userRole/userRole.entity';
import mongoose from 'mongoose';

@Schema({ collection: 'users' })
export class UserEntity {
  @Prop()
  _id: ObjectId;

  @Prop()
  firebase_id: string;

  @Prop({
    type: [{ type: mongoose.Schema.Types.ObjectId, ref: UserRoleEntity.name }],
  })
  roles: UserRoleEntity[];

  @Prop()
  email: string;

  @Prop()
  firstname?: string;

  @Prop()
  lastname?: string;

  @Prop()
  first_login?: boolean;

  @Prop()
  salt?: string;

  @Prop()
  device_token?: string[];

  @Prop()
  last_updated: Date;

  @Prop()
  creation_date: Date;

}

export const UserEntitySchema = SchemaFactory.createForClass(UserEntity);
UserEntitySchema.set('toJSON', {
  transform: (doc, ret) => {
      ret.id = ret._id;
      delete ret._id;
      delete ret.__v;
  }
})