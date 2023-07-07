#FOR SERVICE 
FROM golang:1.19-alpine AS debug

WORKDIR /services/order

COPY . .
 
RUN go build -o main cmd/order/main.go  

FROM alpine:3.17 AS prod

COPY --from=debug /services/order/main /
COPY --from=debug /services/order/configs/ /configs 
COPY --from=debug /services/order/templates/ /templates
COPY --from=debug /services/order/tmp/ /tmp
EXPOSE 8080

CMD /main
