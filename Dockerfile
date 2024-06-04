FROM golang:1.20 as builder

WORKDIR /apps

COPY ./ /apps
RUN export GOPROXY=https://goproxy.cn \
    && go build  -ldflags "-s -w" -o notify \
    && chmod +x notify

FROM alpine
LABEL maintainer="tchua"
COPY --from=builder /apps/notify  /apps/
COPY --from=builder /apps/etc/  /apps/etc/
COPY --from=builder /apps/templates/  /apps/templates/

RUN echo -e  "http://mirrors.aliyun.com/alpine/latest-stable/main\nhttp://mirrors.aliyun.com/alpine/latest-stable/community" >  /etc/apk/repositories \
&& apk  update && apk --no-cache add tzdata gcompat libc6-compat \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata \
&& ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2

WORKDIR /apps

EXPOSE 18010

CMD ["./notify"]