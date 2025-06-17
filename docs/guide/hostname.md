---
title: Hostname
---

# Changing Limascope's Hostname

Limascope's default connection is called localhost. Using the `--hostname` flag, Limascope's name can be changed to anything. This value will be shown on the page title and under the Limascope logo.

Changing the label for localhost also changes the label for the `localhost` connection which is displayed under the multi-host menu. Below is an example of using `--hostname` to change the name of Limascope's subheader to `mywebsite.xyz`.

::: code-group

```sh
docker run --volume=/var/run/docker.sock:/var/run/docker.sock -p 8080:8080 Das-Rabindra/limascope --hostname mywebsite.xyz
```

```yaml [docker-compose.yml]
services:
  limascope:
    image: Das-Rabindra/limascope:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - 8080:8080
    environment:
      DOZZLE_HOSTNAME: mywebsite.xyz
```

:::
