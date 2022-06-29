import { Controller, Get, Param, Post } from '@nestjs/common';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private appService: AppService) {}

  @Get('/current')
  getActivity() {
    return this.appService.getActivity();
  }

  @Get('/current/br')
  getBRActivity() {
    return this.appService.getBrazilActivity();
  }

  @Post('/current')
  async saveActivity() {
    await this.appService.saveActivity();
    await this.appService.saveActivityBR();
  }

  @Get('/history/24h')
  getHistory() {
    return this.appService.getHistoryByMinutes(24 * 60);
  }

  @Get('/history/:minutes')
  getHistoryByFilter(@Param('minutes') minutes: string) {
    return this.appService.getHistoryByMinutes(Number(minutes));
  }

  @Get('/history/br/24h')
  getBRHistory() {
    return this.appService.getBRHistoryByMinutes(24 * 60);
  }

  @Get('/history/br/:minutes')
  getBRHistoryByFilter(@Param('minutes') minutes: string) {
    return this.appService.getBRHistoryByMinutes(Number(minutes));
  }
}
