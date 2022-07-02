import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { readFile } from 'fs/promises';
import { GeoLocator } from 'src/geo/geo.locator';

export class Point {
  lat: number;
  lng: number;
}

interface Fir {
  icao: string;
  region: string;
  points: Array<Point>;
}

@Injectable()
export class FirService {
  private endpoint = 'https://map.vatsim.net/livedata/firboundaries.json';

  private firsMap: Map<string, Fir>;
  private firCountry: Map<string, string>;
  private responseData: any;

  constructor(
    private readonly httpService: HttpService,
    private geoService: GeoLocator,
  ) {
    this.firsMap = new Map();
    this.firCountry = new Map();
  }

  async loadFirData() {
    if (!this.responseData) {
      const response = await this.httpService.get(this.endpoint).toPromise();

      this.responseData = response.data;

      const firs: Array<Fir> = response.data.features.map((feature): Fir => {
        return {
          icao: feature.properties.id,
          points: feature.geometry.coordinates[0][0].map((point) => {
            return {
              lat: point[1],
              lng: point[0],
            };
          }),
          region: feature.properties.region,
        };
      });

      firs.forEach((fir) => this.firsMap.set(fir.icao, fir));

      const data = await readFile(`${__dirname}/firs_countries.json`, 'utf-8');
      const firData = JSON.parse(data);
      firData.forEach((fir) => {
        this.firCountry.set(fir.ICAO, fir.Country);
      });
    }
  }

  detectFir(point: Point): string {
    for (const fir of this.firsMap.values()) {
      if (this.geoService.isInside(fir.points, point)) {
        return fir.icao;
      }
    }

    throw new Error('Fir not founded');
  }

  detectCountryByFirCode(fir: string): string {
    return (this.firCountry.get(fir) || 'UNKNOWN').toUpperCase();
  }

  detectCountryByPoint(point: Point): string {
    for (const fir of this.firsMap.values()) {
      if (this.geoService.isInside(fir.points, point)) {
        return this.detectCountryByFirCode(fir.icao);
      }
    }

    throw new Error('Fir not founded');
  }

  isInsideFir(point: Point, fir: string) {
    if (!this.firsMap.has(fir)) {
      throw new Error();
    }

    return this.geoService.isInside(this.firsMap.get(fir).points, point);
  }
}
