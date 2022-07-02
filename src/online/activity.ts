import { Column, Entity, ObjectIdColumn } from 'typeorm';

export class Activity {
  pilot: number;
  atc: number;
}

@Entity('activity')
export class NetworksActivity {
  @ObjectIdColumn()
  id: string;

  @Column()
  date: Date;

  @Column()
  ivao: Activity;

  @Column()
  vatsim: Activity;

  @Column()
  poscon: Activity;
}

@Entity('br_activity')
export class BrazilNetworksActivity {
  @ObjectIdColumn()
  id: string;

  @Column()
  date: Date;

  @Column()
  ivao: Activity;

  @Column()
  vatsim: Activity;

  @Column()
  poscon: Activity;
}

export class GeoActivity {
  [fir: string]: Activity;
}
@Entity('geo_activity')
export class GeoNetworksActivity {
  @ObjectIdColumn()
  id: string;

  @Column()
  date: Date;

  @Column()
  ivao: GeoActivity;

  @Column()
  vatsim: GeoActivity;

  @Column()
  poscon: GeoActivity;
}
