FROM docker.io/golang:1.10

ARG CODE=github.com/joyrex2001/kie-import

ENV GIT_DESTINATION=/import/${DROOLS_GIT_REPO}
ENV SSH_KNOWN_HOSTS=/tmp/known_hosts

ENV DROOLS_HOST=${DROOLS_HOST}
ENV DROOLS_GIT_SSH_PORT=${DROOLS_GIT_SSH_PORT}
ENV DROOLS_GIT_REPO=${DROOLS_GIT_REPO}
ENV GIT_USERNAME=${GIT_USERNAME}
ENV GIT_PASSWORD=${GIT_PASSWORD}

ADD . /go/src/${CODE}/
RUN cd /go/src/${CODE}       && \
    go build -o /app/main    && \
    touch ${SSH_KNOWN_HOSTS} && \
    /app/main                && \
    ./hack/add_dep.pl ${GIT_DESTINATION}/${DROOLS_PROJECT}/pom.xml ${DEP_KIE_API_VERSION}

CMD ["/app/main"]
