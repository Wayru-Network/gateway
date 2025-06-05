# Git Setup Requirements

This document outlines the Git setup requirements for working with the Wayru Network repositories.

## SSH Access to GitHub

To access private repositories, you need to set up SSH authentication with GitHub:

1. Generate a new SSH key pair:
   - Follow the [GitHub guide for generating a new SSH key](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent)
   - Make sure to add the key to your SSH agent (see `Managing SSH Keys` section below)

2. Add your public key to GitHub:
   - Follow the [GitHub guide for adding a new SSH key](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)
   - Ensure you have the necessary permissions in the Wayru Network organization

### Managing SSH Keys

#### Fish Shell
Add the following to your `~/.config/fish/config.fish`:

```fish
# Start ssh-agent if not running
set -q SSH_AUTH_SOCK; or begin
    eval (ssh-agent -c)
    set -Ux SSH_AUTH_SOCK $SSH_AUTH_SOCK
end

# Add key if not already added
ssh-add -l >/dev/null 2>&1; or ssh-add ~/.ssh/github_id_ed25519
```

#### Bash Shell
Add the following to your `~/.bashrc` or `~/.bash_profile`:

```bash
# Start ssh-agent if not running
if [ -z "$SSH_AUTH_SOCK" ]; then
    eval "$(ssh-agent -s)"
    export SSH_AUTH_SOCK
fi

# Add key if not already added
ssh-add -l >/dev/null 2>&1 || ssh-add ~/.ssh/github_id_ed25519
```

#### Zsh Shell
Add the following to your `~/.zshrc`:

```zsh
# Start ssh-agent if not running
if [ -z "$SSH_AUTH_SOCK" ]; then
    eval "$(ssh-agent -s)"
    export SSH_AUTH_SOCK
fi

# Add key if not already added
ssh-add -l >/dev/null 2>&1 || ssh-add ~/.ssh/github_id_ed25519
```

Note: Replace `~/.ssh/github_id_ed25519` with the path to your actual SSH key if different.

## GOPRIVATE Environment Variable

To allow Go to access private repositories, you need to set the `GOPRIVATE` environment variable.

### Configuration by Shell

#### Bash (.bashrc or .bash_profile)
```bash
export GOPRIVATE=github.com/Wayru-Network/*
```

#### Zsh (.zshrc)
```zsh
export GOPRIVATE=github.com/Wayru-Network/*
```

#### Fish (config.fish)
```fish
set -x GOPRIVATE github.com/Wayru-Network/*
```

### Applying the Configuration

After adding the configuration:

1. Save the file
2. Reload your shell configuration:
   - Bash: `source ~/.bashrc` or `source ~/.bash_profile`
   - Zsh: `source ~/.zshrc`
   - Fish: `source ~/.config/fish/config.fish`

### Verification

You can verify the setting is active by running:
```bash
go env GOPRIVATE
```

The output should show `github.com/Wayru-Network/*` 