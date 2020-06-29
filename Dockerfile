FROM scratch

COPY bin/resque_exporter /

CMD ["/resque_exporter"]