FROM debian:bullseye-slim

COPY ./output/study_xxqg /opt/study_xxqg

RUN mkdir /opt/config/

COPY ./lib/config_default.yml /opt/config/config.yml

RUN apt-get -qq update \
        && apt-get -qq install -y --no-install-recommends ca-certificates curl \
    && apt-get update && apt-get install -y libx11-6 libgbm1 libasound2 libcairo2 libxshmfence1 libatspi2.0-0 libpango-1.0-0 libnss3  \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 libxrandr2 libxfixes3 libxdamage1 libxcomposite1 libxkbcommon0 \
    && chmod -R 777 /opt/study_xxqg && cd /opt/ &&  ./study_xxqg --init
EXPOSE 8080

VOLUME /opt/config

CMD cd /opt && ./study_xxqg