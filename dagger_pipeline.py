import os
import anyio
import dagger


async def test(client, src):
    """Run pytest."""
    python = (
        client.container()
        .from_("python:3.9-slim")
        .with_workdir("/app")
        .with_mounted_directory("/app", src)
        .with_file(
            "/app/requirements.txt",
            src.file("python-app/requirements.txt")
        )
    )

    dependencies = python.with_exec([
        "pip", "install", "-r", "requirements.txt"
    ])
    test_result = dependencies.with_exec(
        ["pytest", "test-app.py"]
    )
    return test_result


async def lint(client, src):
    """Run flake8 linting."""
    python = (
        client.container()
        .from_("python:3.9-slim")
        .with_workdir("/app")
        .with_mounted_directory("/app", src)
        .with_file(
            "/app/requirements.txt",
            src.file("python-app/requirements.txt")
        )
    )

    dependencies = python.with_exec(["pip", "install", "flake8"])
    lint_result = dependencies.with_exec(["flake8", "app.py"])
    return lint_result


async def docker_build(client, src, tag, username):
    """Build Docker image."""
    dockerfile_context = src.directory("python-app")
    print(f"Building Docker image with tag: {tag}")

    try:
        container = (
            client.container()
            .from_("python:3.9-slim")
            .with_mounted_directory("/src", dockerfile_context)
            .with_workdir("/src")
        )

        # Build the image
        await container.with_exec([
            "docker",
            "build",
            "-t",
            f"{username}/docker-dagger-poc:{tag}",  # Push without docker.io prefix
            "."
        ])
        print(f"Successfully built image with tag {tag}")
    except Exception as e:
        print(f"Failed to build Docker image: {str(e)}")
        raise


async def docker_push(client, username, token, tag):
    """Push Docker image to DockerHub."""
    repository = f"{username}/docker-dagger-poc"
    image_tag = f"{repository}:{tag}"
    print(f"Pushing image {image_tag} to DockerHub")
    
    try:
        container = client.container().from_("python:3.9-slim")
        
        # Login to DockerHub with credentials from GitHub secrets
        await container.with_exec([
            "docker",
            "login",
            "-u", username,
            "-p", token,
        ])
        print(f"Successfully logged into DockerHub as {username}")

        # Push the image to your specific repository
        await container.with_exec([
            "docker",
            "push",
            image_tag  # Push using the repository name directly
        ])
        print(f"Successfully pushed image {image_tag}")
        return True
    except Exception as e:
        print(f"Failed to push image to {repository}: {str(e)}")
        raise


async def update_deployment(client, src, tag, username):
    """Update deployment.yml with new image tag."""
    deployment_path = "infrastructure/deployment.yml"

    # Run yq directly on the host
    await client.with_exec([
        "yq",
        "-i",
        f'.spec.template.spec.containers[0].image = "{username}/docker-dagger-poc:{tag}"',
        deployment_path
    ])


async def commit_and_push(client, src, updated_file, tag):
    """Commit and push the updated deployment file."""
    new_src = (
        client.container()
        .from_("alpine")
        .with_mounted_directory("/src", src)
        .with_workdir("/src")
        .with_file("infrastructure/deployment.yml", updated_file)
    )

    git = (
        client.container()
        .from_("alpine/git")
        .with_mounted_directory("/src", new_src.directory("."))
        .with_workdir("/src")
        .with_exec(["git", "config", "--local", "user.email", "srivatsa.rv@one2n.in"])
        .with_exec(["git", "config", "--local", "user.name", "srivatsa.rv"])
        .with_exec(["git", "add", "infrastructure/deployment.yml"])
        .with_exec([
            "git",
            "commit",
            "-m",
            f"Update deployment.yml with image tag {tag}"
        ])
        .with_exec(["git", "push", "origin", "main"])
    )
    return git


async def main():
    """Main execution function."""
    username = os.environ.get("DOCKERHUB_USERNAME")
    token = os.environ.get("DOCKERHUB_TOKEN")

    if not username or not token:
        raise ValueError("DOCKERHUB_USERNAME and DOCKERHUB_TOKEN must be set")

    with open("python-app/VERSION", "r") as f:
        tag = f.read().strip()

    print(f"Starting pipeline with tag: {tag}")

    async with dagger.Connection() as client:
        src = client.host().directory(".", exclude=[".git"])
        src_with_git = client.host().directory(".", include=[".git"])

        if os.environ.get("CI"):
            print("Running in CI environment...")
            
            print("Running tests...")
            await test(client, src)

            print("Running linting...")
            await lint(client, src)

            # Separate build step
            print(f"Building Docker image with tag: {tag}...")
            await docker_build(client, src, tag, username)

            # Separate push step
            print("Pushing to DockerHub...")
            push_success = await docker_push(client, username, token, tag)

            if push_success:
                print("Updating deployment...")
                updated_file = await update_deployment(client, src, tag, username)

                print("Committing and pushing changes...")
                await commit_and_push(client, src_with_git, updated_file, tag)
                
                print("CI pipeline completed successfully")
            else:
                raise Exception("Failed to push to DockerHub, deployment.yml not updated")


if __name__ == "__main__":
    anyio.run(main)
