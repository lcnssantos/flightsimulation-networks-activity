import { HttpModule } from '@nestjs/axios';
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { getConfiguration } from './configuration';
import { NetworksActivity } from './online/activity';
import { IVAOOnline } from './online/ivao.online';
import { PosconOnline } from './online/poscon.online';
import { VatsimOnline } from './online/vatsim.online';

@Module({
  imports: [
    HttpModule.register({}),
    TypeOrmModule.forRoot({
      type: 'mongodb',
      url: getConfiguration().MONGO_URL,
      entities: [NetworksActivity],
      logger: 'simple-console',
      logging: true,
      ssl: getConfiguration().MONGO_SSL === 'true',
      useUnifiedTopology: true,
      useNewUrlParser: true,
    }),
    TypeOrmModule.forFeature([NetworksActivity]),
  ],
  controllers: [AppController],
  providers: [PosconOnline, VatsimOnline, IVAOOnline, AppService],
})
export class AppModule {}
