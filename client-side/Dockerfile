FROM node:14-alpine AS builder

LABEL maintainer="Muhamad Hilmi Hibatullah <hilmihibatullah@gmail.com>"

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.21.0-alpine as production
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx","-g","daemon off;"]