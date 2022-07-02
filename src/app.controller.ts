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

  @Get('/current/geo')
  getGeoActivity() {
    return this.appService.getGeoActivity();
  }

  @Post('/current')
  async saveActivity() {
    await this.appService.saveActivity();
    await this.appService.saveActivityBR();
    await this.appService.saveActivityByRegion();
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

  @Get('/history/geo/24h')
  getGeoHistory() {
    return this.appService.getGeoHistoryByMinutes(24 * 60);
  }

  @Get('/history/geo/:minutes')
  getGeoHistoryByFilter(@Param('minutes') minutes: string) {
    return this.appService.getGeoHistoryByMinutes(Number(minutes));
  }
}
