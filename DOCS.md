# NAME

tmpl - Template automation

# SYNOPSIS

tmpl

```
[--debug|-d]
[--registry|-r]=[value]
[--source|-s]=[value]
```

**Usage**:

```
tmpl [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--debug, -d**: Debug mode

**--registry, -r**="": Registry directory of tmpl (default: /home/jolheiser/.tmpl)

**--source, -s**="": Short-name source to use


# COMMANDS

## download

Download a template

**--branch, -b**="": Branch to clone (default: main)

## init

Initialize a template

## list

List templates in the registry

## remove

Remove a template

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
