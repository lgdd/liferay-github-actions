### STAGE 1: Build ###
FROM node:16 AS build
WORKDIR /usr/src/app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

### STAGE 2: Run ###
FROM nginx:latest
COPY --from=build /usr/src/app/dist/remote-apps /usr/share/nginx/html
WORKDIR /usr/share/nginx/html
RUN ls
RUN touch ./package.js && cat polyfills.*.js runtime.*.js main.*.js > package.js
RUN cp styles.*.css styles.css