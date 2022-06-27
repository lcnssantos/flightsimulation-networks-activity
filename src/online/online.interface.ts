import { Activity } from './activity';

export interface OnlineService {
  getActivity(): Promise<Activity>;
}
