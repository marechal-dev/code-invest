FROM node:lts-alpine AS build

WORKDIR /home/app

COPY package.json .

RUN npm install

COPY . .

CMD ["npm", "run", "start:dev"]
