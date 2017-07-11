
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.0.1/instaunit-v1.0.1-darwin-amd64.tgz"
  sha256 "6ab14f02e5a19663894451d88c2ba69c55623755d1937aae0f003dc8b1763e35"
  version "v1.0.1"
  
  def install
    system "make", "install", "PREFIX=#{prefix}"
  end
end
