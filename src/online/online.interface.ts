import { Activity, GeoActivity } from './activity';

export interface OnlineService {
  getActivity(): Promise<Activity>;
  getBrazilActivity(): Promise<Activity>;
  getActivityByRegion(): Promise<GeoActivity>;
}
