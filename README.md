# Emplidler

This is very simple open source cross platform desktop client to track user idle time and send metrics to backend 

### Support platform:
- windows 7\10
- linux (debian based)
- mac os x (10.13+)

## Noties:

- This is work in progress application
- It is created for some personal needs 
- Server will be published soon

## How to build

*Linux*
```
docker pull alexrett/emplidler-linux-build:0.0.1
docker run --rm -v $(pwd):/opt/emplidler alexrett/emplidler-linux-build:0.0.1 bash -c "cd /opt/emplidler/ && go build ."
```