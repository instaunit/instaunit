
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.2.3/instaunit-v1.2.3-darwin-amd64.tgz"
  sha256 "10f92c8e12e3fa7381d6349b72620053f8bfe848746fa997d43ef1a45c0431e3"
  version "v1.2.3"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
