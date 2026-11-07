FROM golang:1.26-alpine AS build

WORKDIR /src

COPY Backend/go.mod Backend/go.sum ./
RUN go mod download

COPY Backend/ ./
RUN go build -o /out/primarte ./cmd/api

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/primarte ./primarte
COPY Backend/configs    ./configs
COPY Backend/templates  ./templates
COPY L1-site            /L1-site

RUN mkdir -p uploads/avatars

EXPOSE 8080

CMD ["./primarte"]