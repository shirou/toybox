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

RUN ["/toybox", "initialize_toybox", "-r", "/"]
RUN ["/toybox", "mv", "/toybox", "/usr/sbin/toybox"]

ENV PATH "/usr/bin:/usr/sbin"


CMD ["/toybox", "shell"]
