# gimme-cli

Gimme CDN CLI Tool

## Description

The CLI has been designed to publish modules to a CDN instance from the CI tools.

## Configuration

To use the CLI, you must first configure these variables.

| Variable    | Description                                                                                                      |
|-------------|------------------------------------------------------------------------------------------------------------------|
| GIMME_CDN_URL   | CDN instance URL                                                                                                 |
| GIMME_TOKEN | A valid CDN access token ([see create access token doc](https://github.com/gimme-cdn/gimme#create-access-token)) |

## Usage

You can push a module to your CDN instance with the following command :
```shell
./gimme-cli --name=<module-name> --version=<module-version> --file=<file-path>
```

## CI Integration
