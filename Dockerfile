# Docker image for the Drone rsync plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w'
#     docker build --rm=true -t plugins/drone-rsync .

FROM alpine:3.6
ENV TIMEZONE Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIMEZONE /etc/localtime
RUN echo $TIMEZONE > /etc/timezone

RUN apk add -U ca-certificates openssh-client rsync && rm -rf /var/cache/apk/*
ADD drone-rsync /bin/
ENTRYPOINT ["/bin/drone-rsync"]
