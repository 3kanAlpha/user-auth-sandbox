FROM node:18.7.0-alpine3.15

CMD ["/bin/sh"]

WORKDIR /app

RUN apk update && \
    npm install -g @vue/cli

ADD ./ui/package.json /app

RUN yarn install