{
  description = "tmpl";

  inputs.nixpkgs.url = "github:nixos/nixpkgs";

  outputs = {
    self,
    nixpkgs,
  }: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
    tmpl = pkgs.buildGoModule rec {
      pname = "tmpl";
      version = "0.4.0";

      src = ./.;

      vendorHash = "sha256-QNwzHC4fHLAhshOplKmMjRYa9sHNjBLdfBgANbs/iKk=";

      ldflags = ["-s" "-w" "-X=go.jolheiser.com/tmpl/cmd.Version=${version}"];

      postInstall = ''
        mkdir -p $out/share
        cp -vr ./contrib/tmpl-completions.nu $out/share/tmpl-completions.nu
      '';

      meta = with pkgs.lib; {
        description = "";
        homepage = "https://git.jojodev.com/jolheiser/tmpl";
        license = licenses.mit;
        maintainers = with maintainers; [jolheiser];
        mainProgram = "tmpl";
      };
    };
  in {
    packages.x86_64-linux.default = tmpl;
  };
}
