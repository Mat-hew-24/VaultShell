import os, datetime, pty, docker, subprocess, atexit
from db import log_end_session, log_start_session, DATA_DIR

watcher = None


def start_watcher():
    global watcher
    if watcher is None or watcher.poll() is not None:
        watcher = subprocess.Popen(["bash", "/usr/local/bin/watch.sh"])
        atexit.register(lambda: watcher.terminate())


def launch(user, container_name):
    client = docker.from_env()
    try:
        container = client.containers.get(container_name)
        if container.status != "running":
            print(f"Container '{container_name}' not running.")
            return
    except docker.errors.NotFound:
        print(f"Container '{container_name}' not found.")
        return
    start_watcher()
    logs_dir = DATA_DIR / "logs"
    logs_dir.mkdir(exist_ok=True)
    log_path = (
        logs_dir
        / f"{user['username']}_{container_name}_{datetime.datetime.now():%Y%m%d_%H%M%S}.log"
    )
    session_id = log_start_session(user["username"], container_name, str(log_path))
    with open(log_path, "wb") as f:
        pty.spawn(
            ["/usr/bin/docker", "exec", "-it", container_name, "bash"],
            master_read=lambda fd: (lambda d: (f.write(d), d)[1])(os.read(fd, 1024)),
        )
    log_end_session(session_id)
