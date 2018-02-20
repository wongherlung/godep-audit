## godep-audit

This program makes use of the output of `dep ensure -update -dry-run` to determine which packages are outdated. Based on the `json` whitelist file (format given below), it will output the outdated packages in JUnit XML format.

### Whitelisting format

```json
{
  "whitelisted_packages": [
    {
      "name": "<<package name>>",
      "upstream_commit": "<<upstream commit to ignore>>",
      "reason": "<<reason why not updating>>"
    },
    {
      "name": "<<package name>>",
      "upstream_commit": "<<upstream commit to ignore>>",
      "reason": "<<reason why not updating>>"
    }
  ]
}
```
