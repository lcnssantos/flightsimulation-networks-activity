FROM keymetrics/pm2:14-alpine

WORKDIR /usr/src/app

COPY . ./
RUN yarn

RUN yarn build

RUN cp src/firs/firs.json dist/firs

CMD ["yarn", "run", "start:prod"]