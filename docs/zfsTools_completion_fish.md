## zfsTools completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	zfsTools completion fish | source

To load completions for every new session, execute once:

	zfsTools completion fish > ~/.config/fish/completions/zfsTools.fish

You will need to start a new shell for this setup to take effect.


```
zfsTools completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -l, --log-file string   Log file to write to
```

### SEE ALSO

* [zfsTools completion](zfsTools_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 19-Feb-2023