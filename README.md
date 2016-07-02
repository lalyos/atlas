A cli tool for query atlas.hashicorp.com artifacts. It will print each matching atlas artifact version.

## Install

```
curl -L https://github.com/lalyos/atlas/releases/download/v0.0.5/atlas_0.0.5_$(uname)_x86_64.tgz | tar -xz -C /usr/local/bin/
```

## Usage

```
$ atlas -u sequenceiq -a docker -t openstack.image

sequenceiq/docker/openstack.image/1
```
## Filtering

You can filter by any metadata.
```
atlas -u sequenceiq -a cloudbreak -t openstack.image -m cloudbreak_image_version=1.2.0-v1
```

## Custom format

You can use the `--format` or `-f` short option to define a custom template.
See [ArtifactVersion godoc](https://godoc.org/github.com/hashicorp/atlas-go/v1#ArtifactVersion) for available fields.
See [pkg/text/template godoc](https://golang.org/pkg/text/template/) to learn about golang template syntax.
```
atlas -u sequenceiq -a docker -t openstack.image -f '{{json .}}'
```

## License

MIT

