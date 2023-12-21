# example-web-app

## Prerequisite

- [docker](https://docs.docker.com/engine/install/)
  - with [compose](https://docs.docker.com/compose/install/linux/)
- [go-task/task](https://github.com/go-task/task)

## Development

- start services (with hot reloading):

  - prepare 2 shells
    - shell 1: run `task watch`
    - shell 2: run `task up`

- tidy packages:

  - run `task tidy`

- directly connect to postgres db:

  - run `task db`

## Deployment

- run `task deploy`
