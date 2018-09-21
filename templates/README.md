[![Docker Build Status](https://img.shields.io/docker/build/zenika/kotlin.svg)](https://hub.docker.com/r/zenika/kotlin/) [![Docker Pulls](https://img.shields.io/docker/pulls/zenika/kotlin.svg)](https://hub.docker.com/r/zenika/kotlin/)

### Supported tags and respective `Dockerfile` links
#{range $_, $version := .Versions}
#### #{$version.Version}
#{range $_, $jdkVersion := $version.JDKVersions}
 * `#{$version.Version}-jdk#{$jdkVersion.JDKVersion}`#{range $_, $additionalTag := $jdkVersion.Base.AdditionalTags}, `#{$additionalTag}`#{end} [(#{$version.Version}/jdk#{$jdkVersion.JDKVersion}/Dockerfile)](https://github.com/Zenika/docker-kotlin/blob/master/#{$version.Version}/jdk#{$jdkVersion.JDKVersion}/Dockerfile)
#{range $_, $variant := $jdkVersion.Variants}
 * `#{$version.Version}-jdk#{$jdkVersion.JDKVersion}-#{$variant.Base}`#{range $_, $additionalTag := $variant.AdditionalTags}, `#{$additionalTag}`#{end} [(#{$version.Version}/jdk#{$jdkVersion.JDKVersion}/#{$variant.Base}/Dockerfile)](https://github.com/Zenika/docker-kotlin/blob/master/#{$version.Version}/jdk#{$jdkVersion.JDKVersion}/#{$variant.Base}/Dockerfile)
#{end}#{end}#{end}
### What is Kotlin

Kotlin is a statically-typed programming language that runs on the Java virtual machine and also can be compiled to JavaScript source code or use the LLVM compiler infrastructure. Its primary development is from a team of JetBrains programmers based in Saint Petersburg, Russia. While the syntax is not compatible with Java, Kotlin is designed to interoperate with Java code and is reliant on Java code from the existing Java Class Library, such as the collections framework.

See https://en.wikipedia.org/wiki/Kotlin_%28programming_language%29 for more information.

![Kotlin Logo](https://github.com/Zenika/docker-kotlin/raw/master/Kotlin-logo.png)

### Usage

Start using the Kotlin REPL : `docker container run -it --rm zenika/kotlin`

See Kotlin compiler version : `docker container run -it --rm zenika/kotlin kotlinc -version`

See Kotlin compiler help : `docker container run -it --rm zenika/kotlin kotlinc -help`

### Reference

 * Kotlin website : https://kotlinlang.org

 * Where to file issues : https://github.com/Zenika/docker-kotlin/issues

 * Maintained by : https://www.zenika.com
