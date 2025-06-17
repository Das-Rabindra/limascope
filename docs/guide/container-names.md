---
title: Container Names
---

# Container Names

By default, Limascope retrieves container names directly from Docker. This is usually sufficient, as these names can be customized using the `--name` flag in `docker run` commands or through the `container_name` field in Docker Compose services.

## Custom Names

In cases where modifying the container name itself isn't possible, you can override it by adding a `dev.limascope.name` label to your container.

Here is an example using Docker Compose or Docker CLI:

::: code-group

```sh
docker run --label dev.limascope.name=hello hello-world
```

```yaml [docker-compose.yml]
services:
  limascope:
    image: hello-world
    labels:
      - dev.limascope.name=hello
```

:::
