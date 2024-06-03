FROM golang AS builder

RUN go install -v github.com/goreleaser/goreleaser@latest

WORKDIR /home/

COPY . .
RUN GOVERSION=$(shell go version | awk '{print $$3;}') \
	goreleaser build --clean --snapshot --single-target

WORKDIR /app/
RUN tar -zxvf /home/dist/linux_amd64.tar.gz -C /app/

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/* /usr/share/zoneinfo/
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ /app/

CMD ["/app/{{ .PB.GoPackageName }}"]
