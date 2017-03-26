#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smscode .
FROM scratch
ADD smscode /
COPY conf/*   /conf/
ENV ZONEINFO  /conf/zoneinfo.zip
ENV PHONE_DATA_DIR  /conf/
CMD ["/smscode"]