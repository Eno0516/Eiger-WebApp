FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
COPY packages/ packages/
COPY apps/admin/ apps/admin/

RUN npm ci --workspace=apps/admin
RUN npm run build --workspace=apps/admin

FROM nginx:alpine

COPY --from=builder /app/apps/admin/dist /usr/share/nginx/html
COPY apps/admin/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
