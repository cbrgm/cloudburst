FROM node:12-alpine

RUN npm install -g artillery@1.6.2 --allow-root --unsafe-perm=true

WORKDIR /workload

ENV PROCESSOR_DIR /workload

ADD ./processor.js ./processor.js
