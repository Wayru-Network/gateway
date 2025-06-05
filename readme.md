# gateway

## Dependencies
- go
- air
- just
- git

## Requirements
- git configured with SSH access to GitHub
- GOPRIVATE env variable

### SSH access to GitHub
- generate public/private key pair and adding to SSH agent: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent
- adding public key to GitHub: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account

### GOPRIVATE env variable
Add the following to your shell configuration file:

For bash (.bashrc or .bash_profile):
```bash
export GOPRIVATE=github.com/Wayru-Network/*
```

For zsh (.zshrc):
```zsh
export GOPRIVATE=github.com/Wayru-Network/*
```

For fish (config.fish):
```fish
set -x GOPRIVATE github.com/Wayru-Network/*
```

After adding the configuration:
1. Save the file
2. Restart your terminal or reload your shell configuration:
   - bash: `source ~/.bashrc` or `source ~/.bash_profile`
   - zsh: `source ~/.zshrc`
   - fish: `source ~/.config/fish/config.fish`

You can verify the setting is active by running:
```bash
go env GOPRIVATE
```


