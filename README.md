# diff-docker-compose

Inspired by [adamdicarlo/diff-docker-compose](https://github.com/adamdicarlo/diff-docker-compose) but in Go.

Small utility to diff two `docker-compose.yml` files; useful, for instance, if your
`docker-compose.yml` is based off of a template.

By default, running with no arguments will assume the two yaml files to diff are
`docker-compose.yml.template` and `docker-compose.yml`.

Currently, this only gives a high-level overview of how the files compare to each other.

![Screenshot showing output: services locally removed/disabled, services locally adedd/enabled, and
services locally modified](assets/screenshot1.png)
