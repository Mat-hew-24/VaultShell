# VaultShell

VaultShell is a small Linux tool that puts **authentication + simple RBAC** in front of **interactive container shells**.

It works by installing a lightweight wrapper for `docker exec` that:

1. Prompts for a VaultShell username + password
2. Verifies access (admin or assigned container)
3. Opens an interactive `bash` session inside the container
4. Logs the session output and (optionally) converts ANSI logs into plain text

## What it does

- **Auth**: Passwords are stored as **bcrypt** hashes.
- **RBAC**:
  - `admin` users can access any container.
  - `user` accounts can only access containers assigned to them.
- **Session logging**:
  - Raw session logs are written to `/var/lib/vaultshell/logs/`.
  - A watcher can convert those logs to text in `/var/lib/vaultshell/logstxt/`.

## How it works (high level)

1. `make install` installs:
   - `/usr/local/bin/vaultshell` (a PyInstaller one-file binary)
   - `/usr/local/bin/docker` (a wrapper that intercepts `docker exec`)
   - `/usr/local/bin/watch.sh` and `/usr/local/bin/logconverter.sh` (log conversion)
   - `/var/lib/vaultshell/` (SQLite DB + logs directories)
2. When you run `docker exec ...`, the wrapper detects `exec` and runs:
   - `vaultshell enter-container <container>`
3. VaultShell checks credentials and authorization, then runs:
   - `/usr/bin/docker exec -it <container> bash`
   - while capturing all output into a timestamped log file.

## Requirements

### Runtime (on the host)

- Linux
- Docker Engine + daemon running
- Docker CLI at **`/usr/bin/docker`** (the wrapper forwards non-`exec` commands there)

### Build / install-time

- `make`
- Python (recommended: 3.10+)
- `pip`

### Optional (log conversion)

These are also listed in `dump.txt`:

- `inotifywait` (from `inotify-tools`)
- `ansi2txt` (a command available on your distro; see install steps below)

## Install (system-wide, recommended)

> Important: the installer puts a wrapper at `/usr/local/bin/docker`.
> On most systems, `/usr/local/bin` comes before `/usr/bin` in `PATH`, so this will affect _all_ users.

### 1) Install system dependencies

Install Docker (if you don’t already have it), plus optional log conversion dependencies:

- Fedora:
  - `sudo dnf install -y inotify-tools`
- Debian/Ubuntu:
  - `sudo apt-get update && sudo apt-get install -y inotify-tools`

For the `ansi2txt` tool, package names differ by distro.
After installing it, these should work:

```bash
command -v inotifywait
command -v ansi2txt
```

### 2) Clone and set up Python deps

```bash
git clone <your-repo-url>
cd VaultShell

python -m venv .venv
source .venv/bin/activate

pip install -r requirements.txt
```

### 3) Build and install

```bash
make install
```

`make install` will:

- Build a one-file binary using PyInstaller
- Copy binaries/scripts into `/usr/local/bin/`
- Create and permission `/var/lib/vaultshell/` like this:
  - `/var/lib/vaultshell/vaultshell.db`
  - `/var/lib/vaultshell/logs/`
  - `/var/lib/vaultshell/logstxt/`

### 4) Verify the wrapper is active

```bash
which docker
which vaultshell
```

Expected:

- `docker` resolves to `/usr/local/bin/docker`
- `vaultshell` resolves to `/usr/local/bin/vaultshell`

### 5) (Recommended) remove direct Docker access for non-admin users

VaultShell is intended to gate interactive container shells.
If a user is in the `docker` group, they can often bypass controls by calling Docker directly.

From `dump.txt`:

```bash
sudo gpasswd -d "$USER" docker
```

Then log out and log back in so group membership updates.

## First-time setup

VaultShell stores its state in an SQLite database at:

- `/var/lib/vaultshell/vaultshell.db`

Tables are created automatically on first run.

### Create users

```bash
vaultshell add-user alice
vaultshell add-user bob
```

### Assign containers to users

Container access is by **container name**:

```bash
vaultshell assign alice my-container
```

If a container isn’t assigned, non-admin users will get **Access Denied**.

### Make a user an admin (manual)

There is no CLI command yet for role management.
To promote a user to `admin`, update `users.role` in `/var/lib/vaultshell/vaultshell.db`.

If you don’t have a `sqlite3` CLI available, you can do it with Python:

```bash
python - <<'PY'
import sqlite3

db = "/var/lib/vaultshell/vaultshell.db"
username = "alice"

conn = sqlite3.connect(db)
conn.execute("UPDATE users SET role='admin' WHERE username=?", (username,))
conn.commit()
conn.close()
print(f"Promoted {username} to admin")
PY
```

## Usage

### Enter a container shell

With the wrapper installed, use `docker exec` as usual:

