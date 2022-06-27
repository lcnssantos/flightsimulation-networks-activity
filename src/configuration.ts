import { plainToClass } from 'class-transformer';
import { IsNotEmpty, IsString, validate } from 'class-validator';
import { config } from 'dotenv';

config();

export class Configuration {
  @IsString()
  @IsNotEmpty()
  MONGO_URL: string;

  @IsString()
  @IsNotEmpty()
  MONGO_SSL: string;
}

export const validateConfiguration = async () => {
  const object = plainToClass(Configuration, process.env);
  const errors = await validate(object);
  if (errors.length > 0) {
    throw new Error(
      JSON.stringify(
        errors.map((error) => {
          const { target, ...rest } = error;
          return rest;
        }),
      ),
    );
  }

  return object;
};

export const getConfiguration = () => plainToClass(Configuration, process.env);
