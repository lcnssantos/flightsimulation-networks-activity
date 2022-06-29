import { Activity } from './activity';

export interface OnlineService {
  getActivity(): Promise<Activity>;
  getBrazilActivity(): Promise<Activity>;
}
