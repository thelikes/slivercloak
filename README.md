# slivercloak

PoC "framework" to build Sliver

* Two docker environments for building Sliver v1.5 and v1.6
* Go-based build program to clone, run modules, and compile protobufs and the Sliver client and server

## run

```bash
# Build the image for Sliver v1.5.42
docker build -f Dockerfile.1.5 -t cloak:1.5 .

# Build the image for Sliver master
docker build -f Dockerfile.1.6 -t cloak:1.6 .

# For sliver 1.5 (uses Go 1.18)
docker run -v $(pwd)/output:/tmp/output -it cloak:1.5 cloak -modules all

# For master (uses modern Go)
docker run -v $(pwd)/output:/tmp/output -it cloak:1.6 cloak -modules all
```

## modules
Example modules:

* [example module](./builder/mod-example.go)
* [branding module](./builder/mod-branding.go)
* [donutamsi module](./builder/mod-donut.go)