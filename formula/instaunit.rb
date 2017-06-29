
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/untagged-6b913b5df19f6f7bfa8f/instaunit-v1.0-darwin-amd64.tgz"
  sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
  
  def install
    system "make", "install", "PREFIX=#{prefix}"
  end
end
