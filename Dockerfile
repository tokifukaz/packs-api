FROM golang:1.20-bullseye AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y curl
COPY ./ ./

RUN make build

FROM golang:1.20-bullseye AS prod
COPY --from=builder /app/packs-api /app/packs-api
CMD ["/app/packs-api"]