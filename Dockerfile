#  # build env
#  # ---
#  FROM golang:alpine AS build-env
#  RUN apk add --update git
#
#  ADD . /work
#  WORKDIR /work
#  RUN make release
#
#
# toybox
FROM scratch

# COPY --from=build-env /work/toybox /
COPY toybox /

# create directories and sym links
RUN ["/toybox", "--install", "-s", "/"]
# move toybox itself
RUN ["/toybox", "mv", "/toybox", "/usr/sbin/toybox"]

ENV PATH "/usr/bin:/usr/sbin"

RUN ["mkdir", "-p", "/etc/ssl/certs"]
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

CMD ["/toybox", "shell"]
