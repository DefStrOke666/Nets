# build stage
FROM node:lts-alpine as build-stage
ARG API_BASE_URL
ENV API_BASE_URL "$API_BASE_URL"
WORKDIR /frontend
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /frontend/dist /usr/share/nginx/html
COPY --from=build-stage /frontend/changeIP.sh /
EXPOSE 80
CMD ["sh", "changeIP.sh"]