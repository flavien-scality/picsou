FROM alpine:3.5

RUN apk add --no-cache ca-certificates

COPY picsou /usr/local/bin/picsou
COPY assets /usr/local/share/picsou/assets

CMD ["picsou", "report"]
