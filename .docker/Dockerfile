ARG DEBIAN_IMAGE="debian:stable-slim"
FROM ${DEBIAN_IMAGE}

RUN useradd -Ms /bin/bash hr
COPY target/bin/* /usr/local/bin/

USER hr

CMD [ "hr" ]
