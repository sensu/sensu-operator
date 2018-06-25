## How to use codegen

Dependencies:
- Docker

In repo root dir, run:

```sh
./hack/k8s/codegen/update-generated.sh
```

It should print:

```
Generating deepcopy funcs
Generating clientset for sensu:v1beta2 at github.com/kinvolk/sensu-operator/pkg/generated/clientset
Generating listers for sensu:v1beta2 at github.com/kinvolk/sensu-operator/pkg/generated/listers
Generating informers for sensu:v1beta2 at github.com/kinvolk/sensu-operator/pkg/generated/informers
```
