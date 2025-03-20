# SONiC config_db Utils

Commandline tool for generating a `config_db.json` file for SONiC.

## Usage

Create a `sonic-config.yaml` file with config variables.
Then run

```bash
sonic-confidb-utils generate
```

Check the [template-values.yaml](template-values.yaml) file to see examples for all accepted variables.

## Limitations

### LLDP

When configuring LLDP `hello_time`, a simple `config load` may not be enough.
Depending on what SONiC version you are running, the lldp container might need to be restarted for the change to become effective, either with `systemctl restart lldp` or with `config reload`.
