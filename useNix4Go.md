
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

[nix website]: https://nixos.org/
[nix flakes wiki]: https://nixos.wiki/wiki/Flakes/
[nix flakes]: https://www.tweag.io/blog/2020-05-25-flakes/
[go language]: https://go.dev/
[summon website]: https://cyberark.github.io/summon/
