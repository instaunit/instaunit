
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/3d44460/instaunit-3d44460-darwin-amd64.tgz"
  sha256 "054e6d4dd4763b96ea7ead4cf9b845afd49577da83669fc218bb6f9154284ec0"
  version "3d44460"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
