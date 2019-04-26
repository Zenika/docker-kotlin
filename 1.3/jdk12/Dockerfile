FROM openjdk:12-jdk

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
      org.label-schema.version="1.3-jdk12"

RUN yum install -y wget unzip && \
    cd /usr/lib && \
    wget https://github.com/JetBrains/kotlin/releases/download/v1.3.31/kotlin-compiler-1.3.31.zip && \
    unzip kotlin-compiler-*.zip && \
    rm kotlin-compiler-*.zip && \
    rm -f kotlinc/bin/*.bat && \
    yum remove -y wget unzip

ENV PATH $PATH:/usr/lib/kotlinc/bin

CMD ["kotlinc"]
