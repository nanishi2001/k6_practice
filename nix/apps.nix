{ self, system }:
{
  # アプリケーション定義
  api = {
    type = "app";
    program = "${self.packages.${system}.api}/bin/api";
  };

  default = self.apps.${system}.api;
}
