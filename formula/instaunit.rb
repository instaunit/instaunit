
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.0/instaunit-v1.0-darwin-amd64.tgz"
  sha256 "d055b66fed8a7e5cb286a1a42ef057d4980ee59d2d4e487480eedb1227e5e128"
  
  def install
    system "make", "install", "PREFIX=#{prefix}"
  end
end
