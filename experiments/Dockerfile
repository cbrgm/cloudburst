FROM node:12-alpine

RUN npm i -g artillery@1 --allow-root --unsafe-perm=true

WORKDIR /workload

ENV PROCESSOR_DIR /workload

ADD ./processor.js ./processor.js