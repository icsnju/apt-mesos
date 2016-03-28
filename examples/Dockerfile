FROM dockerfile/nodejs

# dddd
RUN npm install -g bower gulp http-server

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . /usr/src/app

CMD ["bash"]