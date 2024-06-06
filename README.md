# Traefik RemoteAddr plugin

[![Build Status](https://github.com/RiskIdent/traefik-remoteaddr-plugin/workflows/Main/badge.svg?branch=main)](https://github.com/RiskIdent/traefik-remoteaddr-plugin/actions)

## Usage

This plugin is very simple: take the **client** IP and port and write them to some headers.
This is done by using the Go field [`net/http.Request.RemoteAddr`](https://pkg.go.dev/net/http#Request)
which is composed of `IP:port` of the client connection.

To mimic nginx's behaviour of `X-Forwarded-Port`, where it sets that header to the client's port, then use the dynamic middleware config:

```yaml
middlewares:
  my-middleware:
    plugin:
      remoteaddr:
        headers:
          port: X-Forwarded-Port
```

Alternatively, you could use the non-standard `X-Real-Port` to not override Traefik's behavior:

```yaml
middlewares:
  my-middleware:
    plugin:
      remoteaddr:
        headers:
          port: X-Real-Port
```

### Configuration

Traefik static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines a plugin:

<details open><summary>File (YAML)</summary>

```yaml
# Static configuration

experimental:
  plugins:
    remoteaddr:
      moduleName: github.com/RiskIdent/traefik-remoteaddr-plugin
      version: v0.1.0
```

</details>

<details><summary>CLI</summary>

```bash
# Static configuration

--experimental.plugins.remoteaddr.moduleName=github.com/RiskIdent/traefik-remoteaddr-plugin
--experimental.plugins.remoteaddr.version=v0.1.0
```

</details>

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

<details open><summary>File (YAML)</summary>

```yaml
# Dynamic configuration

http:
  middlewares:
    my-middleware:
      plugin:
        remoteaddr:
          headers:
            # if set, then set header "X-Real-Address" to the RemoteAddr (e.g "192.168.1.2:1234")
            address: X-Real-Address
            # if set, then set header "X-Real-Ip" to the IP of RemoteAddr (e.g "192.168.1.2")
            ip: X-Real-Ip
            # if set, then set header "X-Real-Port" to the port of RemoteAddr (e.g "1234")
            port: X-Real-Port
```

</details>

<details><summary>Kubernetes</summary>

```yaml
# Dynamic configuration

apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: my-middleware
spec:
  plugin:
    remoteaddr:
      headers:
        # if set, then set header "X-Real-Address" to the RemoteAddr (e.g "192.168.1.2:1234")
        address: X-Real-Address
        # if set, then set header "X-Real-Ip" to the IP of RemoteAddr (e.g "192.168.1.2")
        ip: X-Real-Ip
        # if set, then set header "X-Real-Port" to the port of RemoteAddr (e.g "1234")
        port: X-Real-Port
```

</details>

### Local Mode

Traefik also offers a developer mode that can be used for temporary testing of plugins not hosted on GitHub.
To use a plugin in local mode, the Traefik static configuration must define the module name (as is usual for Go packages) and a path to a [Go workspace](https://golang.org/doc/gopath_code.html#Workspaces), which can be the local GOPATH or any directory.

The plugins must be placed in `./plugins-local` directory,
which should be in the working directory of the process running the Traefik binary.
The source code of the plugin should be organized as follows:

```console
$ tree ./plugins-local/
./plugins-local/
    └── src
        └── github.com
            └── RiskIdent
                └── traefik-remoteaddr-plugin
                    ├── plugin.go
                    ├── plugin_test.go
                    ├── go.mod
                    ├── LICENSE
                    ├── Makefile
                    └── README.md
```

<details open><summary>File (YAML)</summary>

```yaml
# Static configuration

experimental:
  localPlugins:
    remoteaddr:
      moduleName: github.com/RiskIdent/traefik-remoteaddr-plugin
```

</details>

<details><summary>CLI</summary>

```bash
# Static configuration

--experimental.localPlugins.remoteaddr.moduleName=github.com/RiskIdent/traefik-remoteaddr-plugin
```

</details>

(In the above example, the `traefik-remoteaddr-plugin` plugin will be loaded from the path `./plugins-local/src/github.com/RiskIdent/traefik-remoteaddr-plugin`.)
