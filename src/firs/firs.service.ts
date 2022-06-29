import { Injectable } from '@nestjs/common';
import { readFile } from 'fs/promises';
import { GeoLocator } from 'src/geo/geo.locator';

export class Point {
  lat: number;
  lng: number;
}

@Injectable()
export class FirService {
  private firsMap: Map<string, Array<Point>>;

  constructor(private geoService: GeoLocator) {
    this.firsMap = new Map();
  }

  async loadPointsByFIR(icao: string) {
    if (this.firsMap.has(icao)) {
      return;
    }

    const firs = await readFile(`${__dirname}/firs.json`, 'utf8').then((data) =>
      JSON.parse(data),
    );

    const feature = firs.features.find((f) => f.properties.ident === icao);

    const points = feature.geometry.coordinates[0].map((point) => {
      return {
        lat: point[1],
        lng: point[0],
      };
    });

    this.firsMap.set(icao, points);
  }

  isInsideFir(point: Point, fir: string) {
    if (!this.firsMap.has(fir)) {
      throw new Error();
    }

    return this.geoService.isInside(this.firsMap.get(fir), point);
  }
}
