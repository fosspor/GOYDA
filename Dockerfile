# syntax=docker/dockerfile:1

FROM node:20-alpine AS frontend
WORKDIR /fe
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
# Пустой URL → fetch на тот же origin (:8080)
ARG VITE_API_URL=
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
ENV GOFLAGS=-mod=mod
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy && go mod download
COPY --from=frontend /fe/dist /src/internal/spa/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -tags embed -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/server /app/server
COPY migrations /app/migrations
ENV HTTP_ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]
