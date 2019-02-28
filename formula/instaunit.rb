
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.3.1/instaunit-v1.3.1-darwin-amd64.tgz"
  sha256 "777f97dfd29a5d92dcc75ed92235c39532d61542a1bc0d3193d7df140a13a498"
  version "v1.3.1"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
