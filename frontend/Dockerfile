#Build from visual-simd-debugger/frontend

#To build:
#sudo docker build ./ -t visual-simd-front

#To run:
#sudo docker run -d --rm --name v-simd-front --network host visual-simd-front

FROM node:alpine

WORKDIR '/app'

COPY package.json .
RUN npm install
COPY . .
RUN \
yarn build && \
npm install -g serve
CMD ["serve", "-p", "3000", "-s", "build"]



# FROM node:alpine AS builder
# WORKDIR /app
# RUN npm install
# COPY . .
# RUN npm run build

# FROM node:alpine
# WORKDIR /app
# COPY --from=builder /app/build app/build
# RUN npm install -g serve
# EXPOSE 80
# CMD serve -p 80 -s build