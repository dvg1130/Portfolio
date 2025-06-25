# uses base64 and XOR encoding with vairable key to

import base64
import argparse
import random
import string
from cryptography.fernet import Fernet

def xor_encode(data, key):
    return ''.join(chr(ord(c) ^ ord(key[i % len(key)])) for i, c in enumerate(data))

def generate_junk_code(lines=5):
    junk = []
    for _ in range(lines):
        var = ''.join(random.choices(string.ascii_letters, k=5))
        val = random.choice(["0", "'junk'", "False", "[]", "{}"])
        fake_if = f"if {var} == {val}: pass"
        junk.append(f"{var} = {val}")
        if random.choice([True, False]):
            junk.append(fake_if)
    return '\n'.join(junk)

def generate_xor_stub(encoded_data, key):
    obf_exec = '+'.join([f"chr({ord(c)})" for c in "exec"])
    decode_fn = """
def _decode(data, key):
    return ''.join([chr(ord(c) ^ ord(key[i % len(key)])) for i, c in enumerate(data)])
"""
    junk_code = generate_junk_code()

    stub = f"""{junk_code}

{decode_fn}
_encoded = {repr(encoded_data)}
_key = {repr(key)}
eval({obf_exec})(_decode(_encoded, _key))
"""
    return stub



def encrypt_aes(data, key=None):
    if key is None:
        key = Fernet.generate_key()
    fernet = Fernet(key)
    encrypted = fernet.encrypt(data.encode())
    return encrypted.decode(), key.decode()

def generate_aes_stub(encrypted_data, key):
    obf_exec = '+'.join([f"chr({ord(c)})" for c in "exec"])
    junk_code = generate_junk_code(lines=6)

    stub = f"""
{junk_code}

from cryptography.fernet import Fernet
import base64

_key = {repr(key)}
_encrypted = {repr(encrypted_data)}
fernet = Fernet(_key)
decrypted = fernet.decrypt(_encrypted.encode()).decode()
eval({obf_exec})(decrypted)
"""
    return stub



def encode_payload(input_file, method, key=None):
    with open(input_file, "r") as f:
        original_code = f.read()

    if method == "base64":
        encoded_code = base64.b64encode(original_code.encode()).decode()
        stub = f"""
import base64
exec(base64.b64decode("{encoded_code}").decode())
"""
    elif method == "xor":
        if not key:
            raise ValueError("XOR encoding requires a key.")
        encoded = xor_encode(original_code, key)
        stub = generate_xor_stub(encoded, key)

    elif method == "aes":
        encrypted, gen_key = encrypt_aes(original_code, key.encode() if key else None)
        stub = generate_aes_stub(encrypted, gen_key)

    else:
        raise ValueError("Unsupported encoding method.")

    return stub

def main():
    parser = argparse.ArgumentParser(description="Wrap a Python script using encoding.")
    parser.add_argument("-i", "--input", required=True, help="Input Python file")
    parser.add_argument("-o", "--output", required=True, help="Output wrapped file")
    parser.add_argument("--method", choices=["base64", "xor", "aes"], default="base64", help="Encoding method")
    parser.add_argument("--key", help="Key for XOR encoding")

    args = parser.parse_args()
    wrapped_code = encode_payload(args.input, args.method, args.key)

    with open(args.output, "w") as f:
        f.write(wrapped_code)

    print(f"[+] Wrapped file saved as {args.output}")

if __name__ == "__main__":
    main()
