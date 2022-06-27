import { Column, Entity, ObjectIdColumn } from 'typeorm';

export class Activity {
  pilot: number;
  atc: number;
  time: Date;
}

@Entity('activity')
export class NetworksActivity {
  @ObjectIdColumn()
  id: string;

  @Column()
  ivao: Activity;

  @Column()
  vatsim: Activity;

  @Column()
  poscon: Activity;
}
