import { Controller, Get, Param, Post } from '@nestjs/common';
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
    return this.appService.getHistoryByMinutes(24 * 60);
  }

  @Get('/history/:minutes')
  getHistoryByFilter(@Param('minutes') minutes: string) {
    return this.appService.getHistoryByMinutes(Number(minutes));
  }
}
