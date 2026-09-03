# Nox — Project Documentation
 
> An anonymous-first social platform for the Lagos event and nightlife community. Users can post as a ghost (fully anonymous) or as a visible identity, share DJ sets, discuss events, and engage with the local scene — without the social baggage of Instagram.
 
---

## Local processes

Use the root scripts from the repo root:

```bash
yarn dev:all
```

This starts:

- `web`
- `api`
- `notifications`
- `stories`
- `media`

If you only want the backend processes:

```bash
yarn dev:backend
```

Individual processes are also available:

```bash
yarn dev:api
yarn dev:notifications
yarn dev:stories
yarn dev:media
```

