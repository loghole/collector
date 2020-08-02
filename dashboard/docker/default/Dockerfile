# BUILD FRONTEND
FROM node:14.2.0-alpine as ui_builder

WORKDIR /usr/src/app

COPY ./frontend/package.json ./package.json
COPY ./frontend/yarn.lock ./yarn.lock

RUN yarn install

COPY ./frontend/ ./

RUN yarn build

# BUILD BACKEND
FROM golang:1.14.6-alpine as go_builder

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

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download && go mod verify

COPY ./backend ./

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

COPY --from=go_builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=go_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go_builder /etc/passwd /etc/passwd
COPY --from=go_builder /etc/group /etc/group
COPY --from=go_builder /app /app

COPY --from=ui_builder /usr/src/app/dist/ /static/

USER appuser:appuser

# Set frontend static dir
ENV FRONTEND_PATH=./static

CMD ["/app"]
