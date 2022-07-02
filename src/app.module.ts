import { HttpModule } from '@nestjs/axios';
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { getConfiguration } from './configuration';
import { FirService } from './firs/firs.service';
import { GeoLocator } from './geo/geo.locator';
import {
  BrazilNetworksActivity,
  GeoNetworksActivity,
  NetworksActivity,
} from './online/activity';
import { IVAOOnline } from './online/ivao.online';
import { PosconOnline } from './online/poscon.online';
import { VatsimOnline } from './online/vatsim.online';

@Module({
  imports: [
    HttpModule.register({}),
    TypeOrmModule.forRoot({
      type: 'mongodb',
      url: getConfiguration().MONGO_URL,
      entities: [NetworksActivity, BrazilNetworksActivity, GeoNetworksActivity],
      logger: 'simple-console',
      logging: true,
      ssl: getConfiguration().MONGO_SSL === 'true',
      useUnifiedTopology: true,
      useNewUrlParser: true,
    }),
    TypeOrmModule.forFeature([
      NetworksActivity,
      BrazilNetworksActivity,
      GeoNetworksActivity,
    ]),
  ],
  controllers: [AppController],
  providers: [
    GeoLocator,
    FirService,
    PosconOnline,
    VatsimOnline,
    IVAOOnline,
    AppService,
  ],
})
export class AppModule {}
