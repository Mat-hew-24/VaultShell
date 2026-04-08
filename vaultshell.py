import click
from auth import authenticate, hashed_pass
from rbac import is_authorized
from session import launch
from db import init_db, create_user, assign_container


@click.group()
def cli():
    pass


@cli.command()
@click.argument("container_name")
def enter_container(container_name):
    username = click.prompt("Username")
    password = click.prompt("Password", hide_input=True)
    user = authenticate(username, password)
    if not user:
        click.echo("Authentication Failed.")
        return
    if not is_authorized(user, container_name):
        click.echo("Access Denied.")
        return
    launch(user, container_name)


@cli.command()
@click.argument("username")
def add_user(username):
    password = click.prompt("Password", hide_input=True, confirmation_prompt=True)

    try:
        create_user(username, hashed_pass(password))
        click.echo(f"User '{username}' created.")
    except Exception as e:
        click.echo(f"Error: {e}")


@cli.command()
@click.argument("username")
@click.argument("container_name")
def assign(username, container_name):
    try:
        assign_container(username, container_name)
        click.echo(f"Success: Container '{container_name}' assigned to '{username}'.")
    except Exception as e:
        click.echo(f"Error: {e}")


if __name__ == "__main__":
    init_db()
    cli()
