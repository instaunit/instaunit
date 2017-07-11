
class Instaunit < Formula
  homepage "https://github.com/instaunit/instaunit"
  url "https://github.com/instaunit/instaunit/releases/download/v1.0.1/instaunit-v1.0.1-darwin-amd64.tgz"
  sha256 "c13ae4690b8e4f8d25ad200b64eb6c7c21b6e03d2e73177119406d7cc7810e89"
  
  def install
    system "make", "install", "PREFIX=#{prefix}"
  end
end
