# Go Repos Sync

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=nolte_go-repos-sync&metric=bugs)](https://sonarcloud.io/dashboard?id=nolte_go-repos-sync) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=nolte_go-repos-sync&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=nolte_go-repos-sync) [![CI Application Tests](https://github.com/nolte/go-repos-sync/actions/workflows/build-application.yaml/badge.svg)](https://github.com/nolte/go-repos-sync/actions/workflows/build-application.yaml)

---

Tool for keep Local Repository in Sync with public/private Remote Hosts like Github/Gitlab.

## Usage

```sh
  go-repos-sync import bulk
```

### Configuration

For Daily usage is it recommendet to Preconfigure the Import Commands, with Settings like: `Default Checkut Protocol`, `Checkut Base Directory`, and many more.

**Default Config Location:** `~/.repos-sync/config.yaml`  
**Example:** [`./examples/config.yaml`](./examples/config.yaml)

### Bulk Checkout Configuration

For `go-repos-sync import bulk` you will need a Configuration File with the Selected Repository for Sync, from different Remotes.

Examples:  
 - [`./examples/minimal-projects.yaml`](./examples/minimal-projects.yaml)
 - [`./examples/complex-projects.yaml`](./examples/complex-projects.yaml)

The Bulk Configs can be configure at the [Configuration](#configuration), by adding a Path/URL at `settings.bulkElements`, or you use the Commandline Paramaeter `--bulkConfig` for overwrite the [Configuration](#configuration) from file.


## Development

```sh
asdf plugin-add golangci-lint https://github.com/hypnoglow/asdf-golangci-lint.git
asdf plugin-add goreleaser https://github.com/kforsthoevel/asdf-goreleaser.git

asdf install
```


### Build


```sh
goreleaser --snapshot --skip-publish --rm-dist --skip-sign
```



* https://www.source-fellows.de/go-datenbank-orm-association/
