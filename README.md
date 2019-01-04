# Sensu operator (forked from sensu)

[![CircleCI](https://circleci.com/gh/objectrocket/sensu-operator.svg?style=svg)](https://circleci.com/gh/objectrocket/sensu-operator)

Status: Proof of concept

The Sensu operator manages Sensu 2.0 clusters deployed to [Kubernetes][k8s-home] and automates tasks related to operating a Sensu cluster.

It is based on and heavily inspired by the [etcd-operator](https://github.com/coreos/etcd-operator).

## Setup

Start Minikube with CNI plugins enabled and install Calico for network policies to take effect:

```bash
minikube start --memory=3072 --kubernetes-version v1.10.0 --extra-config=controller-manager.cluster-cidr=192.168.0.0/16 --extra-config=controller-manager.allocate-node-cidrs=true --network-plugin=cni --extra-config=kubelet.network-plugin=cni
kubectl apply -f https://docs.projectcalico.org/v3.1/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
kubectl apply -f https://docs.projectcalico.org/v3.1/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
```

Network policies will get installed automatically with a Sensu cluster.

For testing, a NetworkPolicy capable CNI plugin is not necessary, the operator will install the policy regardless without effect.

```bash
minikube start --memory=3072 --kubernetes-version v1.10.0
```

### Prerequisites

Build the binaries:

```bash
make build
```

Since there is no official, public `sensu-operator` container image
yet, i.e. you have to build your own:

```bash
#### Make sure the container image is build with the Minikube Docker
#### instance so that it's available for the kubelet later:
eval $(minikube docker-env)

#### Build the container:
make container
```

### Installation

Create a role and role binding:

```bash
./example/rbac/create-role
```

Create a `sensu-operator` deployment:

```bash
kubectl apply -f example/deployment.yaml
```

You should end up with three running pods, e.g.:

```bash
$ kubectl get pods -l name=sensu-operator
NAME                              READY     STATUS    RESTARTS   AGE
sensu-operator-6444f68845-54bvs   1/1       Running   0          1m
sensu-operator-6444f68845-p74zn   1/1       Running   0          1m
sensu-operator-6444f68845-vpkxj   1/1       Running   0          1m
```

## Usage example

Create your first `SensuCluster`:

```bash
kubectl apply -f example/example-sensu-cluster.yaml
```

From within the cluster, the Sensu cluster agent should now be reachable
via:

```bash
ws://example-sensu-cluster-agent.default.svc.cluster.local:8081
```

To reach the Sensu cluster's services via `NodePort` do:

```bash
kubectl apply -f example/example-sensu-cluster-service-external.yaml

$ curl -Li http://$(minikube ip):31980/health
HTTP/1.1 200 OK
Date: Thu, 21 Jun 2018 14:44:47 GMT
Content-Length: 0
```

Let's deploy a dummy agent:

```bash
kubectl apply -f example/dummy-agent-deployment.yaml
```

The Sensu dashboard (via `http://192.168.99.100:31900/default/default/entities`)
should now show you two entities. `192.168.99.100` is the IP of the
Minikube instance and could be different on your system, see
`minikube ip`.

## Backup & restore

### Setup backup/restore operators

Sensu backup and restore operators can be set up to backup and
restore the state of a `SensuCluster` to and from S3.

Deploy the Sensu backup and restore operators:

```bash
kubectl apply -f example/backup-operator/deployment.yaml
kubectl apply -f example/restore-operator/deployment.yaml
```

Create a S3 bucket and an AWS IAM user with at least the following privileges:

```bash
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

```bash
mkdir -p s3creds

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

```bash
$ ./example/backup-operator/create-backup --aws-bucket-name=YOUR_BUCKET --backup-name=sensu-cluster-backup-$(date +%s)
Backup of cluster 'example-sensu-cluster' with backup named 'sensu-cluster-backup-1529593491'
sensubackup.objectrocket.com "sensu-cluster-backup-1529593491" created
```

### Restore

To restore the state of a `SensuCluster`

* deploy a new clean `SensuCluster` and
* use the `restore-backup` helper script to restore a previously
  created backup.

For example:

```bash
kubectl apply -f example/example-sensu-cluster.yaml

$ ./example/restore-operator/restore-backup --cluster-name=example-sensu-cluster --aws-bucket-name=YOUR_BUCKET --backup-name=sensu-cluster-backup-1529593491
Restore of cluster 'example-sensu-cluster' with backup named 'sensu-cluster-backup-1529593491'
sensurestore.objectrocket.com "example-sensu-cluster" created
```

If everything went well, delete the `SensuRestore` resource, e.g.:

```bash
kubectl delete sensurestore example-sensu-cluster
```

## Testing

For example, to run the e2e tests (`PASSES="e2e"`):

```bash
minikube start --kubernetes-version v1.10.0
eval $(minikube docker-env)
make
./example/rbac/create-role
KUBECONFIG=~/.kube/config \
  OPERATOR_IMAGE=objectrocket/sensu-operator:v0.0.1 \
  TEST_NAMESPACE=default \
  TEST_AWS_SECRET=sensu-backups-aws-secret \
  TEST_S3_BUCKET=sensu-backup-test \
  PASSES="e2e" \
  ./hack/test
```

[k8s-home]: http://kubernetes.io
