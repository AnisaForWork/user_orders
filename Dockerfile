#FOR SERVICE 
FROM golang:1.19-alpine AS gited

RUN apk fix && \
    apk --no-cache --update add git git-lfs gpg less openssh patch && \
    git lfs install

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

FROM gited AS debug

WORKDIR /services/order

COPY . /services/order
  
RUN go clean --modcache
 
RUN go build -o main cmd/order/main.go 

CMD [ "/go/bin/dlv", "--listen=:3305", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/services/order/main" ]

FROM alpine:3.17 AS prod 

COPY --from=debug /services/order/main /
COPY --from=debug /services/order/configs/ /configs  
EXPOSE 8080

CMD  /main
