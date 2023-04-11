A Kubernetes operator built with [operator-builder](https://github.com/nukleros/operator-builder).

## Quick Start

Deploy the operator (be sure to substitute the appropriate version) from a release:

    kubectl apply -f https://github.com/scottd018/ocm-log-forwarder-operator/releases/download/${VERSION}/deploy.yaml


## Local Development & Testing

To install the custom resource/s for this operator, make sure you have a
kubeconfig set up for a test cluster, then run:

    make install

To run the controller locally against a test cluster:

    make run

You can then test the operator by creating the sample manifest/s:

    kubectl apply -f config/samples

To clean up:

    make uninstall


## Deploy the Controller Manager (from Source)

First, set the image (be sure the substitute the appropriate version):

    export IMG=ghcr.io/scottd018/ocm-log-forwarder-operator:latest

Now you can build and push the image:

    make docker-build
    make docker-push

Then deploy:

    make deploy

To clean up:

    make undeploy


## Companion CLI

To build the companion CLI:

    make build-cli

The CLI binary will get saved to the bin directory.  You can see the help
message with:

    ./bin/ocmlogctl help


## Deploy the Operator Lifecycle Manager Bundle

First, build the bundle.  The bundle contains metadata that makes it 
compatible with Operator Lifecycle Manager and also makes the operator 
importable into OpenShift OperatorHub:

    make bundle

Next, set the bundle image.  This is the image that contains the packaged 
bundle:

    export BUNDLE_IMG=ghcr.io/scottd018/ocm-log-forwarder-operator-bundle:latest

Now you can build and push the bundle image:

    make bundle-build
    make bundle-push

To deploy the bundle (requires OLM to be running in the cluster):

    make operator-sdk
    bin/operator-sdk bundle validate $BUNDLE_IMG
    bin/operator-sdk run bundle $BUNDLE_IMG

To clean up:

    bin/operator-sdk cleanup --delete-all $BUNDLE_IMG
