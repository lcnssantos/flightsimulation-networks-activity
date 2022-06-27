import { Controller, Get, Post } from '@nestjs/common';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private appService: AppService) {}

  @Get('/current')
  getActivity() {
    return this.appService.getActivity();
  }

  @Post('/current')
  saveActivity() {
    return this.appService.saveActivity();
  }

  @Get('/history/24h')
  getHistory() {
    return this.appService.getHistory();
  }
}
