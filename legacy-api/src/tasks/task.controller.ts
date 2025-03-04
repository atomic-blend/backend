import { Body, Controller, Delete, Get, Param, Post, Put, Req } from "@nestjs/common";
import { TaskService } from "./tasks.service";
import { TaskEntity } from "./task.entity";
import { UserService } from "../users/user/user.service";
import { ObjectId } from "mongodb";

@Controller('tasks')
export class TaskController {
    constructor(
        private readonly taskService: TaskService,
        private readonly userService: UserService,
    ) { }

    @Get()
    async getTasks(@Req() request: any) {
        const user = await this.userService.getUserByFirebaseId(request.user.uid);
        return this.taskService.getAll(user.at(0));
    }

    @Post()
    async createTask(@Req() req: any, @Body() task: TaskEntity) {
        const user = await this.userService.getUserByFirebaseId(req.user.uid);
        return this.taskService.create(task, user.at(0));
    }

    @Put('/:id')
    async editTask(@Req() req: any, @Body() task: TaskEntity, @Param('id') id: string) {
        const user = await this.userService.getUserByFirebaseId(req.user.uid);
        return this.taskService.edit(new ObjectId(id), task, user.at(0));
    }

    @Delete('/:id')
    async deleteTask(@Req() req: any, @Param('id') id: string) {
        const user = await this.userService.getUserByFirebaseId(req.user.uid);
        return this.taskService.delete(new ObjectId(id), user.at(0));
    }
}
