# Step 1: Build the Go binary
FROM golang:1.22.5-alpine AS build

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY .. .

RUN go build -o volley-live-score ./cmd/matches

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/volley-live-score /app/volley-live-score

CMD ["./volley-live-score"]

EXPOSE 8080
