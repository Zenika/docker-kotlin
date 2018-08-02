[![Docker Build Status](https://img.shields.io/docker/build/zenika/alpine-kotlin.svg)](https://hub.docker.com/r/zenika/alpine-kotlin/) [![Docker Pulls](https://img.shields.io/docker/pulls/zenika/alpine-kotlin.svg)](https://hub.docker.com/r/zenika/alpine-kotlin/)

### Supported tags and respective `Dockerfile` links

 * `1.2.60-jdk8`, `1.2.60`, `1.2-jdk8`, `1.2`, `1-jdk8`, `1`, `jdk8`, `latest` [(jdk8/1.2/Dockerfile)](https://github.com/Zenika/alpine-kotlin/blob/master/jdk8/1.2/Dockerfile)

 * `1.3-M1-jdk8`, `1.3-M1`, `1.3-jdk8`, `1.3` [(jdk8/1.3/Dockerfile)](https://github.com/Zenika/alpine-kotlin/blob/master/jdk8/1.3/Dockerfile)

 * `1.1.61-jdk8`, `1.1.61`, `1.1-jdk8`, `1.1` [(jdk8/1.1/Dockerfile)](https://github.com/Zenika/alpine-kotlin/blob/master/jdk8/1.1/Dockerfile)

### What is Kotlin

Kotlin is a statically-typed programming language that runs on the Java virtual machine and also can be compiled to JavaScript source code or use the LLVM compiler infrastructure. Its primary development is from a team of JetBrains programmers based in Saint Petersburg, Russia. While the syntax is not compatible with Java, Kotlin is designed to interoperate with Java code and is reliant on Java code from the existing Java Class Library, such as the collections framework.

See https://en.wikipedia.org/wiki/Kotlin_%28programming_language%29 for more information.

![Kotlin Logo](https://github.com/Zenika/alpine-kotlin/raw/master/Kotlin-logo.png)

### Usage

Start using the Kotlin REPL : `docker container run -it --rm zenika/alpine-kotlin`

See Kotlin compiler version : `docker container run -it --rm zenika/alpine-kotlin kotlinc -version`

See Kotlin compiler help : `docker container run -it --rm zenika/alpine-kotlin kotlinc -help`

### Reference

 * Kotlin website : https://kotlinlang.org

 * Where to file issues : https://github.com/Zenika/alpine-kotlin/issues

 * Maintained by : https://www.zenika.com
