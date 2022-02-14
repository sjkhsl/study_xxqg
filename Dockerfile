FROM golang:bullseye AS build_study_xxqg

ARG TARGETARCH
ARG BUILDARCH

ENV CGO_ENABLED="0"
ENV GOOS="linux"
ENV GOPROXY="https://goproxy.cn"
ARG UPX_VERSION="3.96"

RUN  sed -i "s@http://deb.debian.org@https://mirrors.163.com@g" /etc/apt/sources.list && apt update && apt install -y xz-utils git gcc curl
COPY / /study_xxqg
RUN cd /study_xxqg && \
    version=$(git describe --tags --long --always) && \
    echo "${version}" && \
    GOARCH=${TARGETARCH} go build -ldflags "-s -w -X main.version=${version}" -trimpath -o study_xxqg
RUN target="${BUILDARCH}" && \
    curl -sSL https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${target}_linux.tar.xz | tar xvJf - -C / && \
    cp -f /upx-${UPX_VERSION}-${target}_linux/upx /usr/bin/ && \
    /usr/bin/upx -9 -v /study_xxqg/study_xxqg


FROM debian:bullseye-slim

ARG TARGETARCH

ENV TIMEZONE="Asia/Shanghai"

RUN apt update \
    && apt install -y --no-install-recommends ca-certificates tzdata libglib2.0-0 libnss3 libatk1.0-0 libcups2 libatk-bridge2.0-0 libdrm2 libxcb1 libxkbcommon0 libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libgbm1 libpango-1.0-0 libcairo2 libasound2 libxshmfence1 \
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
# install study_xxqg
COPY --from=build_study_xxqg /study_xxqg/study_xxqg /study_xxqg

VOLUME ["/config"]

ENTRYPOINT ["/study_xxqg"]
