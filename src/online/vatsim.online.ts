import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { Activity } from './activity';
import { OnlineService } from './online.interface';

@Injectable()
export class VatsimOnline implements OnlineService {
  private whazzupHost = 'https://data.vatsim.net/v3/vatsim-data.json';

  constructor(private readonly httpService: HttpService) {}

  getActivity(): Promise<Activity> {
    return this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => {
        const data = response.data;

        return {
          atc: data.controllers.length,
          pilot: data.pilots.length,
        };
      });
  }
}
