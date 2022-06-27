FROM keymetrics/pm2:14-alpine

WORKDIR /usr/src/app

COPY . ./
RUN yarn

RUN yarn build

CMD ["yarn", "run", "start:prod"]