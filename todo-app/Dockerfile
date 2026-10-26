FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /todo-app .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /todo-app /todo-app
ENTRYPOINT ["/todo-app"]
