# Using Deploy-Workflows with Masa-Oracle

This repository, `deploy-workflows`, contains GitHub Actions workflows used by the `masa-oracle` open-source project. To ensure these workflows remain private while `masa-oracle` is public, we use Git submodules. Here's how to set it up:

## Step-by-Step Instructions

### Step 1: Add Deploy-Workflows as a Submodule to Masa-Oracle

1. Clone your `masa-oracle` repository locally (if not already done):
   ```
   git clone git@github.com:masa-finance/masa-oracle.git
   cd masa-oracle
   ```

2. Add `deploy-workflows` as a submodule within the `.github` directory of `masa-oracle`:
   ```
   git submodule add git@github.com:masa-finance/deploy-workflows.git .github/deploy-workflows
   git commit -m "Added deploy-workflows submodule"
   git push
   ```

### Step 2: Configure Submodule to Include Specific Workflow Files

Since Git submodules are designed to include whole repositories, you can't add individual files directly. However, you can control which files are available to the `masa-oracle` repository by managing the contents of the `.github/deploy-workflows` directory.

1. After adding the submodule, navigate to the submodule directory:
   ```
   cd .github/deploy-workflows
   ```

2. Check out a specific branch or tag that contains only the required workflow files. Alternatively, manually delete or adjust files in this directory as needed.

3. Commit these changes in the `masa-oracle` repository:
   ```
   cd ../../
   git add .github/deploy-workflows
   git commit -m "Configured deploy-workflows submodule"
   git push
   ```

### Step 3: Usage

When cloning or forking the `masa-oracle` repository, the `deploy-workflows` submodule will not be included automatically. To fetch the submodule's contents:

```
git submodule init
git submodule update
```

Note: Only users with access to `deploy-workflows` can fetch its contents.

### Conclusion

This setup allows the `masa-oracle` project to be open source while keeping specific GitHub Actions workflows private within `deploy-workflows`.
