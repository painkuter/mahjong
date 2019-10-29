FROM golang:1.11.4-alpine

#ENV GOROOT=/usr/local/go
#ENV GOPATH=/go
#ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
ENV CGO_ENABLED=0

RUN echo "Start" \
#&& which service \
&& apk update && apk add gcc openssh-client openssh git \
#&& service rsyslog start \
#&& git clone https://github.com/mmcgrana/gobyexample \
#&& cd gobyexample \
#&& go build examples/hello-world/hello-world.go \
&& git clone https://github.com/painkuter/mahjong.git \
&& cd mahjong \
&& git checkout actions \
&& go build cmd/single/single.go \
&& ls
#&& cd .. \
#&& service syslog start
#&& mkdir dev \
#&& mkdir log \




#COPY service_key /root/.ssh/id_rsa
#RUN chmod 600 /root/.ssh/id_rsa

#RUN apk update && apk add --no-cache openssh-client openssh bash git coreutils make \
#graphviz font-bitstream-type1 ttf-freefont ca-certificates \
#&& go get -u github.com/golang/dep/cmd/dep \
#&& go get -u -insecure golang.org/x/lint/golint
#
#RUN ssh-keyscan -H gitlab.ozon.ru >> /root/.ssh/known_hosts

EXPOSE 80:8080
CMD ["./mahjong/single"]
#CMD ["./gobyexample/hello-world"]