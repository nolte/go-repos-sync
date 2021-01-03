# BoilerCode

```sh

asdf shell python 3.7.1

cookiecutter gh:nolte/cookiecutter-gh-project \
    module_slug="go-repos-sync" \
    topics="management, checkout, gitlab, github" \
    description="Commandline Tool for Sync a set of remote Repos with your local FileSystem." \
    template_issues="y" \
    template_pull_request="y" \
    dependabot_github_actions="y" \
    dependabot_pip="n" \
    dependabot_gitsubmodule="n" \
    dependabot_docker="n" \
    -f --no-input
```
