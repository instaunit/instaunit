
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.0.2/instaunit-v1.0.2-darwin-amd64.tgz"
  sha256 "c9ccd87f13d804d7367a6ff3129258152a4db7cff0a2d99da59ea1b2ce6ad4ce"
  version "v1.0.2"
  
  def install
    system "make", "install", "PREFIX=#{prefix}"
  end
end
