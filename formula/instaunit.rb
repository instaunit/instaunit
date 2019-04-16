
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/790d65b/instaunit-790d65b-darwin-amd64.tgz"
  sha256 "0d3ae41b447835dee519a70fd5f083e15caef1bce131063c9d8e148f0c30d96e"
  version "790d65b"
  
  def install
    bin.install "bin/instaunit"
  end
end
