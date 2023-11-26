FROM golang:alpine

RUN apk --no-cache add ca-certificates libc6-compat

EXPOSE 9091

COPY bin/pooling /

COPY bin/.env.conf /

# RUN chmod +x /pooling
 
ENTRYPOINT [ "/pooling" ]

