# BUILD BACKEND
FROM golang:1.14.6-alpine as builder

RUN apk add --no-cache ca-certificates tzdata git

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src

COPY go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG APP_NAME
ARG SERVICE_NAME
ARG CMD_NAME
ARG VERSION
ARG GIT_HASH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
    -X $APP_NAME/config.ServiceName=$SERVICE_NAME \
    -X $APP_NAME/config.AppName=$APP_NAME \
    -X $APP_NAME/config.GitHash=$GIT_HASH \
    -X $APP_NAME/config.Version=$VERSION \
    -X $APP_NAME/config.BuildAt=$(date --utc +%FT%TZ) \
    " -o /app cmd/$SERVICE_NAME/*.go

# BUILD FINAL CONTAINER
FROM scratch as final
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app /app

USER appuser:appuser

ENTRYPOINT ["/app"]
