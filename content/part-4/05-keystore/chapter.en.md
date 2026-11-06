# Private Key Custody: KDF, Cold/Warm Storage & Strategy

The private key IS the asset; so "how you store it" decides the money's risk. This section covers two layers: **encrypted storage** (technical) and **cold/warm separation** (strategic).

## Not Just Stored, but Encrypted: KDF

Raw private keys sitting on disk/cloud are a disaster the moment they leak. Professional wallets first encrypt the private key into ciphertext with a **passphrase + KDF** before persisting:

- **salt**: random each time, so even the same passphrase yields different ciphertext;
- **KDF** (e.g. PBKDF2/scrypt/Argon2): runs a **huge number of repeated hashes** on "passphrase + salt" to derive a key, making passphrase guessing extremely slow;
- encrypt the private key with this key using **AES-GCM** — GCM provides **authentication**, so a wrong passphrase is rejected with a clear error, never silently returning a tampered fake private key.

> Run the code on the right: encrypt a wallet file with a strong passphrase → correct passphrase unlocks exactly; **a wrong passphrase is blocked by GCM authentication**. The point of a KDF: passphrases are human-memorized and weak; a slow hash drags each brute-force attempt "down to seconds", turning a million tries into weeks — flattening the gap between "easy to guess" and "hard".

## Cold/Warm Separation Strategy

- **Warm storage**: held by internet-connected devices (phone, exchange) — convenient but a large attack surface;
- **Cold storage**: entirely offline device / paper / hardware wallet — never networked, the safest for private keys;
- in practice: **encrypt the wallet with a strong passphrase**, and let the "private key" exist only on an offline machine, with only ciphertext online. Even if servers are dumped or hard drives stolen, without the passphrase the key can't be unwrapped.

## Custody Checklist (this book's security baseline)

1. Mnemonics/private keys are **never** photographed, screenshotted, cloud-stored, or sent in plaintext;
2. encrypt with a strong passphrase + KDF; don't use a common word (apply last section's "randomness" mindset);
3. large sums go to **cold storage**, small amounts to hot wallets;
4. periodically rotate/re-verify backups, keep multiple copies in different places;
5. remember: **only holding the private key means truly owning the asset** — an exchange can freeze your balance; self-custody is complete.