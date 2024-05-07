FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
COPY cmd/templates /templates

CMD [ "/app/mailApp" ]