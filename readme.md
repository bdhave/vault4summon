# vault4summon: this project implements a [CyberArk Summon][summon website] provider using [Hashicorp Vault][Hashicorp Vault] as secrets provider.
>This project is distributed under the terms of the European Union Public Licence (EUPL) Version 1.2 or newer.`
You can find the latest version of the EUPL licence [here][EUPL].

> The work on this software project is in no way associated with my employer nor with the role I'm having at my employer. Any requests for changes will be decided upon exclusively by myself based on my personal preferences. I maintain this project as much or as little as my spare time permits.

>**WARNING: the current code was only tested with a Hashicorp Vault Server in development mode.
It must therefore be considered as a Proof of Concept, and it is not intended to be used in production
until [issue 1](https://github.com/bdhave/vault4summon/issues/1#issue-798122084) is  closed**

## Quick start

* Get vault4summon
  * download it
  * or build it
  * or use [Nix](useNix4Go.md) to create your Go environment to build it
* Install
  * Install [Summon][summon website] if you don't have it already installed
  * Copy `vault4summon` to `/usr/local/lib/summon/`
* Configure
  * Set the [environment variables](https://www.vaultproject.io/docs/commands#environment-variables)
    to access Hashicorp Vault
    * `VAULT_ADDR`: e.g. http://127.0.0.1:8200/.
    * `VAULT_TOKEN`: e.g. 00000000-0000-0000-0000-000000000000
* Create a [secrets.yml](secrets.yml) file
* Use [Summon][summon website]
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
  by [AWS Secrets Manager provider][AWS-summon]
* secret/name/mysecretkeypath as used
  by [Keepass kdbx database file provider][Keepass] or [Gopass provider][Gopass]

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
use [cla-assistant.io][cla]. You will be prompted to sign once a pull request is opened.

[go language]: https://go.dev/
[summon website]: https://cyberark.github.io/summon/
[Hashicorp Vault]: https://www.vaultproject.io/
[AWS-summon]: https://github.com/cyberark/summon-aws-secrets
[Keepass]: https://github.com/mskarbek/summon-keepass
[Gopass]: https://github.com/gopasspw/gopass-summon-provider

[EUPL]: https://ec.europa.eu/isa2/solutions/european-union-public-licence-eupl_en/
[cla]: https://cla-assistant.io
