import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { FirService } from 'src/firs/firs.service';
import { Activity, GeoActivity } from './activity';
import { OnlineService } from './online.interface';

const UNKNOWN = 'UNKNOWN';
class Count {
  private map: Map<string, Activity> = new Map();

  public increment(id: string, type: string) {
    const count = this.map.get(id) || { atc: 0, pilot: 0 };
    if (type === 'atc') {
      this.map.set(id, { ...count, atc: count.atc + 1 });
    } else if (type === 'pilot') {
      this.map.set(id, { ...count, pilot: count.pilot + 1 });
    }
  }

  public get() {
    return this.map;
  }
}

@Injectable()
export class IVAOOnline implements OnlineService {
  private whazzupHost = 'https://api.ivao.aero/v2/tracker/whazzup';

  constructor(
    private readonly httpService: HttpService,
    private firService: FirService,
  ) {}

  async getActivityByRegion(): Promise<GeoActivity> {
    await this.firService.loadFirData();

    const data = await this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => response.data);

    const count = new Count();

    for (const pilot of data.clients.pilots) {
      try {
        if (!pilot.lastTrack) {
          throw new Error();
        }

        const country = this.firService.detectCountryByPoint({
          lat: pilot.lastTrack.latitude,
          lng: pilot.lastTrack.longitude,
        });
        count.increment(country, 'pilot');
      } catch {
        count.increment(UNKNOWN, 'pilot');
      }
    }

    for (const atc of data.clients.atcs) {
      try {
        if (!atc.lastTrack) {
          throw new Error();
        }

        const country = this.firService.detectCountryByPoint({
          lat: atc.lastTrack.latitude,
          lng: atc.lastTrack.longitude,
        });

        count.increment(country, 'atc');
      } catch {
        count.increment(UNKNOWN, 'atc');
      }
    }

    const map = count.get();

    return [...map.keys()].reduce((acc, key) => {
      acc[key] = map.get(key);
      return acc;
    }, {});
  }

  private getFirActivity(data: any, icao: string): Activity {
    const pilots = data.clients.pilots.filter((pilot) => {
      if (!pilot.lastTrack) {
        return false;
      }

      return this.firService.isInsideFir(
        {
          lat: pilot.lastTrack.latitude,
          lng: pilot.lastTrack.longitude,
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

    const atcs = data.clients.atcs.filter((atc) => {
      return atc.callsign.startsWith('SB') || atc.callsign.startsWith('SD');
    });

    await this.firService.loadFirData();

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
          atc: data.clients.atcs.length,
          pilot: data.clients.pilots.length,
        };
      });
  }
}
