import { Inject, Injectable } from "@nestjs/common";
import { InjectModel } from "@nestjs/mongoose";
import { Model } from "mongoose";
import { TaskEntity } from "./task.entity";
import { ObjectId } from "mongodb";
import { UserEntity } from "../users/user/user.entity";
import { UserService } from "../users/user/user.service";

@Injectable()
export class TaskService {
    
    constructor(
        @InjectModel('TaskEntity') private readonly taskModel: Model<TaskEntity>,
        private readonly userService: UserService,
    ) { }

    create(task: TaskEntity, user: UserEntity) {
        task._id = new ObjectId();
        task.user = user;
        return this.taskModel.create(task);
    }

    getAll(user: UserEntity) {
        return this.taskModel.find({ user: user }).populate('user');
    }

    edit(_id: ObjectId, task: TaskEntity, user: UserEntity) {
        return this.taskModel.updateOne({ _id: _id, user: user }, task);
    }

    delete(_id: ObjectId, user: UserEntity) {
        return this.taskModel.deleteOne({ _id: _id, user: user });
    }
}