FROM node:9.2

RUN mkdir -p /app
WORKDIR /app

COPY package.json yarn.lock /app/
RUN yarn install --no-progress && yarn cache clean
