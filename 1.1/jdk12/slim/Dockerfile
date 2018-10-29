FROM openjdk:12-jdk-slim

ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="Kotlin docker images built upon official openjdk images" \
      org.label-schema.name="kotlin" \
      org.label-schema.schema-version="1.0.0-rc1" \
      org.label-schema.usage="https://github.com/Zenika/docker-kotlin/blob/master/README.md" \
      org.label-schema.vcs-url="https://github.com/Zenika/docker-kotlin" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor="Zenika" \
      org.label-schema.version="1.1-jdk12-slim"

RUN apt-get update && apt-get install -y --no-install-recommends wget && \
    rm -rf /var/lib/apt/lists/* && \
    cd /usr/lib && \
    wget https://github.com/JetBrains/kotlin/releases/download/v1.1.61/kotlin-compiler-1.1.61.zip && \
    apt-get remove -y wget && \
    unzip kotlin-compiler-*.zip && \
    rm kotlin-compiler-*.zip && \
    rm -f kotlinc/bin/*.bat

ENV PATH $PATH:/usr/lib/kotlinc/bin

CMD ["kotlinc"]
