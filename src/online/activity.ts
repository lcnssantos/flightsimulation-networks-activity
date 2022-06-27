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
