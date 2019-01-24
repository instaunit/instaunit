
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/1.1/instaunit-1.1-darwin-amd64.tgz"
  sha256 "c7c59ab63089fa07db7c45caab3d3b40a53c3ecb49de9c2613f19e78f5b3c1be"
  version "1.1"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
