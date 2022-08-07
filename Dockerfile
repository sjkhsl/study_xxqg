FROM debian:bullseye-slim

ARG TARGETARCH

RUN  ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
     echo 'Asia/Shanghai' >/etc/timezone && \
     apt-get -qq update && \
     apt-get -qq install -y --no-install-recommends ca-certificates curl && \
     apt-get install -y libx11-6 libgbm1 libasound2 libcairo2 libxshmfence1 libatspi2.0-0 libpango-1.0-0 libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 libxrandr2 libxfixes3 libxdamage1 libxcomposite1 libxkbcommon0 && \
     apt-get clean && \
     rm -rf /var/lib/apt/lists/* && \
     mkdir /opt/config/

COPY ./dist/docker_linux_$TARGETARCH*/study_xxqg /opt/study_xxqg

COPY conf/config_default.yml /opt/config/config.yml

RUN  chmod -R 777 /opt/study_xxqg && \
     cd /opt/ && \
     ./study_xxqg --init

EXPOSE 8080

VOLUME /opt/config

CMD cd /opt && ./study_xxqg