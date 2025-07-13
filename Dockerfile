FROM golang:1.24-alpine AS build

WORKDIR /app

COPY . .
RUN go build -o merge-sort .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/merge-sort .

CMD ["./merge-sort"]