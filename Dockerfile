FROM golang

RUN mkdir -p /app
WORKDIR /app
COPY bin/main .
EXPOSE 3000
CMD ["main"]