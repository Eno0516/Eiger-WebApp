FROM node:20-alpine AS builder

WORKDIR /app

# Copy workspace root
COPY package*.json ./
# Copy packages and this app only (avoid copying other apps)
COPY packages/ packages/
COPY apps/user/ apps/user/

RUN npm ci --workspace=apps/user
RUN npm run build --workspace=apps/user

FROM nginx:alpine

COPY --from=builder /app/apps/user/dist /usr/share/nginx/html
COPY apps/user/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
