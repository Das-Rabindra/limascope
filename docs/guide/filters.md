---
title: Filter
---

# Filtering Containers

Limascope supports conditional filtering similar to Docker's [--filter](https://docs.docker.com/reference/cli/docker/container/ls/#filter) with `DOZZLE_FILTER` or `--filter`. Filters are passed directly to Docker to limit what Limascope can see. For example, filtering by label is supported with `--filter "label=color"`, which is similar to `docker ps` command with `docker ps --filter "label=color"`.

::: code-group

```sh
docker run --volume=/var/run/docker.sock:/var/run/docker.sock -p 8080:8080 Das-Rabindra/limascope --filter label=color
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
      DOZZLE_FILTER: label=color
```

:::

Common filters are `name` or `label` to limit Limascope's access to containers.

## UI, Agents, and User Filters <Badge type="tip" text="New" />

Limascope supports multiple filters to limit the containers it can see. Filters can be set at the UI, agent, or user level.

1. **UI Filters**: These filters are applied to the Limascope UI instance and sent to Docker to restrict the visible containers. They affect all agents and users who do not have their own filters.
2. **Agent Filters**: These filters are set at the agent level and sent to Docker to limit the containers exposed by that agent. Agent filters and UI filters work together to restrict the containers.
3. **User Filters**: These filters are set at the user level and determine which containers the user can see. If user filters are not defined, Limascope defaults to using the UI filters.

For more information on setting filters for specific users, see [user filters](/guide/authentication#setting-specific-filters-for-users). For details on setting filters for agents, see [agent filters](/guide/agent#setting-up-filters).

> [!WARNING]
> It is important to understand that multiple filters are combined to limit the containers. For example, if you set `--filter label=color` at the UI level and `--filter label=type` at the agent level, Limascope will only display containers that have both the `color` and `type` labels.
