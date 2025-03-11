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

Open source SONiC does not support configuring LLDP.
To enable configuring the LLDP `hello_time` a patch was added to our build.
However, for the LLDP config to become effective the `lldp` container needs to be restarted, e.g. with `systemctl restart lldp` or `config reload -y`.
