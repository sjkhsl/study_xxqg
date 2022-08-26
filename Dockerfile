FROM ubuntu:jammy

ARG DEBIAN_FRONTEND=noninteractive
ARG TARGETARCH
ARG TZ="Asia/Shanghai"

RUN  apt-get -qq update && \
     apt-get -qq install -y --no-install-recommends tzdata ca-certificates libglib2.0-0 libnss3 libnspr4 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 \
         libdbus-1-3 libexpat1 libxcb1 libxkbcommon0 libx11-6 libxcomposite1 libxdamage1 libxext6 libxfixes3 libxrandr2 libgbm1 libpango-1.0-0 \
         libcairo2 libasound2 libatspi2.0-0 && \
     ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime && \
     echo ${TZ} > /etc/timezone && \
     dpkg-reconfigure --frontend noninteractive tzdata &&\
     rm -rf /var/lib/apt/lists/* && \
     mkdir /opt/config/

COPY conf/config_default.yml /opt/config/config.yml
COPY conf/QuestionBank.db /opt/QuestionBank.db

COPY ./dist/docker_linux_$TARGETARCH*/study_xxqg /opt/study_xxqg

RUN  chmod -R 777 /opt/study_xxqg && \
     cd /opt/ && \
     ./study_xxqg --init

EXPOSE 8080

VOLUME /opt/config

CMD cd /opt && ./study_xxqg