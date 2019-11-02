FROM golang:1.11.4-alpine

ENV CGO_ENABLED=0
ENV DOCKER_RUN=true

RUN echo "Start" \
&& apk update && apk add gcc openssh-client openssh git \
&& git clone https://github.com/painkuter/mahjong.git \
&& cd mahjong \
&& git checkout actions \
&& go build cmd/single/single.go \
&& ls

CMD ["./mahjong/single"]