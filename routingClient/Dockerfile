FROM gradle:8-alpine as build

COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle bootJar

FROM openjdk:17-alpine

EXPOSE 8081
RUN mkdir /app

COPY --from=build /home/gradle/src/build/libs/*.jar /app/route-client.jar

ENTRYPOINT ["java","-jar","/app/route-client.jar"]