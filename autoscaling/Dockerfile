FROM debian:latest
LABEL authors="rotscher"

ADD autoscaling /usr/local/bin/autoscaling
RUN chmod 755 /usr/local/bin/autoscaling

ENTRYPOINT ["autoscaling"]