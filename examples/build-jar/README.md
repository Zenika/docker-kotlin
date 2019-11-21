# Building an app into a JAR file

Run `docker container run -v $(pwd):/app --rm zenika/kotlin kotlinc /app -include-runtime -d /app/hello.jar`
