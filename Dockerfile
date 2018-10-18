FROM golang:1.9-alpine3.7

ENV SRC_DIR=/go/src/github.com/deFarro/letsdoit_backend/app
ENV PORT=9090

COPY ./app $SRC_DIR

WORKDIR $SRC_DIR

RUN apk add --no-cache git curl \
  && go get ./

RUN go build

HEALTHCHECK --interval=5s --timeout=1s --retries=3 \
  CMD curl --fail http://localhost:$PORT/version || exit 1

EXPOSE $PORT

CMD app