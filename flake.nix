{
  description = "tmpl";

  inputs.nixpkgs.url = "github:nixos/nixpkgs";

  outputs = {
    self,
    nixpkgs,
  }: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
    tmpl = pkgs.buildGoModule {
      pname = "tmpl";
      version = "0.4.0";

      src = ./.;

      vendorHash = "sha256-QNwzHC4fHLAhshOplKmMjRYa9sHNjBLdfBgANbs/iKk=";

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
