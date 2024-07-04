# builder
FROM node:18.18.0-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN yarn install --frozen-lock-file

COPY . .

RUN yarn build

# runner
FROM nginx:alpine

WORKDIR /usr/local/bin

COPY --from=builder /app/dist /usr/share/nginx/html

COPY .docker_files/webapp.conf /etc/nginx/conf.d/

EXPOSE 80

ENTRYPOINT [ "nginx", "-g", "daemon off;"]