FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /ping-pong .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /ping-pong /ping-pong
ENTRYPOINT ["/ping-pong"]
