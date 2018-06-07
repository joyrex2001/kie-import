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

## References

* [Your first Business Rules application on OpenShift: from Zero to Hero in 30 minutes](https://developers.redhat.com/blog/2017/06/13/your-first-business-rules-application-on-openshift-from-zero-to-hero-in-30-minutes/)
