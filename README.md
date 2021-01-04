# k8s-client-sandbox

## Module Lists

- [list-services](list-services)
- [list-nodes](list-nodes)
- [list-pods](list-pods)


## Develop modules
### Build from source
```
cd $GOPATH/src/github.com/
mkdir -p yokawasa
cd yokawasa/
git clone git@github.com:yokawasa/k8s-client-sandbox.git
cd k8s-client-sandbox
make
```

### Add new module

```
git clone git@github.com:yokawasa/k8s-client-sandbox.git
cd k8s-client-sandbox
mkdir new-module
vi main.go
go mod init new-module
go build
```
