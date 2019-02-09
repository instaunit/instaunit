
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.3.1/instaunit-v1.3.1-darwin-amd64.tgz"
  sha256 "6c72c2fa339d5f4e5f96edac9c356a7234c30f5e5978eda253cf493fcfa5c893"
  version "v1.3.1"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
