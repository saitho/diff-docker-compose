# diff-docker-compose

Inspired by [adamdicarlo/diff-docker-compose](https://github.com/adamdicarlo/diff-docker-compose) but in Go.

Small utility to diff two `docker-compose.yml` files; useful, for instance, if your
`docker-compose.yml` is based off of a template.

By default, running with no arguments will assume the two yaml files to diff are
`docker-compose.yml.template` and `docker-compose.yml`.

Currently, this only gives a high-level overview of how the files compare to each other.

![Screenshot showing output: services locally removed/disabled, services locally adedd/enabled, and
services locally modified](assets/screenshot1.png)

## Go API

### `lib.DiffYaml(oldYaml map[string]interface{}, newYaml map[string]interface{}) YamlDiffResult`

DiffYaml compares two map[string]interface{} objects (e.g. read via YAML parser).

Arguments:
* `oldYaml`: map structure of the YAML file considered the old or current state
* `newYaml`: map structure of the YAML file considered the new or desired state

Returns:
* `YamlDiffResult` object

### `YamlDiffResult`

This is the result object returned by the `DiffYaml` method.

Properties:

* `Diffs []YamlDiffEntry`: list of all nodes with differences.
* `Structure map[string]*YamlDiffStructure`: tree structure of both YAML files combined

Methods:
* `HasChanged(path []string) bool`: returns true when there differences for the node
* `Get(path []string) []YamlDiffEntry`: returns differences for the node identified by the selector (e.g. `Get([]string{"services"})` or `Get([]string{"services", "myapp"})`)
* `GetStructure(path []string) *YamlDiffStructure`: returns the structure for the node identified by the path. Access the diff for that node or its children for diff information.

*Note:* When a node is added, the whole node will be listed as difference in `YamlDiffEntry` (obtained via `Get()`).
If you need to traverse the node's children and look for changes, use `GetStructure`.

Examples:

```go
result := lib.DiffYaml(composeTemplate, composeActual)
fmt.Println(result.Get([]string{"services"}))

// this will get all differences of the "service" node
fmt.Println(result.GetStructure([]string{"services"}).GetDiff())

// both lines below will get only the differences of the "app2" service
fmt.Println(result.GetStructure([]string{"services", "app2"}).GetDiff())
result.GetStructure([]string{"services"}).GetChildren()["app2"].GetDiff()
```

```
[{[services app3] <nil> map[image:myapp:latest]} {[services db-service] <nil> map[image:db:latest]} {[services user-service environment] [SECRET=SET_ME] [SECRET=MYSECRET123]} {[services app2] map[image:myapp:latest] <nil>}]
{[services] map[app2:map[image:myapp:latest] user-service:map[environment:[SECRET=SET_ME] image:user-service:latest] web:map[build:. ports:[8000:5000]]] map[app3:map[image:myapp:latest] db-service:map[image:db:latest] user-service:map[environment:[SECRET=MYSECRET123] image:user-service:latest] web:map[build:. ports:[8000:5000]]]}
{[services app2] map[image:myapp:latest] <nil>}
{[services app2] map[image:myapp:latest] <nil>}
```
