FROM    golang:1.21.4-alpine3.18
LABEL   authors="pmaw"
WORKDIR /src
COPY     . /src
RUN     go build -o /main .
ENTRYPOINT ["/main"]