FROM openjdk:8-jdk-alpine

ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="Kotlin docker images built upon official openjdk alpine images" \
      org.label-schema.name="alpine-kotlin" \
      org.label-schema.schema-version="1.0.0-rc1" \
      org.label-schema.usage="https://github.com/Zenika/alpine-kotlin/blob/master/README.md" \
      org.label-schema.vcs-url="https://github.com/Zenika/alpine-kotlin" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor="Zenika" \
      org.label-schema.version="1.1-jdk8"

RUN apk add --no-cache bash && \
    apk add --no-cache -t build-dependencies wget && \
    cd /usr/lib && \
    wget https://github.com/JetBrains/kotlin/releases/download/v1.1.61/kotlin-compiler-1.1.61.zip && \
    unzip kotlin-compiler-*.zip && \
    rm kotlin-compiler-*.zip && \
    rm kotlinc/bin/*.bat && \
    apk del --no-cache build-dependencies

ENV PATH $PATH:/usr/lib/kotlinc/bin

CMD ["kotlinc"]
