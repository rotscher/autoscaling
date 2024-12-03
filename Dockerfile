FROM golang:1.23 AS build

WORKDIR /src
COPY autoscaling/* /src/
RUN go build -o /usr/local/bin/autoscaling ./main.go

FROM debian:latest
LABEL authors="rotscher"
RUN useradd -u 1001 nonroot

COPY --from=0 /usr/local/bin/autoscaling /usr/local/bin/autoscaling

USER nonroot
EXPOSE 3333

ENTRYPOINT ["/usr/local/bin/autoscaling"]