FROM scratch

ARG coherence_server_address

ENV COHERENCE_SERVER_ADDRESS=127.0.0.1:1408

COPY runner /files/runner

ENTRYPOINT ["/files/runner"]
CMD ["-h"]


