# Redbase Plugin

Redbase Plugin is a middleware for [Traefik](https://github.com/traefik/traefik) that redirects according to records in an SQLite database served by a `Redbase` daemon [Redbase](https://hub.docker.com/vtacquet/redbase)

## Configuration

## Static

**traefik.yml**
```yaml
experimental:
  plugins:
    redbase:
      moduleName: github.com/vtacquet/redbase-plugin
      version: "v0.1.0"
```

## Dynamic

To configure the `Redbase` Plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example illustates
the usage of `Redbase` plugin middleware. The `redbaseurl` parameter is where the Redbase daemon is running. The `defaulturl` is 
the fallback url when there is no redirection entry in the database or when the `Redbase` daemon is not available. Since this is
a redirection middleware, the service can be set to `noop@internal`

**http.yml**
```yaml
http:
  routers:
    my-site:
      entrypoints: "in443
      tls:
        certresolver: "letsencrypt""
      rule: "Host(`my-server.com`) && PathPrefix(`/go`)"
      middlewares:
        - redbase
      service: noop@internal

  middlewares:
    redbase:
      plugin:
        redbase:
          redbaseurl: "redbase:9922"
          defaulturl: "https://www.tacquet.be"
```
