# Security Policy

## Supported versions

Security fixes are applied to the latest `main` branch.

## Reporting a vulnerability

Please open a private GitHub security advisory, or contact the maintainer via https://github.com/Fracizz.

Do **not** open a public issue for credential leaks or remote code execution.

## Threat model

### Assets

- SSH passwords stored in `~/.sshfrac/servers.json`
- Private key paths referenced by `key_file` (keys themselves live on disk)
- Ability to run remote commands as configured users

### Trust boundaries

| Attacker | What they can do |
|----------|------------------|
| Network eavesdropper | Cannot read SSH session contents (SSH/SFTP). Host key verification (default) resists naive MITM; `--insecure` disables this. |
| Same OS user, no master password | Can decrypt **enc:v1** inventory (machine-derived key). |
| Same OS user, with master password | Needs `SSHFRAC_MASTER_PASSWORD` / `--master-password` (and matching `--bind-machine` if used) to decrypt **enc:v2**. |
| Other OS users | Should not read `servers.json` (mode 0600) or your private keys if permissions are correct. |
| Malware as your user | Can read env, keys, and drive the CLI — treat the workstation as trusted. |

### Encryption schemes

| Prefix | KDF / key | When used |
|--------|-----------|-----------|
| `enc:v1:` | SHA-256(machine material) → AES-256-GCM | Default when no master password is set (legacy / convenience). |
| `enc:v2:` | Argon2id(master password [+ optional machine bind], random salt) → AES-256-GCM | When `--master-password` or `SSHFRAC_MASTER_PASSWORD` is set. |

**Recommendations**

1. Prefer **enc:v2** with a strong master password on shared or multi-user machines.
2. Optionally set `--bind-machine` / `SSHFRAC_BIND_MACHINE=1` so ciphertext will not decrypt on another host even with the same master password.
3. Agents should pass the master password via env, not commit it.
4. Do not use `--insecure` outside trusted labs.
5. Prefer SSH public keys (`key_file`) over passwords when possible.

OS keychain integration is not built-in; wrap sshfrac with a script that loads the secret from your keychain into `SSHFRAC_MASTER_PASSWORD`.
