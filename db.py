import sqlite3
from pathlib import Path

DB_PATH = Path.home() / "vaultshell.db"


def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn


def init_db():
    conn = get_db()
    cursor = conn.cursor()
    cursor.executescript(
        """
            CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            role TEXT NOT NULL DEFAULT 'user'
        );

        CREATE TABLE IF NOT EXISTS container_assignments (
            container_name TEXT NOT NULL,
            username TEXT NOT NULL,
            FOREIGN KEY (username) REFERENCES users(username)
        );

        CREATE TABLE IF NOT EXISTS session_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            container_name TEXT NOT NULL,
            started_at TEXT,
            ended_at TEXT,
            log_file TEXT
                );
                """
    )
    conn.commit()
    conn.close()


def get_user(username: str):
    conn = get_db()
    user = conn.execute(
        "SELECT * from users WHERE username = ?", (username,)
    ).fetchone()
    conn.close()
    return user


def create_user(username: str, password_hash: str, role: str = "user"):
    conn = get_db()
    conn.execute(
        "INSERT INTO users (username,password_hash,role) VALUES (?,?,?)",
        (username, password_hash, role),
    )
    conn.commit()
    conn.close()


def get_container_owner(container_name: str):
    conn = get_db()
    row = conn.execute(
        "SELECT username FROM container_assignments where container_name = ?",
        (container_name,),
    ).fetchone()
    conn.close()
    return row["username"] if row else None


def assign_container(username: str, container_name: str):
    conn = get_db()
    conn.execute(
        "INSERT INTO container_assignments (container_name,username) VALUES (?,?)",
        (container_name, username),
    )
    conn.commit()
    conn.close()


def log_start_session(username: str, container_name: str, log_file: str) -> int:
    conn = get_db()
    cursor = conn.execute(
        "INSERT INTO session_logs (username, container_name, started_at, log_file) VALUES (?,?,datetime('now'),?)",
        (username, container_name, log_file),
    )
    session_id = cursor.lastrowid
    conn.commit()
    conn.close()
    return session_id


def log_end_session(session_id: int):
    conn = get_db()
    conn.execute(
        "UPDATE session_logs SET ended_at = datetime('now') WHERE id = ?",
        (session_id,),
    )
    conn.commit()
    conn.close()
