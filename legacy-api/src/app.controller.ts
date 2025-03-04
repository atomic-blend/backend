import { Controller, Options } from '@nestjs/common';
import { Public } from './utils/public_annotation';

@Controller('*')
export class AppController {
  constructor() {}

  @Options('')
  @Public()
  catchall() {
    return true;
  }
}
