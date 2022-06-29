import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { FirService } from 'src/firs/firs.service';
import { Activity } from './activity';
import { OnlineService } from './online.interface';

@Injectable()
export class PosconOnline implements OnlineService {
  private whazzupHost = 'https://hqapi.poscon.net/online.json';

  constructor(
    private readonly httpService: HttpService,
    private firService: FirService,
  ) {}

  private getFirActivity(data: any, icao: string): Activity {
    const pilots = data.flights.filter((pilot) => {
      return this.firService.isInsideFir(
        {
          lat: pilot.position.lat,
          lng: pilot.position.long,
        },
        icao,
      );
    });

    const atcs = data.atc.filter((atc) => atc.fir === icao);

    return {
      atc: atcs.length,
      pilot: pilots.length,
    };
  }

  async getBrazilActivity(): Promise<Activity> {
    const data = await this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => response.data);

    const firs = ['SBBS', 'SBCW', 'SBRE', 'SBAZ', 'SBAO'];

    await Promise.all(firs.map((fir) => this.firService.loadPointsByFIR(fir)));

    return Promise.all(firs.map((fir) => this.getFirActivity(data, fir))).then(
      (activities) =>
        activities.reduce(
          (acc, activity) => {
            return {
              atc: acc.atc + activity.atc,
              pilot: acc.pilot + activity.pilot,
            };
          },
          { atc: 0, pilot: 0 },
        ),
    );
  }

  getActivity(): Promise<Activity> {
    return this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => {
        const data = response.data;

        return {
          atc: data.atc.length,
          pilot: data.flights.length,
        };
      });
  }
}
