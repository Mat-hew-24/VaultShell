from db import get_container_owner


def is_authorized(user: dict, container_name: str) -> bool:
    if user["role"] == "admin":
        return True
    owner = get_container_owner(container_name)
    return owner == user["username"]
