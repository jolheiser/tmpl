# NAME

tmpl - Template automation

# SYNOPSIS

tmpl

```
[--registry|-r]=[value]
[--source|-s]=[value]
```

# DESCRIPTION

Template automation

**Usage**:

```
tmpl [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--registry, -r**="": Registry directory of tmpl (default: ~/.tmpl)

**--source, -s**="": Short-name source to use


# COMMANDS

## download

Download a template

**--branch, -b**="": Branch to clone (default: main)

## env

Show tmpl environment variables

### set, add

Set a tmpl environment variable (stored plaintext)

### unset, delete

Unsets a tmpl environment variable

## init

Initialize a template

## list

List templates in the registry

## remove

Remove a template

## restore

Restore missing templates

## save

Save a local template

## source

Commands for working with sources

### list

List available sources

### add

Add a source

### remove

Remove a source

## test

Test if a directory is a valid template

## update

Update a template

## use

Use a template

**--defaults**: Use template defaults

**--force**: Overwrite existing files
