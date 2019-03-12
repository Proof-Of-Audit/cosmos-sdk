# sdkch
Simple tool to maintain modular changelogs

# Usage

```
$ sdkch -help
usage: sdkch COMMAND

Maintain unreleased changelog entries in a modular fashion.

Commands:
    init
    add SECTION STANZA            Add an entry file.
                                  Read from stdin until it
                                  encounters EOF.
    generate [VERSION]            Generate a changelog in
                                  Markdown format and print it
                                  to stdout. VERSION defaults
                                  to UNRELEASED.

    Sections             Stanzas
         ---                 ---
    breaking                gaia
    features             gaiacli
improvements            gaiarest
    bugfixes                 sdk
                      tendermint

Flags:
  -d string
    	entry files directory (default "./.pending")
```

## Add a new entry

You can either drop a text file in the appropriate directory or use the `add` command:

```bash
$ sdkch add features gaiacli '#3452 New cool gaiacli command'
```

If no message is provided, a new entry file is opened in an editor is started

## Generate the full changelog

```bash
$ sdkch generate
```

The `-prune` flag would remove the old entry files after the changelog is generated.
