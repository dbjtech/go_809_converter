FROM registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:go_ent AS build

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/release
# 添加项目文件
ADD . .
# 复制go.sum 避免再下载一次依赖
COPY --from=registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:go_ent /go/src/go.mod /go/release/.
COPY --from=registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:go_ent /go/src/go.sum /go/release/.

RUN chmod 755 ./mark_info.sh && ./mark_info.sh
RUN GOOS=linux GOARCH=amd64 go build -mod=mod  -ldflags '-s -w' -o converter.exe ./converter

FROM registry.cn-qingdao.aliyuncs.com/from_docker_hub/alpine:latest

WORKDIR project

COPY --from=build /go/release/converter.exe /project
COPY --from=build /go/release/converter/static /project/converter/static
COPY --from=build /go/release/git_info.txt /project
COPY --from=build /go/release/version.json /project
COPY --from=build /go/release/config/configuration.toml.template /project/config/configuration.toml

ENV TZ=Asia/Shanghai
RUN echo "http://mirrors.aliyun.com/alpine/v3.19/main/" > /etc/apk/repositories \
    && apk --no-cache add tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone
