FROM node:latest as builder

WORKDIR /app
COPY . .
RUN npm install && \
    npm run build:tailwind && \
    npm run build


FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
