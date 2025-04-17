{
  description = "Language server for NetLinx programming language";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Get version from git or env
        version = self.shortRev or "dev";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "netlinx-language-server";
          inherit version;

          src = ./.;

          vendorHash = null; # Will need to be updated

          ldflags = [
            "-s" "-w"
            "-X main.version=${version}"
            "-X main.commit=${version}"
            "-X main.date=unstable"
          ];

          meta = with pkgs.lib; {
            description = "Language server for NetLinx";
            homepage = "https://github.com/Norgate-AV/netlinx-language-server";
            license = licenses.mit;
            mainProgram = "netlinx-language-server";
          };
        };

        # Development environment
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            golangci-lint
            goreleaser
          ];
        };

        # App definition for nix run
        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };
      });
}
