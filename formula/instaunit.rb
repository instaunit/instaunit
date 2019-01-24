
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.0.3/instaunit-v1.0.3-darwin-amd64.tgz"
  sha256 "ad42ac112cd5f45f5f266d557b4aaab54f7688a9d818bf99299446c8d86e02ce"
  version "v1.0.3"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
