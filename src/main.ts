import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { validateConfiguration } from './configuration';

async function bootstrap() {
  await validateConfiguration();
  const app = await NestFactory.create(AppModule);
  app.enableCors();
  await app.listen(process.env.PORT || 80);
}
bootstrap();
