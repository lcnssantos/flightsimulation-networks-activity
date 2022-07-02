import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { FirService, Point } from 'src/firs/firs.service';
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
export class VatsimOnline implements OnlineService {
  private whazzupHost = 'https://data.vatsim.net/v3/vatsim-data.json';
  private transceiverHost = 'https://data.vatsim.net/v3/transceivers-data.json';
  private transceiverData: Map<string, Point> = new Map();

  constructor(
    private readonly httpService: HttpService,
    private firService: FirService,
  ) {
    this.httpService
      .get(this.transceiverHost)
      .toPromise()
      .then((r) => {
        const data = r.data;
        for (const transceiver of data) {
          if (transceiver.transceivers.length > 0) {
            this.transceiverData.set(transceiver.callsign, {
              lat: transceiver.transceivers[0].latDeg,
              lng: transceiver.transceivers[0].lonDeg,
            });
          }
        }
      });
  }

  async getActivityByRegion(): Promise<GeoActivity> {
    await this.firService.loadFirData();

    const data = await this.httpService
      .get(this.whazzupHost)
      .toPromise()
      .then((response) => response.data);

    const count = new Count();

    for (const pilot of data.pilots) {
      try {
        if (!pilot.latitude || !pilot.longitude) {
          throw new Error();
        }

        const fir = this.firService.detectCountryByPoint({
          lat: pilot.latitude,
          lng: pilot.longitude,
        });
        count.increment(fir, 'pilot');
      } catch {
        count.increment(UNKNOWN, 'pilot');
      }
    }

    for (const atc of data.controllers) {
      try {
        const point = this.transceiverData.get(atc.callsign);

        if (!point) {
          throw new Error();
        }

        const country = this.firService.detectCountryByPoint(point);

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
    const pilots = data.pilots.filter((pilot) => {
      if (!pilot.latitude || !pilot.longitude) {
        return false;
      }

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
          atc: data.controllers.length,
          pilot: data.pilots.length,
        };
      });
  }
}
