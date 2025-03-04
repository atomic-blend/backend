import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { ObjectId } from "mongodb";
import { UserEntity } from "../users/user/user.entity";
import mongoose from "mongoose";

@Schema({ collection: 'tasks', timestamps: true })
export class TaskEntity {
    @Prop()
    _id: ObjectId;

    @Prop()
    title: string;

    @Prop(
        {
            type: mongoose.Schema.Types.ObjectId,
            ref: UserEntity.name
        }
    )
    user: UserEntity;

    @Prop()
    description?: string;

    @Prop()
    start_date?: string;

    @Prop()
    completed?: boolean;

    @Prop()
    created_at: Date;

    @Prop()
    updated_at: Date;
}

export const TaskEntitySchema = SchemaFactory.createForClass(TaskEntity);

// convert the _id to id
TaskEntitySchema.set('toJSON', {
    transform: (doc, ret) => {
        ret.id = ret._id;
        delete ret._id;
        delete ret.__v;
    }
});