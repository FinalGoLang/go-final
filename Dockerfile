FROM golang:1.18-alpine AS builder
WORKDIR /source
COPY . /source
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o film-subscription ./cmd/.

FROM alpine:3.9
RUN mkdir /app
WORKDIR /app
RUN mkdir config && mkdir schema
COPY --from=builder /source/config /app/config
COPY --from=builder /source/schema /app/schema
COPY --from=builder /source/film-subscription /usr/local/bin
RUN chmod a+x /usr/local/bin/film-subscription


ENTRYPOINT [ "film-subscription" ]
