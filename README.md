# japanese-notes-parser

This is a tool create to parse notes on
https://github.com/alexhokl/notes/tree/master/spoken_languages/japanese.

## Commands

To create a database from a notes file

```sh
japanese-notes-parser create -f vocabulary.md -d vocab.db --overwrite
```

## Running with Docker

This repository ships a multi-stage `Dockerfile` that produces a static-CGO
binary running on a `scratch` base image as a non-root user (uid/gid `1000`),
with a writable volume at `/mnt/data`.

### Usage

Equivalent of the native command:

```sh
japanese-notes-parser create \
  -f spoken_languages/japanese/vocabulary.md \
  -d database/japanese_vocab.db \
  --overwrite
```

run inside the container:

```sh
docker run --rm \
  -v spoken_languages/japanese/vocabulary.md:/vocabulary.md:ro \
  -v database:/mnt/data \
  asia-east1-docker.pkg.dev/craffy/public/japanese-notes-parser create \
    -f /vocabulary.md \
    -d /mnt/data/japanese_vocab.db \
    --overwrite
```

Notes:
- The input notes file is bind-mounted read-only at `/vocabulary.md` inside the
  container.
- `database` is bind-mounted at `/mnt/data` (the image's declared writable
  volume), so the resulting `japanese_vocab.db` appears on the host at
  `database/japanese_vocab.db`.
- If your host directory is not owned by uid `1000`, add
  `--user $(id -u):$(id -g)` to the `docker run` command (see Caveats below).

### Caveats

1. **Host dir ownership.** When bind-mounting `./data:/mnt/data`, the host
   directory must be writable by uid `1000` (or whichever uid you run the
   container as). If your host directory is owned by a different uid, override
   the container user at run time, e.g. `--user $(id -u):$(id -g)`. The
   pre-created `/mnt/data` ownership only matters when relying on the image's
   anonymous volume; bind mounts take their ownership from the host.

2. **No shell in the image.** The `scratch` runtime has no shell, so
   `docker exec -it <container> sh` is not available. For interactive
   debugging, run a one-off `alpine` container mounting the same volumes, or
   add a separate debug target to the `Dockerfile`.

3. **UID choice.** The image defaults to uid/gid `1000:1000`. To use a
   different uid (e.g. `65532` to match distroless nonroot, or your host uid
   to avoid `--user` overrides at run time), rebuild with a different value.
   Run with `--user <uid>:<gid>` to override per-invocation without
   rebuilding.
