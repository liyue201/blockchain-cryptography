# 阶段一：构建前端（React + Vite + Monaco）
FROM node:20-alpine AS front
WORKDIR /app
COPY web/package.json ./
RUN npm install
COPY web/ .
RUN npm run build

# 阶段二：编译后端（需要阶段一产出的 web/dist 供 go:embed）
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY --from=front /app/dist web/dist
COPY embed.go ./
COPY content content
COPY examples examples
COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -o /out/bkcc-server ./cmd/server

# 阶段三：运行镜像
# 保留 Go toolchain，用于进程内沙箱(exec)编译用户提交的示例代码。
FROM golang:1.25-alpine
WORKDIR /app

COPY --from=build /out/bkcc-server /usr/local/bin/bkcc-server

# 以非 root 身份运行。
# 沙箱(exec)需要可写且可"执行"的工作区（注意：tmpfs/noexec 挂载点无法运行二进制），
# 因此统一放到 overlay 文件系统的 /workspace 下。
RUN adduser -D -H -u 10001 appuser \
 && mkdir -p /workspace/tmp /workspace/gocache /workspace/gomodcache \
 && chown -R appuser /workspace
ENV HOME=/workspace \
    TMPDIR=/workspace/tmp \
    GOCACHE=/workspace/gocache \
    GOMODCACHE=/workspace/gomodcache

USER appuser
EXPOSE 8080
ENTRYPOINT ["bkcc-server"]
CMD ["-addr", ":8080"]