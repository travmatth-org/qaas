FROM travmatth/amazonlinux-golang-test:latest

RUN mkdir /usr/qaas \
    && mkdir /etc/qaas

WORKDIR /usr/qaas

COPY . .

RUN make build.test.all \
    && unzip -o dist/assets.zip -d /srv \
    && mv configs/httpd.yml /etc/qaas \
    && mv dist/httpd /usr/sbin