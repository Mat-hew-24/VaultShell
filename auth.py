import bcrypt
from db import get_user


def hashed_pass(password: str) -> str:
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()


def authenticate(username: str, password: str):
    user = get_user(username)
    if not user:
        return None
    if bcrypt.checkpw(password.encode(), user["password_hash"].encode()):
        return user
    return None
