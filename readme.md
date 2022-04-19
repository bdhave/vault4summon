$# vault4summon: this project implement a [CyberArk Summon][summon website] provider using [Hashicorp Vault][Hashicorp Vault] as secrets provider.
>This project is distributed under the terms of the European Union Public Licence (EUPL) Version 1.2 or newer.`
You can find the latest version of the EUPL licence [here][EUPL].

> The work on this software project is in no way associated with my employer nor with the role I'm having at my employer. Any requests for changes will be decided upon exclusively by myself based on my personal preferences. I maintain this project as much or as little as my spare time permits.

**WARNING: the current code was only tested with a Hashicorp Vault Server in development mode.
**It must therefore be considered as a Proof of Concept, and it is not intended to be used in production
until [issue 1](https://github.com/bdhave/vault4summon/issues/1#issue-798122084) is  closed**

## Quick start

* Build or download vault4summon
  * `go build`
* Install
  * Install [Summon][summon website] if you don't hzve it already
  * Copy `vault4summon` to `/usr/local/lib/summon/`
* Configure
  * Set the [environment variables](https://www.vaultproject.io/docs/commands#environment-variables)
    to access Hashicorp Vault
    * `VAULT_ADDR`: e.g. http://127.0.0.1:8200/.
    * `VAULT_TOKEN`: e.g. 00000000-0000-0000-0000-000000000000
* Use Summon

### Using Nix

[Nix][nix website] is a tool that takes a unique approach to package
management and system configuration.

[Nix Flakes][nix flakes wiki] are an upcoming feature of the Nix package manager.

[Flakes][nix flakes] allow to define inputs (*you can think of them as dependencies*) and outputs of packages in a declarative way.

You will notice similarities to what you find in package definitions for other languages and like many language package managers flakes also introduce dependency pinning using a lockfile (`flake.lock`).

If you're willing to contribute and develop in `vault4summon`, a Flake file is shipped within this project.

The development environment provides the following tools:

* [go][go language]
* [summon][summon website]

To enter in a development environment, run:

```bash
nix develop
```

or

```bash
nix shell
```

If you just want to just build `vault4summon` locally, run:

```bash
nix build
```
##### note
Flakes and commands are still experimental features, they will be piushed to production level in 2022. So you still need a ~/.config/nix/nix.conf file with at least this line:
```
experimental-features = nix-command flakes
```
## Summon provider contract

Providers for Summon are easy to write. Given the identifier of a secret, they either return its value or an error.

There is the contract:

* They take one and only one argument, the identifier of a secret (a string). The argument can also be a flag with value
  -v or --version. The provider must return his version on stdout.

* If retrieval is successful, they return the value on stdout with exit code 0.

* If an error occurs, they return an error message on stderr with a non-0 exit code.

* The default path for providers is /usr/local/lib/summon/. If one provider is in that path, summon will use it. If
  multiple providers are in the path, you can specify which one to use with the --provider flag, or the environment
  variable SUMMON_PROVIDER. If your providers are placed outside the default path, give summon the full path to them.

* Variable IDs are used as identifiers for fetching Secrets. These are made up of a secret name (required) and secret
  key path (optional).

The Vault CLI to retrieve a secret is

`vault kv get -field=mysecretkeypath secret/name`

This provider has 2 implemented formats for Variable ID:

* secret/name#mysecretkeypath as used
  by [AWS Secrets Manager provider](https://github.com/cyberark/summon-aws-secrets)
* secret/name/mysecretkeypath as used
  by [Keepass kdbx database file provider](https://github.com/mskarbek/summon-keepass)

So the two commands below return the same value

`
summon --provider vault4summon --yaml 'hello: !var secret/name#mysecretkeypath' printenv hello
`

`
summon --provider vault4summon --yaml 'hello: !var secret/name/mysecretkeypath' printenv hello
`

## Contributing guidelines

If you would like to contribute code to vault4summon you can do so through GitHub by forking the repository and sending
a pull request.

When submitting code, please make efforts to follow existing conventions and style in order to keep the code as readable
as possible. Please also make sure your code compiles and passes tests.

Before your code can be accepted into the project, you must also sign the Individual Contributor License Agreement. I
use [cla-assistant.io](https://cla-assistant.io). You will be prompted to sign once a pull request is opened.

[nix website]: https://nixos.org/
[nix flakes wiki]: https://nixos.wiki/wiki/Flakes/
[nix flakes]: https://www.tweag.io/blog/2020-05-25-flakes/
[go language]: https://go.dev/
[summon website]: https://cyberark.github.io/summon/
[Hashicorp Vault]: https://www.vaultproject.io/
[EUPL]: https://ec.europa.eu/isa2/solutions/european-union-public-licence-eupl_en/
