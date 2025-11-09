# Quick Start: Making HSCTL Available via Homebrew

The Homebrew formula doesn't exist yet because no release has been made. Here's how to fix it:

## Solution: Create Your First Release

GoReleaser will automatically create the Homebrew formula when you make your first release.

### Step 1: Ensure Everything is Committed

```bash
git add .
git commit -m "chore: prepare for initial release"
git push
```

### Step 2: Create and Push a Release Tag

```bash
# Create the tag
git tag -a v0.1.0 -m "Initial release"

# Push the tag (this triggers GitHub Actions)
git push origin v0.1.0
```

### Step 3: Wait for GitHub Actions

GitHub Actions will:
1. Build binaries for all platforms
2. Create a GitHub release
3. **Automatically generate and commit the Homebrew formula to `obay/homebrew-tap`**
4. Generate Scoop manifest for `obay/scoop-bucket`

### Step 4: Verify and Install

After the workflow completes (usually 2-5 minutes):

```bash
# Update your tap
brew tap obay/homebrew-tap

# Now install
brew install hsctl
```

## Alternative: Generate Formula Locally (For Testing)

If you want to test the formula generation before making a release:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Generate a snapshot (doesn't publish)
goreleaser release --snapshot --skip-publish

# Check the generated formula
cat dist/Formula/hsctl.rb
```

Then manually copy `dist/Formula/hsctl.rb` to your `obay/homebrew-tap` repository at `Formula/hsctl.rb`.

## Troubleshooting

### "No available formula with the name 'hsctl'"

This means:
- The formula hasn't been created yet (no release made)
- OR the tap hasn't been updated after a release

**Solution**: Create a release or wait for GoReleaser to finish after a release.

### Formula exists but installation fails

- Check that the release exists: https://github.com/obay/hsctl/releases
- Verify the URLs in the formula point to actual release assets
- Check SHA256 hashes match (GoReleaser updates these automatically)

### GoReleaser can't commit to tap

- Ensure `GITHUB_TOKEN` has write access to `obay/homebrew-tap`
- Check GitHub Actions logs for errors
- Verify the tap repository exists and is accessible

