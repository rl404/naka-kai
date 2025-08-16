# Golang base image
FROM golang:1.25 AS go_builder
LABEL stage="naka-kai_builder"
WORKDIR /naka-kai
COPY . .
WORKDIR /naka-kai/cmd/naka-kai
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o naka-kai

# New stage from scratch
FROM jrottenberg/ffmpeg:5-alpine
LABEL stage="naka-kai_binary"
COPY --from=go_builder /naka-kai/cmd/naka-kai/naka-kai naka-kai
ENTRYPOINT []