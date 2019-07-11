# Development

Prerequisites:
* [Set up local k8s](https://github.com/kubernetes/community/blob/master/contributors/devel/running-locally.md)


Running `unload` locally:
```
$ kubectl config use-context local
$ go mod vendor
$ docker build -f Dockerfile-local -t local-unload .
$ kubectl apply -f examples/local-deploy.yml
$ kubectl apply -f examples/fortune-teller.yaml
```

Add `teller.local` to local resolver:
```
$ vi /etc/hosts
127.0.0.1   teller.local
```

`grpcurl` fortune-teller via `unload`:
```
git clone --depth=1 https://github.com/kubernetes/ingress-nginx.git
grpcurl -v \
        -plaintext \
        -proto ingress-nginx/images/grpc-fortune-teller/proto/fortune/fortune.proto \
        teller.local:50051 build.stack.fortune.FortuneTeller/Predict
```