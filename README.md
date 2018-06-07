# KIE import

KIE import is tool that will allow automation of creating KIE servers, created
with the Drools Workbench, in an OpenShift environment.

The Drools Workbench contains an internal git repository that can be cloned
and, with small adjustments on the maven pom file, can be compiled as a KIE
server.

One of the issues that comes with cloning and automating this repo, is that
it is based on username / password combinations, and uses an legacy key
algorithm. This will make the configuration, and automation a bit more
challenging. KIE import will make this easier.

## Config

KIE import can be used to clone this repo automatically based on the given
parameters (as environment variables):

* DROOLS_HOST
* DROOLS_GIT_SSH_PORT
* DROOLS_GIT_REPO
* GIT_USERNAME
* GIT_PASSWORD
* DESTINATION_FOLDER

## Using in OpenShift build

An OpenShift template, based on the ```decisionserver64-basic-s2``` template
is included. This template adds an intermediate imagestream which will clone
the code from the workbench repo first, before continuing with the actual
build.

Create the project with:

```bash
oc process -f openshift.yaml \
    -p APPLICATION_NAME="loan-demo" \
    -p KIE_SERVER_USER="brmsAdmin" \
    -p KIE_SERVER_PASSWORD="jbossbrms@01" \
    -p KIE_CONTAINER_DEPLOYMENT="container-loan10=com.redhat.demos:loandemo:1.0" \
    -p DROOLS_HOST="businesscentral" \
    -p DROOLS_GIT_REPO="loan" \
    -p DROOLS_PROJECT="loandemo" \
    -p GIT_USERNAME="erics" \
    -p GIT_PASSWORD='jbossbrms1!' | oc create -f -
```

And start the build:

```bash
oc start-build loan-demo --from-dir=.
```

## References

* [Micro-rules on OpenShift: The CoolStore just became even cooler!](https://developers.redhat.com/blog/2016/10/05/micro-rules-on-openshift-the-coolstore-just-became-even-cooler/)
* [Your first Business Rules application on OpenShift: from Zero to Hero in 30 minutes](https://developers.redhat.com/blog/2017/06/13/your-first-business-rules-application-on-openshift-from-zero-to-hero-in-30-minutes/)
