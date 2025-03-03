import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { ObjectId } from 'mongodb';

@Schema({ collection: 'user_roles' })
export class UserRoleEntity {
  @Prop()
  id: ObjectId;

  @Prop()
  name: string;

  @Prop()
  description?: string;

  @Prop()
  last_updated: Date;

  @Prop()
  creation_date: Date;
}

export const UserRoleEntitySchema =
  SchemaFactory.createForClass(UserRoleEntity);
