
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.2/instaunit-v1.2-darwin-amd64.tgz"
  sha256 "ca4f6b8678f2b3c82447a6af67975feaef52fb43189530dfa442fea266c21f70"
  version "v1.2"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
