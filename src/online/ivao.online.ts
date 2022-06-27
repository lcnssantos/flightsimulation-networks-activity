import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { Activity } from './activity';
import { OnlineService } from './online.interface';

@Injectable()
export class IVAOOnline implements OnlineService {
  private whazzupHost = 'https://api.ivao.aero/v2/tracker/whazzup';

  constructor(private readonly httpService: HttpService) {}

  getActivity(): Promise<Activity> {
    return this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => {
        const data = response.data;
        return {
          atc: data.clients.atcs.length,
          pilot: data.clients.pilots.length,
        };
      });
  }
}
