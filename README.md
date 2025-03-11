# SONiC config_db Utils

Commandline tool for generating a `config_db.json` file for SONiC.

## Usage

Create a `sonic-config.yaml` file with config variables.
Then run

```bash
sonic-confidb-utils generate
```

Check the [template-values.yaml](template-values.yaml) file to see examples for all accepted variables.

## Supported version and limitations

`sonic-configdb-utils` was only tested with our own [build](https://github.com/metal-stack/sonic-build) of SONiC.

### Limitations

When configuring LLDP `hello_time`, a simple `config load` is not enough.
The lldp container needs to be restarted for the change to become effective, either with `systemctl restart lldp` or with `config reload`.
