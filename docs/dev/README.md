# Development

Prerequisites:
* [Set up local k8s](https://github.com/kubernetes/community/blob/master/contributors/devel/running-locally.md)


```
$ kubectl config use-context local
$ go mod vendor
$ docker build -f Dockerfile-local -t local-unload .
$ kubectl apply -f examples/local-deploy.yml
$ kubectl exec -it [pod name] /bin/bash
```
