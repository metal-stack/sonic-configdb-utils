# SONiC config_db Utils

Commandline tool for generating a `config_db.json` file for SONiC.

## Usage

Create an `input.yaml` file with config variables.
Then run

```bash
sonic-confidb-utils generate -i input.yaml -o config_db.json -f path/to/platform.json
```

Check the [template-values.yaml](template-values.yaml) file to see examples for all accepted variables.
