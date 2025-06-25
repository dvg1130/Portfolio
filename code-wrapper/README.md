# Code Wrapper & Obfuscator Tool (XOR / AES / Base64)

This project demonstrates techniques used to wrap and obfuscate Python code. The tool simulates behaviors commonly found in malware and red team tools, while also providing defenders a chance to analyze how such obfuscation affects visibility and static detection.

## Overview

This wrapper script reads a Python file and outputs an obfuscated version using one of several encoding methods. It can optionally insert junk code to add noise and interfere with static analysis.

## Features

- Supports three encoding techniques:
  - base64 encoding
  - XOR encoding with custom or random key
  - AES encryption using Python’s cryptography (Fernet) module
- Obfuscates critical functions such as `exec()` using character-wise `chr()` calls
- Adds junk code to waste time and complicate reverse engineering

## Red team

Using Base64 with XOR provides agile, low-footprint obfuscation ideal for evading static detection in early-stage payloads. AES encryption, while heavier, is valuable for safeguarding second-stage logic during post-exploitation, especially in realistic APT-style simulations.


## Blue team

AES encryption can harden defensive logic by concealing response triggers or detection algorithms, making reverse engineering significantly more difficult—especially when keys are derived dynamically and decrypted only in-memory during runtime.


## Usage

Base64
 - python wrap.py -i payload.py -o wrapped.py
 - python3 wrap.py -i payload.py -o wrapped.py

XOR
 - python wrap.py -i payload.py -o wrapped_xor.py --method xor --key secret123
 - python3 wrap.py -i payload.py -o wrapped_xor.py --method xor --key secret123


AES
 - Random key
   - python wrap.py -i payload.py -o wrapped_aes.py --method aes
   - python3 wrap.py -i payload.py -o wrapped_aes.py --method aes
 - Custom Key
   - python wrap.py -i payload.py -o wrapped_aes.py --method aes --key "kKpdiLG-7p4DfrMyBxA6GIDuBTVfxmptYQ9MyVKcyMA="
   - python3 wrap.py -i payload.py -o wrapped_aes.py --method aes --key "kKpdiLG-7p4DfrMyBxA6GIDuBTVfxmptYQ9MyVKcyMA="



