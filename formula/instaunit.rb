
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.3/instaunit-v1.3-darwin-amd64.tgz"
  sha256 "caddd8d5826d24b57e3313a9b38cab22b86865bba8d2aa41c9ef522b0adaaff3"
  version "v1.3"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
