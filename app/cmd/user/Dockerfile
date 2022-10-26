FROM golang:1.18.2 as builder

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN useradd -u 10001 scratchuser

WORKDIR /advapi

COPY go.mod go.sum ./
COPY app app
COPY db db

RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o=/go/bin ./app/cmd/user/main.go

EXPOSE 8080
CMD ["/go/bin/main"]


FROM scratch

COPY --from=builder /go/bin/main .
COPY --from=builder /etc/passwd /etc/passwd

USER scratchuser
EXPOSE 8080
CMD ["./main"]
