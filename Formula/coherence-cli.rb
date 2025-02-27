class CoherenceCli < Formula
  desc "Oracle Coherence CLI"
  homepage "https://github.com/oracle/coherence-cli"
  version "1.8.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/Oracle-Coherence-CLI-1.8.0-darwin-amd64.pkg"
      sha256 "318526990eac30a49fda8133cb7f08109e4fc035774a2e6ffb64e0e79500748d"
      def install
        system "installer", "-pkg", "Oracle-Coherence-CLI-1.8.0-darwin-amd64.pkg", "-target", "/"
      end
    else
      url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/Oracle-Coherence-CLI-1.8.0-darwin-arm64.pkg"
      sha256 "264e377527f30d0b4b3ec5ab3e644033db81c61682c294812559662303ac8efd"
      def install
        system "installer", "-pkg", "Oracle-Coherence-CLI-1.8.0-darwin-arm64.pkg", "-target", "/"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_32_bit?
        url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/cohctl-1.8.0-linux-386"
        sha256 "5adaea0aa9195942294dd6b73350d766d0be0afc729473f2b183f9472e9db3fd"
      else
        url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/cohctl-1.8.0-linux-amd64"
        sha256 "72f443f9432a932a676408027aa77c2fd401e4702a7f293ad78e0243a3e0bb4e"
      end
    else
      url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/cohctl-1.8.0-linux-arm64"
      sha256 "870a45e9f5679b858450f61add0064dc29ec1aac48ac2c6a9b3bd8d718d1b284"
    end
    def install
      mv Dir["*"].first, "cohctl"
      bin.install "cohctl"
    end
  end

  on_windows do
    if Hardware::CPU.intel?
      url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/cohctl-1.8.0-windows-amd64.exe"
      sha256 "39dc981a94cd938deb3c60dfbd9b797ad316f64df39f1aec6e017553d5a4fac0"
    else
      url "https://github.com/oracle/coherence-cli/releases/download/1.8.0/cohctl-1.8.0-windows-arm.exe"
      sha256 "39dc981a94cd938deb3c60dfbd9b797ad316f64df39f1aec6e017553d5a4fac0"
    end
    def install
      mv Dir["*"].first, "cohctl.exe"
      bin.install "cohctl.exe"
    end
  end

  test do
    system "#{bin}/cohctl", "version"
  end
end