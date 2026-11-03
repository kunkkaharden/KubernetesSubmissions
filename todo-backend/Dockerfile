FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /todo-backend .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /todo-backend /todo-backend
ENTRYPOINT ["/todo-backend"]
