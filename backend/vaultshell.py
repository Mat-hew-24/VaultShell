import click
from auth import authenticate
from rbac import is_authorized
from session import launch


@click.group()
def cli():
    pass


@cli.command()
@click.argument("container_name")
def enter_container(container_name):
    username = click.prompt("Username")
    password = click.prompt("Password")
    user = authenticate(username, password)
    if not user:
        click.echo("Authentication Failed.")
        return
    if not is_authorized(user, container_name):
        click.echo("Access Denied.")
        return
    launch(user, container_name)


if __name__ == "__main__":
    cli()
