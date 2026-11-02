FROM golang:alpine AS build
WORKDIR /
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN GOARCH=arm64 go build -tags lambda.norpc -ldflags='-s -w'

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y \
    chromium \
    ca-certificates \
    fonts-noto \
    fonts-noto-color-emoji \
    && rm -rf /var/lib/apt/lists/*
ENV CHROME_PATH=/usr/lib/chromium/chromium
ENV HOME=/tmp
ENV XDG_CONFIG_HOME=/tmp/.config
ENV XDG_CACHE_HOME=/tmp/.cache
ENV FONTCONFIG_PATH=/etc/fonts
ENV FONTCONFIG_FILE=/etc/fonts/fonts.conf
RUN mkdir -p /tmp/.config /tmp/.cache /tmp/chromedp /tmp/fontconfig
COPY --from=build /main /main
CMD ["/main"]