```bash
docker exec -it my-container bash
```

VaultShell will prompt for credentials and then open a shell (if authorized).

You can also call VaultShell directly:

```bash
vaultshell enter-container my-container
```

### `docker exec` wrapper limitations (important)

The wrapper installed at `/usr/local/bin/docker` is intentionally minimal.

- It only intercepts `docker exec ...` (not `docker container exec`, `docker compose exec`, etc.).
- It always launches `bash` inside the container (it does not forward the command you typed after the container name).
- Avoid `docker exec` flags that take a separate value (for example `-u root` / `--user root`), because the wrapper may mis-detect the container name.

If you need more control, prefer running `vaultshell enter-container <container>` directly.

### Session logs

When a session starts, VaultShell writes a file like:

- `/var/lib/vaultshell/logs/<username>_<container>_YYYYMMDD_HHMMSS.log`

It also records metadata (start/end timestamps + log path) into `session_logs`.

### Plain-text log conversion (watcher)

On session start, VaultShell tries to start `/usr/local/bin/watch.sh`, which runs `logconverter.sh` in the background.
`logconverter.sh`:

- Watches `/var/lib/vaultshell/logs/` using `inotifywait`
- Converts ANSI logs to `.txt` using `ansi2txt`
- Writes results to `/var/lib/vaultshell/logstxt/`

If `inotifywait` or `ansi2txt` is missing, sessions still work, but `.txt` conversion won’t.

## Repo map (what each file is for)

| File                | Purpose                                                                                         |
| ------------------- | ----------------------------------------------------------------------------------------------- |
| `vaultshell.py`     | Click-based CLI entrypoint (`add-user`, `assign`, `enter-container`)                            |
| `auth.py`           | Password hashing (bcrypt) + authentication                                                      |
| `rbac.py`           | Authorization rules (admin vs container owner)                                                  |
| `db.py`             | SQLite schema + CRUD helpers; stores state under `/var/lib/vaultshell/`                         |
| `session.py`        | Launches `docker exec` with a PTY and logs session output                                       |
| `docker-wrapper.sh` | Wrapper installed as `/usr/local/bin/docker` to intercept `docker exec`                         |
| `watch.sh`          | Starts the log conversion watcher in the background                                             |
| `logconverter.sh`   | Watches logs dir and converts ANSI logs to `.txt`                                               |
| `Makefile`          | Builds via PyInstaller and installs binaries + permissions                                      |
| `systool.sh`        | Convenience script: `make build`, `make install`, `make clean`                                  |
| `requirements.txt`  | Python dependencies used for building/testing                                                   |
| `dump.txt`          | Notes about required system packages and group permissions                                      |
| `.gitignore`        | Ignores local artifacts like `.venv/` and `vaultshell.db`                                       |
| `vaultshell.db`     | Repo-local SQLite DB (dev artifact; system installs use `/var/lib/vaultshell/vaultshell.db`)    |
| `watch.log`         | Repo-local watcher output (dev artifact; system installs write `/var/lib/vaultshell/watch.log`) |

## Uninstall

```bash
sudo rm -f /usr/local/bin/vaultshell
sudo rm -f /usr/local/bin/watch.sh /usr/local/bin/logconverter.sh

# removes the docker wrapper (restores system docker resolution to /usr/bin/docker)
sudo rm -f /usr/local/bin/docker

# optional: delete state (users, assignments, logs)
sudo rm -rf /var/lib/vaultshell
```

## Troubleshooting

- **`Authentication Failed.`**
  - Username doesn’t exist or password is wrong.
- **`Access Denied.`**
  - User isn’t `admin` and the container isn’t assigned to them.
- **`Error: Run 'sudo vaultshell' once to initialize.`**
  - VaultShell couldn’t create `/var/lib/vaultshell/`. Run `make install` (recommended), or create the directory with the permissions shown in the `Makefile`.
- **`Container '<name>' not running.` / `not found.`**
  - Start the container, and verify the name matches `docker ps`.
- **After removing yourself from the `docker` group, `docker ps` fails**
  - Expected: the wrapper forwards most Docker commands to `/usr/bin/docker` _without_ elevated group permissions. Use `sudo docker ...` (admin tasks) or keep operators in the `docker` group.
- **Docker CLI not at `/usr/bin/docker`**
  - Update the wrapper in `docker-wrapper.sh` (it hard-codes `/usr/bin/docker`).
- **Shell fails because container has no `bash`**
  - VaultShell currently launches `bash` explicitly. Install `bash` in the image, or change `session.py` to use `sh`.

## Security notes (read this)

- Access to the Docker socket is effectively **root-equivalent** on most systems.
- VaultShell’s default install uses group permissions (`root:docker` + setgid) to allow controlled access.
- Treat `admin` users as highly privileged.
