name: Build and Release pearlS

on:
  push:
    tags:
      - 'v*' # This triggers the action only when you push a tag like "v1.0.0"

jobs:
  build:
    name: Build Binaries
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
            name: pearls-linux-amd64
          - os: windows
            arch: amd64
            name: pearls-windows-amd64.exe
          - os: darwin
            arch: amd64
            name: pearls-mac-amd64
          - os: darwin
            arch: arm64
            name: pearls-mac-arm64

    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21' # Use your Go version here

      - name: Build Binary
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
        run: go build -o ${{ matrix.name }} main.go

      - name: Upload Release Asset
        uses: softprops/action-gh-release@v1
        with:
          files: ${{ matrix.name }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
