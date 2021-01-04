

```
kubens default
kubectl run -it --rm=true busybox --image=yauritux/busybox-curl --restart=Never
./list-pods
kubectl delete pod busybox
./list-pods
```
