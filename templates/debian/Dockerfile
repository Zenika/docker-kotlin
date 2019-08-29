FROM openjdk:#{.JDKVersion}-jdk

ARG CIRCLE_BUILD_DATE
ARG CIRCLE_SHA1

LABEL org.label-schema.build-date=$CIRCLE_BUILD_DATE \
      org.label-schema.description="Kotlin docker images built upon official openjdk images" \
      org.label-schema.name="kotlin" \
      org.label-schema.schema-version="1.0.0-rc1" \
      org.label-schema.usage="https://github.com/Zenika/docker-kotlin/blob/master/README.md" \
      org.label-schema.vcs-url="https://github.com/Zenika/docker-kotlin" \
      org.label-schema.vcs-ref=$CIRCLE_SHA1 \
      org.label-schema.vendor="Zenika" \
      org.label-schema.version="#{.Version}-jdk#{.JDKVersion}"

RUN cd /usr/lib && \
    wget -q #{.CompilerURL} && \
    unzip kotlin-compiler-*.zip && \
    rm kotlin-compiler-*.zip && \
    rm -f kotlinc/bin/*.bat

ENV PATH $PATH:/usr/lib/kotlinc/bin

CMD ["kotlinc"]
