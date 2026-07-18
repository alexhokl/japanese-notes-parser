# syntax=docker/dockerfile:1

# ---- builder ----
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build \
      -ldflags="-s -w -extldflags=-static" \
      -o /out/japanese-notes-parser \
      . \
 && mkdir -p /out/mnt/data \
 && chown 1000:1000 /out/mnt/data \
 && echo "app:x:1000:1000::/home/app:/sbin/nologin" > /out/etc-passwd \
 && echo "app:x:1000:" > /out/etc-group

# ---- runtime (bare) ----
FROM scratch

COPY --from=builder /out/etc-passwd /etc/passwd
COPY --from=builder /out/etc-group /etc/group
COPY --from=builder --chown=1000:1000 /out/mnt/data /mnt/data
COPY --from=builder /out/japanese-notes-parser /japanese-notes-parser

USER 1000:1000
VOLUME ["/mnt/data"]
ENTRYPOINT ["/japanese-notes-parser"]