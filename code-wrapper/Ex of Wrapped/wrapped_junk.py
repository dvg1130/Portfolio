ElYGM = {}
if ElYGM == {}: pass
lDzbD = False
VXvfD = 'junk'
qengj = []
ZXjim = 'junk'


def _decode(data, key):
    return ''.join([chr(ord(c) ^ ord(key[i % len(key)])) for i, c in enumerate(data)])

_encoded = '\x03\x17\n\x1c\x11\\\x13f[\x1a\x16C\x1b\x16TP\x12C\x12\x1c\x0f\x1d\x04\x10\x11AP\x01\x0c\x13\x06KV\x188'
_key = 'secret123'
eval(chr(101)+chr(120)+chr(101)+chr(99))(_decode(_encoded, _key))
