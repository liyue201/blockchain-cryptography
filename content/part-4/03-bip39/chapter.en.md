# Mnemonic BIP39: Turning a Private Key into 12 Words

Having you memorize `f917f280……` as a private key isn't realistic. **BIP39 mnemonics** solve that: encode **entropy** into a human-friendly sequence of English words (commonly 12 or 24). It's far friendlier than copying a long hex string — these are those 12 words you're asked to write down when a wallet first launches.

## Core Mechanism

The real BIP39 skeleton is "**entropy + checksum + lookup**":

1. take 128 bits of entropy (12 words) or 256 bits (24 words);
2. append a few **checksum bits** at the end (`checksum length = entropy length / 32`);
3. split the whole string into chunks of 11 bits, each chunk being an index;
4. look up a word in the **standard 2048-word list** by that index.

The checksum is thoughtful: **if you copy a word wrong, restoring it immediately raises the alarm**, so you don't ruin the mnemonic and lose your coins.

> The code on the right replays exactly the same skeleton with a **mini 16-word list** (4 bits per word): generate words, restore the entropy from the words, checksum passes; then **deliberately miscopy one word** and — sure enough — the checksum fails. That's why real BIP39 can "detect writing typos".

## Recovering the Private Key: "Mnemonic → Seed → Private Key"

In the real flow the mnemonic isn't directly the private key; a seed is first derived:

- **seed = PBKDF2(mnemonic, salt = "mnemonic"+passphrase, 2048 iterations)**;
- then BIP32 derives the master private key from the seed (next section).

The optional **passphrase** means "words right but one extra word typed" yields a completely different seed — an extra line of defense once the mnemonic is captured.

> ⚠️ A mnemonic equals the private key. **Whoever gets your mnemonic can move your assets.** Once written down, don't store it in the cloud, don't photograph it, don't screenshot it.

## Why the Checksum Matters

Without a checksum, a missed/miscopied word might not be noticed until the money is long gone; with it, the restoration errors out **at save time** — "fail early". That's why mnemonics, combining "friendliness + self-check", have become the de-facto standard of modern wallets.