FROM alpine as build
RUN apk add go git
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download 
ADD main.go .
RUN go build -ldflags='-s -w' -o /main

FROM alpine
COPY --from=build /main /main
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /rie
RUN chmod 755 /rie
COPY entry.sh /
RUN chmod 755 /entry.sh
ENTRYPOINT [ "/entry.sh" ]
CMD ["/main"]