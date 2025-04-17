#!/usr/bin/env sh

# Get hash for vendored dependencies
VENDOR_HASH=$(go mod vendor && nix hash path --type sha256 --sri vendor)

# Update flake.nix with the correct version and hash
sed -i "s/version = \".*\";/version = \"$VERSION\";/" flake.nix
sed -i "s/vendorHash = \".*\";/vendorHash = \"$VENDOR_HASH\";/" flake.nix
