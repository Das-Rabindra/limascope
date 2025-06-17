---
title: Container Groups
---

# Container Groups

Limascope performs automatic grouping of containers based on their stack name or service name. You can also create custom groups using labels.

## Default Groups

By default, containers are grouped by their stack name in host mode. If `com.docker.swarm.service.name` label is present, Limascope will automatically enable a "swarm mode" where all containers with the same service name will be joined together.

## Custom Groups

Additionally, you can create custom groups by adding a label to your container. The label is `dev.limascope.group` and the value is the name of the group. All containers with the same group name will be joined together in the UI. For example, if you have a group named `myapp`, all containers with the label `limascope.group=myapp` will be joined together.

Here is an example using Docker Compose or Docker CLI:

::: code-group

```sh
docker run --label dev.limascope.group=myapp hello-world
```

```yaml [docker-compose.yml]
services:
  limascope:
    image: hello-world
    labels:
      - dev.limascope.group=myapp
```

:::
