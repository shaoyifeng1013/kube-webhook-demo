FROM base:v1

RUN mkdir -p /home/syf/server
ADD ./ /home/syf/server/

WORKDIR /home/syf/server

CMD ["./server", "webhook"]