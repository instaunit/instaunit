# Maintainer: Brian W. Wolter <brianwolter@gmail.com>
pkgname=instaunit
pkgver=v1.9.1
pkgrel=1
pkgdesc="Instaunit tests your Web APIs"
arch=('x86_64')
url="https://github.com/instaunit/instaunit"
license=('BSD')
groups=()
depends=('go: Go')
makedepends=()
optdepends=()
provides=()
conflicts=()
replaces=()
backup=()
options=()
install=
changelog=
source=($pkgname-$pkgver.tar.gz)
noextract=()
md5sums=() #autofill using updpkgsums

build() {
  cd "$pkgname-$pkgver"
  make
}

package() {
  cd "$pkgname-$pkgver"
  make install DESTDIR="$pkgdir/"
}