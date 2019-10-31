# Sensu operator (forked from sensu)

[![CircleCI](https://circleci.com/gh/objectrocket/sensu-operator.svg?style=svg)](https://circleci.com/gh/objectrocket/sensu-operator)

Status: Proof of concept

The Sensu operator manages Sensu 2.0 clusters deployed to [Kubernetes][k8s-home] and automates tasks related to operating a Sensu cluster.

It is based on and heavily inspired by the [etcd-operator](https://github.com/coreos/etcd-operator).

## Setup instructions for ObjectRocket contributors

Start minikube:

```bash
minikube start --memory=3072 --kubernetes-version v1.13.1
```

Note: the operator has a NetworkPolicy capable CNI plugin, but not needed around here.

### Prerequisites

Build the binaries into a local image in minikube

```bash
#### Make sure the container image is build with the Minikube Docker
#### instance so that it's available for the kubelet later:
eval $(minikube docker-env)

#### Build the container:
make docker-build
```

Before proceding to the next steps make sure the image build is completed and shows:

```bash
docker image ls
REPOSITORY                                TAG                 IMAGE ID            CREATED             SIZE
objectrocket/sensu-operator               latest              6cc444763341        2 hours ago         60.2MB
```

### Installation

Create a role and role binding:

```bash
./example/rbac/create-role --namespace=sensu
```

Create a `sensu-operator` deployment:

```bash
kubectl apply -f example/deployment.yaml
```

You should end up with 1 running pods, e.g.:

```bash
$ kubectl -n sensu get pods -l name=sensu-operator
NAME                              READY     STATUS    RESTARTS   AGE
sensu-operator-6444f68845-54bvs   1/1       Running   0          1m
```

## Usage example

Create your first `SensuCluster` in the `sensu` namespace:

```bash
kubectl apply -f example/example-sensu-cluster-objectrocket.yaml
```

From within the cluster, the Sensu cluster agent should now be reachable
via:

```bash
ws://example-sensu-cluster-agent.sensu.svc.cluster.local:8081
```

If needed, login into the sensu backend UI:

```bash
$ kubectl -n sensu port-forward service/platdev0-dashboard 3000:3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
```

Open your browser at port 3000
User: `admin`
PW: `P@ssw0rd!`

To reach the Sensu cluster's services via `NodePort` do:

```bash
kubectl apply -f example/example-sensu-cluster-service-external.yaml

$ curl -Li http://$(minikube ip):31980/health
HTTP/1.1 200 OK
Date: Thu, 21 Jun 2018 14:44:47 GMT
Content-Length: 0
```

Finally create a sensu asset and check with the following command:

```bash
kubectl apply -f example/example-sensu-check-objectrocket.yaml
```

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
secret "sensu-backups-aws-secret" created
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

## Development

### OpenAPI specs

We use automatically generated OpenAPI specs for our custom resource definitions. If you make changes to the types, you can run this to regenerate specs:

```console
./hack/k8s/openapi-gen/openapi-gen.sh
```
