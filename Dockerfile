#FOR SERVICE 
FROM golang:1.19-alpine AS gited

RUN apk fix && \
    apk --no-cache --update add git git-lfs gpg less openssh patch && \
    git lfs install

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

FROM gited AS debug

WORKDIR /services/user

COPY . /services/user

ARG GITHUB_PAT

RUN git config --global url."https://${GITHUB_PAT}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go clean --modcache
 
RUN go build -o main cmd/user/main.go cmd/user/closable.go

CMD [ "/go/bin/dlv", "--listen=:3305", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/services/user/main" ]

FROM alpine:3.17 AS prod

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

COPY --from=debug /services/user/main /
COPY --from=debug /services/user/configs/ /configs 
COPY --from=debug /services/user/server.crt /
COPY --from=debug /services/user/server.key /
EXPOSE 8080

CMD /wait && /main
