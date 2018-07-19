# Sensu operator

[![CircleCI](https://circleci.com/gh/kinvolk/sensu-operator.svg?style=svg)](https://circleci.com/gh/kinvolk/sensu-operator)

Status: Proof of concept

The Sensu operator manages Sensu 2.0 clusters deployed to [Kubernetes][k8s-home] and automates tasks related to operating a Sensu cluster.

It is based on and heavily inspired by the [etcd-operator](https://github.com/coreos/etcd-operator).

## Setup

Tested with a Kubernetes v1.10 Minikube setup:

```
$ minikube start --kubernetes-version v1.10.0
```

### Prerequisites

Build the binaries:

```
$ make build
```

Since there is no official, public `sensu-operator` container image
yet, i.e. you have to build your own:

```
#### Make sure the container image is build with the Minikube Docker
#### instance so that it's available for the kubelet later:
$ eval $(minikube docker-env)

#### Build the container:
$ make container
```

### Installation

Create a role and role binding:

```
$ ./example/rbac/create-role
```

Create a `sensu-operator` deployment:

```
$ kubectl apply -f example/deployment.yaml
```

You should end up with three running pods, e.g.:

```
$ kubectl get pods -l name=sensu-operator
NAME                              READY     STATUS    RESTARTS   AGE
sensu-operator-6444f68845-54bvs   1/1       Running   0          1m
sensu-operator-6444f68845-p74zn   1/1       Running   0          1m
sensu-operator-6444f68845-vpkxj   1/1       Running   0          1m
```

## Usage example

Create your first `SensuCluster`:

```
$ kubectl apply -f example/example-sensu-cluster.yaml
```

From within the cluster, the Sensu cluster agent should now be reachable
via:

```
ws://example-sensu-cluster-agent.default.svc.cluster.local:8081
```

To reach the Sensu cluster's services via `NodePort` do:

```
$ kubectl apply -f example/example-sensu-cluster-service-external.yaml

$ curl -Li http://$(minikube ip):31980/health
HTTP/1.1 200 OK
Date: Thu, 21 Jun 2018 14:44:47 GMT
Content-Length: 0
```

Let's deploy a dummy agent:

```
$ kubectl apply -f example/dummy-agent-deployment.yaml
```

The Sensu dashboard (via `http://192.168.99.100:31900/default/default/entities`)
should now show you two entities. `192.168.99.100` is the IP of the
Minikube instance and could be different on your system, see
`minikube ip`.

## Backup & restore

**Please note**: backup & restore is experimental and currently relies on
`etcdctl` being availabe in the `sensu/sensu` image. For that you
have to build your own.


#### Build sensu container with etcdctl

Fetch and checkout this branch:

https://github.com/kinvolk/sensu-go/tree/schu/dockerfile-add-etcdctl

Then run

```
eval $(minikube docker-env)
./build.sh build
./build.sh docker
```

### Setup

Sensu backup and restore operators can be set up to backup and
restore the state of a `SensuCluster` to and from S3.

Deploy the Sensu backup and restore operators:

```
$ kubectl apply -f example/backup-operator/deployment.yaml
$ kubectl apply -f example/restore-operator/deployment.yaml
```

Create a S3 bucket and an AWS IAM user with at least the following privileges:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:ListAllMyBuckets",
            "Resource": "arn:aws:s3:::*"
        },
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::YOUR_BUCKET",
                "arn:aws:s3:::YOUR_BUCKET/*"
            ]
        }
    ]
}
```

Create AWS S3 credentials like follows:

```
$ mkdir -p s3creds

$ cat <<EOF >s3creds/credentials
[default]
aws_access_key_id = YOUR_ACCESS_KEY_ID
aws_secret_access_key = YOUR_SECRES_ACCESS_KEY
EOF

$ cat <<EOF >s3creds/config
[default]
region = YOUR_BUCKET_REGION
EOF

$ kubectl create secret generic sensu-backups-aws-secret --from-file s3creds/credentials --from-file s3creds/config
```

### Backup

The `create-backup` helper script can be used to create backups:

```
$ ./example/backup-operator/create-backup --aws-bucket-name=YOUR_BUCKET --backup-name=sensu-cluster-backup-$(date +%s)
Backup of cluster 'example-sensu-cluster' with backup named 'sensu-cluster-backup-1529593491'
sensubackup.sensu.io "sensu-cluster-backup-1529593491" created
```

### Restore

To restore the state of a `SensuCluster`

* deploy a new clean `SensuCluster` and
* use the `restore-backup` helper script to restore a previously
  created backup.

For example:

```
$ kubectl apply -f example/example-sensu-cluster.yaml

$ ./example/restore-operator/restore-backup --cluster-name=new-sensu-cluster --aws-bucket-name=YOUR_BUCKET --backup-name=sensu-cluster-backup-1529593491
Restore of cluster 'example-sensu-cluster' with backup named 'sensu-cluster-backup-1529593491'
sensurestore.sensu.io "example-sensu-cluster" created
```

If everything went well, delete the `SensuRestore` resource, e.g.:

```
kubectl delete sensurestore example-sensu-cluster
```

## Testing

For example, to run the e2e tests (`PASSES="e2e"`):

```
$ minikube start --kubernetes-version v1.10.0
$ eval $(minikube docker-env)
$ make
$ ./example/rbac/create-role
$ KUBECONFIG=~/.kube/config \
  OPERATOR_IMAGE=kinvolk/sensu-operator:v0.0.1 \
  TEST_NAMESPACE=default \
  TEST_AWS_SECRET=sensu-backups-aws-secret \
  TEST_S3_BUCKET=sensu-backup-test \
  PASSES="e2e" \
  ./hack/test
```

[k8s-home]: http://kubernetes.io
