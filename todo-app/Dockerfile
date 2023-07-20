# base image
FROM node:latest as builder
WORKDIR /react-app
COPY package*.json .
RUN yarn install
COPY . .
RUN yarn run build

FROM nginx:stable-alpine as final
COPY ./nginx/default.conf /etc/nginx/conf.d
COPY --from=builder /react-app/build /usr/share/nginx/html
EXPOSE 80
STOPSIGNAL SIGTERM
CMD ["nginx", "-g", "daemon off;"]