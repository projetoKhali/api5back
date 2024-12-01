FROM /prom/prometheus:v2.22.0 
COPY ./prometheus.yml /etc/prometheus/prometheus.yml
ENV TZ=America/Sao_Paulo
EXPOSE 9090

ENTRYPOINT [ "prometheus", \
    "--config.file=/etc/prometheus/prometheus.yml", \
    "--storage.tsdb.path=/prometheus", \
    "--storage.tsdb.retention=365d", \
    "--web.console.libraries=/usr/share/prometheus/console_libraries", \
    "--web.console.templates=/usr/share/prometheus/consoles", \
    "--web.external-url=http://localhost:9090", \
    "--log.level=info", \
    "--log.format=logger:stderr" ]