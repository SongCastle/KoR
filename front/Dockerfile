FROM node:17.3-slim

WORKDIR /opt/front
COPY ./package.json ./package-lock.json* ./

RUN apt-get update && \
    apt-get install -y curl iproute2 procps && \
    mkdir -p /opt/front && \
    chown -R node:node /opt/front

USER node
RUN npm cache clean --force && npm install

CMD ["npm", "start"]
