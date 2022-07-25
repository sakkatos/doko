# Golang TDD


### [Setup the doko nix development environment](#setup-doko-nix)

DOKO_ROOT will henceforth be the directory onto which you have cloned git@github.com:nerds-odd-e/doko.git
As non-root user, go into `$DOKO_ROOT` and execute `./setup-doko-env.sh`.
Exit the above shell when done and start a new shell.

- NB: Should you hit any issues installing nix for doko development, please refer to [official nix installation guide](https://nixos.org/download.html#download-nix)

### [Start doko nix development environment](#start-doko-nix)

Ensure your OS (WSL2/Ubuntu/Fedora, etc) has `/bin/sh` point to `bash`.
If you are using Ubuntu where `/bin/sh` is symlinked to `dash`, please run `sudo dpkg-reconfigure dash` and answer "No" to reconfigure to `bash` as default.
After successful [doko nix environment setup](#setup-doko-nix) from above, you may enter the doko nix environment with:
`nix develop`

### [Run TDD DB migrations](#db-migrations)

We are using [go-pg/pg](https://github.com/go-pg/pg) for orm and [go-pg/migrations](https://github.com/go-pg/migrations) for DB migrations.
A local .env file has been created when entering doko nix environment which sets up the needed environment variables.
Run `make migration` to perform the PostgreSQL DB migrations.

### [Run full Unit Test Suite](#run-test-suite)

After successful completion of [TDD DB migrations](#db-migrations), run the full unit test suite with `make test`,
