# run with docker

## Modes

You can run the backend in docker in 2 modes : 
- dev: deploy only the backend with hot reload for all the APIs + all the necessary services (db, redis, rabbit...)
- prod: deploy the complete platform (back and front) + all the necessary services (db, redis, rabbit...)

## Setup

1. Copy `.env.example` into `.env`.
2. Edit the user config section with the right values. 
3. Set the version for each component, so the latest version is deployed. (Updates are also handled that way)

## `prod`

1. To self-host, you can simply run the docker compose like this
```
docker compose -f docker/docker-compose.yaml up -d
```


## `dev`

**DO NOT USE IF YOU ARE NOT GOING TO EDIT THE BACKEND (DEVELOPERS ONLY)**

1. At the root folder of the `backend` repository:
```
docker compose -f docker/docker-compose-dev.yaml up -d
```