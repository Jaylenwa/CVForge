ARG TARGETPLATFORM
ARG TARGETARCH

FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
ARG VITE_OAUTH_ALLOWED_ORIGINS=
ENV VITE_OAUTH_ALLOWED_ORIGINS=$VITE_OAUTH_ALLOWED_ORIGINS
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
