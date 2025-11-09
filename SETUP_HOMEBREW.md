# Setting Up Homebrew Tap for HSCTL

This guide will help you set up the initial Homebrew formula in your `obay/homebrew-tap` repository.

## Option 1: Wait for First Release (Recommended)

GoReleaser will automatically generate and commit the formula to your tap repository when you create your first release. This is the easiest approach:

1. **Create your first release**:
   ```bash
   git tag -a v0.1.0 -m "Initial release"
   git push origin v0.1.0
   ```

2. **GoReleaser will automatically**:
   - Generate the formula
   - Commit it to `obay/homebrew-tap` at `Formula/hsctl.rb`
   - Update it with correct SHA256 hashes

3. **Then users can install**:
   ```bash
   brew tap obay/homebrew-tap
   brew install hsctl
   ```

## Option 2: Manual Setup (Before First Release)

If you want to make the formula available before the first release, you can manually create it:

### Step 1: Clone Your Tap Repository

```bash
git clone https://github.com/obay/homebrew-tap.git
cd homebrew-tap
```

### Step 2: Create the Formula Directory

```bash
mkdir -p Formula
```

### Step 3: Create the Formula File

Create `Formula/hsctl.rb` with this template:

```ruby
# typed: false
# frozen_string_literal: true

class Hsctl < Formula
  desc "A CLI tool for managing HubSpot contacts"
  homepage "https://github.com/obay/hsctl"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/obay/hsctl/releases/download/v0.1.0/hsctl_0.1.0_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER"  # Will be updated by GoReleaser

      def install
        bin.install "hsctl"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/obay/hsctl/releases/download/v0.1.0/hsctl_0.1.0_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER"  # Will be updated by GoReleaser

      def install
        bin.install "hsctl"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/obay/hsctl/releases/download/v0.1.0/hsctl_0.1.0_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER"  # Will be updated by GoReleaser

      def install
        bin.install "hsctl"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/obay/hsctl/releases/download/v0.1.0/hsctl_0.1.0_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER"  # Will be updated by GoReleaser

      def install
        bin.install "hsctl"
      end
    end
  end

  test do
    system "#{bin}/hsctl --version"
  end
end
```

### Step 4: Commit and Push

```bash
git add Formula/hsctl.rb
git commit -m "Add hsctl formula"
git push origin main
```

### Step 5: Test Installation

```bash
brew tap obay/homebrew-tap
brew install hsctl
```

**Note:** The formula won't work until you create the first release with the actual binaries. The SHA256 hashes will be placeholders until GoReleaser updates them.

## Option 3: Use GoReleaser Snapshot (For Testing)

You can test the formula generation locally:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Generate a snapshot release (doesn't publish)
goreleaser release --snapshot --skip-publish

# Check the generated formula
cat dist/Formula/hsctl.rb
```

Then manually copy the generated formula to your tap repository.

## After First Release

Once you create your first release, GoReleaser will:
1. Generate the formula with correct SHA256 hashes
2. Automatically commit it to `obay/homebrew-tap/Formula/hsctl.rb`
3. Update it on subsequent releases

## Troubleshooting

### "No available formula with the name 'hsctl'"

This means the formula doesn't exist in your tap yet. Either:
- Wait for the first release (GoReleaser will create it)
- Manually create it using Option 2 above

### Formula exists but installation fails

- Check that the release assets exist at the URLs in the formula
- Verify SHA256 hashes are correct (GoReleaser updates these automatically)
- Ensure the version in the formula matches an actual release

### GoReleaser can't commit to tap

- Verify `GITHUB_TOKEN` has write access to `obay/homebrew-tap`
- Check GitHub Actions logs for permission errors
- Ensure the tap repository exists and is accessible

