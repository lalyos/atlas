A cli tool for query atlas.hashicorp.com artifacts. It will print each matching atlas artifact version.

The [terraform atlas provider](https://www.terraform.io/docs/providers/atlas/r/artifact.html) can query
the atlas artifact versions, but i wanted to do it from a stadalone cli tool.

## Install

```
curl -L https://github.com/lalyos/atlas/releases/download/v0.0.5/atlas_0.0.5_$(uname)_x86_64.tgz | tar -xz -C /usr/local/bin/
```

## Usage

You can specify the artifact in 2 ways:
- specify 3 parameters: user, artifact_name, and artifact_type
- specify 1 parameter, the 3 above combined into one: user/artifact/type. Its called `slug`

```
atlas -s sequenceiq/cbd/amazon.image
```

```
$ atlas -u sequenceiq -a docker -t openstack.image

sequenceiq/docker/openstack.image/1
```

## Latest version only

By default all matching versions are listed, you can use the `-l` or `--last` options to list only the latest version.

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

Right now the following template functions are added:
 - json
 - add
 - subtract

## License

MIT

