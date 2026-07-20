{
  description = "forge — unified CLI and library for interacting with git forges";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forEachSystem (pkgs:
        let
          mkForge = { version, rev, hash, vendorHash }:
            pkgs.buildGoModule {
              pname = "forge";
              inherit version vendorHash;
              src = pkgs.fetchFromGitHub {
                owner = "git-pkgs";
                repo = "forge";
                inherit rev hash;
              };

              go = pkgs.go_1_26 or pkgs.go;

              nativeBuildInputs = [ pkgs.installShellFiles ];

              postInstall = ''
                installShellCompletion --cmd forge \
                  --bash <($out/bin/forge completion bash) \
                  --fish <($out/bin/forge completion fish) \
                  --zsh <($out/bin/forge completion zsh)
              '';

              subPackages = [ "cmd/forge" ];

              meta = with pkgs.lib; {
                description = "Unified CLI for GitHub, GitLab, Gitea, Forgejo and Bitbucket";
                homepage = "https://github.com/git-pkgs/forge";
                license = licenses.mit;
                mainProgram = "forge";
              };
            };
        in
        rec {
          default = latest;

          v0_6_0 = mkForge {
            version = "0.6.0";
            rev = "v0.6.0";
            hash = "sha256-kVKDHcrtXbOqqZoiKb/SxOKbTy2A7oHomlUImkcnxmA=";
            vendorHash = "sha256-sduEepxhOCLk7/YMJbIwtt78Bo9UJ5olb8po7drxPZw=";
          };

          latest = mkForge {
            version = "0.6.0-unstable-2026-07-20";
            rev = "b3eb7489772377d24e9687a95f0d5bfdd5027d2d";
            hash = "sha256-ztb2KyDGs+Lkt6gvc1BRTGeSlM8HZq0zlJAjrBZMsiw=";
            vendorHash = "sha256-nG/p5oJI5oJP82z0GvEc+qBTS/5X5lokk2WRUkmwvBk=";
          };
        });

      devShells = forEachSystem (pkgs: {
        default = pkgs.mkShell {
          inputsFrom = [ self.packages.${pkgs.system}.default ];
          packages = with pkgs; [
            gopls
            act
            (writeShellScriptBin "forge-install-profile" ''
              set -e
              staged_flake=0
              if ! git ls-files --error-unmatch flake.nix &>/dev/null || \
                 git diff --name-only | grep -q '^flake\.nix$'; then
                git add flake.nix
                staged_flake=1
              fi
              if nix profile list | grep -q 'forge'; then
                nix profile remove forge
              fi
              nix profile add .#default
              if [ "$staged_flake" -eq 1 ]; then
                git restore --staged flake.nix
              fi
            '')
          ];
        };
      });
    };
}
