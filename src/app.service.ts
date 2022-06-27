import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { MongoRepository } from 'typeorm';
import { NetworksActivity } from './online/activity';
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
    private repository: MongoRepository<NetworksActivity>,
  ) {}

  async getActivity(): Promise<NetworksActivity> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivity(),
      this.vatsimOnline.getActivity(),
      this.posconOnline.getActivity(),
    ]);

    return this.repository.create({ ivao, vatsim, poscon });
  }

  async saveActivity(): Promise<void> {
    const [ivao, vatsim, poscon] = await Promise.all([
      this.ivaoOnline.getActivity(),
      this.vatsimOnline.getActivity(),
      this.posconOnline.getActivity(),
    ]);

    await this.repository.save({ ivao, poscon, vatsim, date: new Date() });
  }

  getHistory() {
    return this.repository.find({
      where: { date: { $gt: new Date(Date.now() - 24 * 60 * 60 * 1000) } },
    } as any);
  }
}
