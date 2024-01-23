FROM    golang:1.21.4-alpine3.18
LABEL   authors="pmaw"
WORKDIR /src
ADD     . /src
RUN     go install github.com/codegangsta/gin@latest
ENTRYPOINT ["gin"]