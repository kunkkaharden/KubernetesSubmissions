FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /log-output .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /log-output /log-output
ENTRYPOINT ["/log-output"]
