
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.3.3/instaunit-v1.3.3-darwin-amd64.tgz"
  sha256 "fc5235ebf6351974c9faa81d57016e5a837e72094ceb8f9d30295fdcd6e49edd"
  version "v1.3.3"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
