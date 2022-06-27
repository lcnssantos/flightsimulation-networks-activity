import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { Activity } from './activity';
import { OnlineService } from './online.interface';

@Injectable()
export class PosconOnline implements OnlineService {
  private whazzupHost = 'https://hqapi.poscon.net/online.json';

  constructor(private readonly httpService: HttpService) {}

  getActivity(): Promise<Activity> {
    return this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => {
        const data = response.data;

        return {
          atc: data.atc.length,
          pilot: data.flights.length,
          time: new Date(),
        };
      });
  }
}
