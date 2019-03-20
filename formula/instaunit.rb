
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.3.2/instaunit-v1.3.2-darwin-amd64.tgz"
  sha256 "6e473b73164b979e63441fbc6c7f3a7718eab35ef2c583e8f16951ee16e3dc23"
  version "v1.3.2"
  
  def install
    system "install", "-d", "#{bin}"
    system "install", "-m", "0755", "bin/instaunit", "#{bin}/instaunit"
  end
end
