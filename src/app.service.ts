import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
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
    private repository: Repository<NetworksActivity>,
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

    await this.repository.save({ ivao, poscon, vatsim });
  }
}
