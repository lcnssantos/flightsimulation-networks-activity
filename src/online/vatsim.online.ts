import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { FirService } from 'src/firs/firs.service';
import { Activity } from './activity';
import { OnlineService } from './online.interface';

@Injectable()
export class VatsimOnline implements OnlineService {
  private whazzupHost = 'https://data.vatsim.net/v3/vatsim-data.json';

  constructor(
    private readonly httpService: HttpService,
    private firService: FirService,
  ) {}

  private getFirActivity(data: any, icao: string): Activity {
    const pilots = data.pilots.filter((pilot) => {
      return this.firService.isInsideFir(
        {
          lat: pilot.latitude,
          lng: pilot.longitude,
        },
        icao,
      );
    });

    return {
      atc: 0,
      pilot: pilots.length,
    };
  }

  async getBrazilActivity(): Promise<Activity> {
    const data = await this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => response.data);

    const firs = ['SBBS', 'SBCW', 'SBRE', 'SBAZ', 'SBAO'];

    const atcs = data.controllers.filter((atc) => {
      return atc.callsign.startsWith('SB') || atc.callsign.startsWith('SD');
    });

    await Promise.all(firs.map((fir) => this.firService.loadPointsByFIR(fir)));

    return Promise.all(firs.map((fir) => this.getFirActivity(data, fir)))
      .then((activities) =>
        activities.reduce(
          (acc, activity) => {
            return {
              atc: acc.atc + activity.atc,
              pilot: acc.pilot + activity.pilot,
            };
          },
          { atc: 0, pilot: 0 },
        ),
      )
      .then((a) => ({ atc: atcs.length, pilot: a.pilot }));
  }

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
