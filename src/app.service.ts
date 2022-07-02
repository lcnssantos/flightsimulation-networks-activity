import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { MongoRepository } from 'typeorm';
import {
  BrazilNetworksActivity,
  GeoNetworksActivity,
  NetworksActivity,
} from './online/activity';
import { IVAOOnline } from './online/ivao.online';
import { PosconOnline } from './online/poscon.online';
import { VatsimOnline } from './online/vatsim.online';

@Injectable()
export class AppService {
  constructor(
    private ivaoOnline: IVAOOnline,
    private vatsimOnline: VatsimOnline,
    private posconOnline: PosconOnline,
    @InjectRepository(NetworksActivity)
    private activityRepository: MongoRepository<NetworksActivity>,
    @InjectRepository(BrazilNetworksActivity)
    private brazilActivityRepository: MongoRepository<BrazilNetworksActivity>,
    @InjectRepository(GeoNetworksActivity)
    private geoActivityRepository: MongoRepository<GeoNetworksActivity>,
  ) {}

  async getActivity(): Promise<NetworksActivity> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivity(),
      this.vatsimOnline.getActivity(),
      this.posconOnline.getActivity(),
    ]);

    return this.activityRepository.create({ ivao, vatsim, poscon });
  }

  async saveActivity(): Promise<void> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivity(),
      this.vatsimOnline.getActivity(),
      this.posconOnline.getActivity(),
    ]);

    await this.activityRepository.save({
      ivao,
      poscon,
      vatsim,
      date: new Date(),
    });
  }

  async saveActivityBR(): Promise<void> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getBrazilActivity(),
      this.vatsimOnline.getBrazilActivity(),
      this.posconOnline.getBrazilActivity(),
    ]);

    await this.brazilActivityRepository.save({
      ivao,
      poscon,
      vatsim,
      date: new Date(),
    });
  }

  async saveActivityByRegion(): Promise<void> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivityByRegion(),
      this.vatsimOnline.getActivityByRegion(),
      this.posconOnline.getActivityByRegion(),
    ]);

    await this.geoActivityRepository.save({
      ivao,
      vatsim,
      poscon,
      date: new Date(),
    });
  }

  getHistoryByMinutes(minutes: number) {
    return this.activityRepository.find({
      where: { date: { $gt: new Date(Date.now() - minutes * 60 * 1000) } },
    } as any);
  }

  getBRHistoryByMinutes(minutes: number) {
    return this.brazilActivityRepository.find({
      where: { date: { $gt: new Date(Date.now() - minutes * 60 * 1000) } },
    } as any);
  }

  getGeoHistoryByMinutes(minutes: number) {
    return this.geoActivityRepository.find({
      where: { date: { $gt: new Date(Date.now() - minutes * 60 * 1000) } },
    } as any);
  }

  async getBrazilActivity() {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getBrazilActivity(),
      this.vatsimOnline.getBrazilActivity(),
      this.posconOnline.getBrazilActivity(),
    ]);

    return { ivao, vatsim, poscon };
  }

  async getGeoActivity() {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivityByRegion(),
      this.vatsimOnline.getActivityByRegion(),
      this.posconOnline.getActivityByRegion(),
    ]);

    return { ivao, vatsim, poscon };
  }
}
