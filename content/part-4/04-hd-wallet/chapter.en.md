# HD Wallet BIP32/BIP44: A Tree of Keys

If every address had to store its own private key separately, a wallet would become a sack of keys and backups a chore. **HD (Hierarchical Deterministic) wallets** solve it: **deterministically derive a whole tree of keys from a single master private key** — as many addresses as you want, as far down the tree as you need.

## Derivation Formula (BIP32)

From a parent private key `d` and parent chain code `c`, derive a child private key:

- append `data = 0x00 ‖ ser256(d) ‖ ser32(index)` (`0x00` requests a private key);
- compute `(IL, IR) = HMAC-SHA512(key=c, data)` (deterministic, so repeated derivations give the same result);
- child private key = `(IL + d) mod n`, child chain code = `IR`.

As long as the root private key is fixed, every `index` stably yields a deterministic child key — that's what "deterministic" means.

> Run the code on the right to derive layer by layer along **m/44'/0'/0'/0/0** from the root key; each layer yields a unique deterministic child key, assembling a key tree.

## Paths and BIP44

BIP44 assigns each layer a **meaning** so wallets can recognize one another:

```
m / 44' / coin' / account' / change / address_index
```

- `44'`: fixed "coin layer" marker;
- `coin'`: distinguishes BTC, ETH, etc.;
- `account'`: separates accounts (privacy: each account kept apart);
- `change`: 0 = receiving address, 1 = change address;
- `address_index`: which address.

The ones with `'` are **hardened derivations** (use the parent private key); normal derivation only needs the parent public key — so HD wallets can also expose public keys externally as "watch-only wallets" that can't move coins.

## Why It Changes the Experience

- **Back up once**: just save the master seed (those 12 mnemonics) and the whole tree restores — no matter how many addresses, one backup suffices;
- **Privacy**: each transaction can use a fresh address; coins don't all sit on one key;
- **Ecosystem unification**: any BIP44-compliant wallet importing the same seed sees the same addresses and balances.

> ⚠️ Note determinism is "same seed, same result" — so the master seed (mnemonics) IS the root of the wallet; lose it and you lose the whole tree.