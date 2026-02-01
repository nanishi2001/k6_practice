{ pkgs, self, system }:
{
  # API サーバーパッケージ
  api = pkgs.buildGoModule {
    pname = "k6-practice-api";
    version = "1.0.0";
    src = ../api;
    vendorHash = "sha256-deZ8L5aju1JraGTnjIW3vR1zm5Jc15F6D+goi1pVLpU=";

    meta = with pkgs.lib; {
      description = "k6 practice API server";
      license = licenses.mit;
    };
  };

  default = self.packages.${system}.api;
}
