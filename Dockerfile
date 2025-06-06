FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /qa-api

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian11:latest
COPY --from=build /qa-api /

CMD ["/qa-api"]
