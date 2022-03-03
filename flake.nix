{
  description = "Provide a CyberArk Summon provider using Hashicorp Vault as secrets provider";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        revision = "${self.lastModifiedDate}-${self.shortRev or "dirty"}";
        pkgs = nixpkgs.legacyPackages.${system};

        vault4summon = pkgs.buildGoModule {
          pname = "vault4summon";
          version = revision;
          src = self;
          subPackages = [ "." ];
          vendorSha256 = "sha256-EaZ+orSes8a38gn8txIu6EWDX0hwTBNwbGfGjgxDCsU=";
        };
      in
      {
        # Nix develop
        devShell = pkgs.mkShellNoCC {
          name = "vault4summon-" + revision;
          buildInputs = [
            vault4summon
            pkgs.go
            pkgs.summon
          ];
        };

        # nix build, nix shell or nix run
        defaultPackage = pkgs.buildEnv {
          name = "vault4summon-" + revision;
          paths = [
            vault4summon
            pkgs.go
            pkgs.summon
          ];
        };

        # nix build .#oci
        packages = {
          oci =
            let
              port = "8000";
            in
            pkgs.dockerTools.buildLayeredImage
              {
                name = "vault4summon-test";
                tag = revision;

                contents = [
                  vault4summon
                  pkgs.summon
                ];

                extraCommands = ''
                  mkdir -p usr/local/lib/summon/
                  cp -ar ${vault4summon}/bin/vault4summon usr/local/lib/summon/vault4summon
                '';

                config = {
                  Cmd = [
                    "summon"
                    "--provider"
                    "vault4summon"
                    "--yaml"
                    "hello: !var secret/hello#foo"
                    "printenv"
                    "hello"
                  ];
                  ExposedPorts = {
                    "${port}/tcp" = { };
                  };
                };
              };
        };
      }
    );
}
