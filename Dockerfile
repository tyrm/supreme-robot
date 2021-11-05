FROM golang:1.17 AS builder
RUN go get github.com/markbates/pkger/cmd/pkger

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN pkger && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o supreme-robot

FROM scratch
COPY --from=builder /etc/ssl /etc/ssl

COPY --from=builder /app/supreme-robot /supreme-robot
CMD ["/supreme-robot"]