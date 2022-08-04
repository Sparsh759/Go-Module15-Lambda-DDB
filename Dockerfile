FROM public.ecr.aws/lambda/provided:al2 as build-env

RUN yum install -y golang

ENV GO111MODULE=on
ENV GOSUMDB=off
ENV GOPROXY=direct

WORKDIR /go/src
ADD / ./
RUN ls -ltr

RUN go mod download
RUN go build -o cmd

FROM public.ecr.aws/lambda/provided:al2

COPY --from=build-env /go/src/cmd /opt

WORKDIR /opt
ENTRYPOINT /opt/cmd