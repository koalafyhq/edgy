# edgy

The top layer of edgyfn.app that powering every website on Koalafy's Dedicated Hosting

## Notes

This project is still under heavy development and experimentally run koalafy-edge.edgyfn.app
landing page.

## Concept

I will just draw some simple diagram for now:

![https://s3.edgyfn.app/koalafy/misc/Untitled-2020-05-10-1200.png](https://s3.edgyfn.app/koalafy/misc/Untitled-2020-05-10-1200.png)

Here is some points:

1. User visit koalafy-edge.edgyfn.app
2. If the entrypoint is not exists, send 404 not found
3. Otherwise, we'll get "endpoint" that reflecting the "deployment id" on our s3 storage
4. Every request URI will have this format: `{endpoint}{path}`, e.g: `ke4furi8dh2pgpx/index.html`
5. We'll get the content from our redis db by querying `cache:{endpoint}{path}`
6. If we don't have any cache yet, request to the "origin server", e.g: https://s3.edgyfn.app/bundles/ke4furi8dh2pgpx/index.html
7. Save the cache into our redis db using this format: `cache:{endpoint}{path}`
8. If we serve user from cache, `X-Edgy-Cache` will have `HIT` value
9. Otherwise, it will return `MISS`
10. Client-side cache is done via `Etag` value (computed by minio)
11. We only cache <= 10MB assets
12. `S3_GATEWAY` will determine which "origin server" should talk to
13. TBD

## Development

First of all, you need to clone this repo, obviously.

Next, you need to have Go installed on your machine, mine is `v1.14.2`.

To make our life easier, you need to have `make(1)` on your machine,
if you won't install it, you can take a look into `Makefile` file.

If all is good, run `make`. The distributed binary is live under `dist` directory.

This server is only run on HTTPS protocol since we'll use HTTP/2. To generate the certificate
for development purpose, you need to install [`mkcert`](https://github.com/FiloSottile/mkcert) on your machine.

And then run `make certs`, the certificate is live under `certs` directory.

To run the server—and to make our life easier, hopefully—run some required services (redis & minio)
we use redis for caching-thing and minio (a s3-compatible object storage) as the "origin server" to run it,
execute this command:

```bash
$ docker-compose -f docker/docker-compose.dev.yml up # or run it in detached mode with -d
```

I believe you've installed docker on your machine, right? ;)

Also, I've add some data persistence-thing to make ~~our~~ my life easier, by running command above
the `Warehouse` directory will be created on your home directory. Don't worry, you can delete it
anytime you want.

To access minio instance, you can access it on https://localhost:3001 with `minioadmin` as access & secret
key

To access redis instance, you can access it on default redis uri string (localhost:6379).

Now, let's run the server. To make your life easier (again), just paste command below on your terminal.

```bash
$ REDIS_URI=localhost:6379 EDGY_REGION=dev S3_GATEWAY=localhost:3001 ./dist/edgy
```

Now you can access the server on https://localhost:3000

## Deployment

Just ask [@faultable](https://twitter.com/faultable) for now!
