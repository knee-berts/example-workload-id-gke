FROM golang:1.16 as build
WORKDIR /root/src/app
ADD . /root/src/app

RUN go get -d -v
RUN go build -o /root/bin/app

FROM gcr.io/distroless/base-debian10
COPY --from=build /root/bin/app /
CMD ["/app"]
