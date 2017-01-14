#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smscode .
FROM scratch
ADD smscode /
COPY etc/*   /etc/
CMD ["/smscode"]