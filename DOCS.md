# tmpl templates

This documentation aims to cover FAQs and setup.

## Setting up a template

A "valid" tmpl template only requires two things

1. A `template.toml` file in the root directory.
2. A `template` directory that serves as the "root" of the template.

## template.toml

**NOTE:** The template.toml file will be expanded, though not with the full power of the template itself.  
The template.toml file will only expand environment variables with syntax `$USER` or `${USER}`.  
For full documentation on the syntax, see [os.ExpandEnv](https://golang.org/pkg/os/#ExpandEnv).

```toml
# Key-value pairs can be simple
# The user will receive a basic prompt asking them to fill out the variable
project = "my-project"

# Extended properties MUST be added after any simple key-value pairs (due to how TOML works)

# The "key" is enclosed in braces
[author]
# prompt is what will be shown to prompt the user
prompt = "The name of the author of this project"
# help would be extra information (generally seen by giving '?' to a prompt)
help = "Who will be primarily writing this project"
# default is the "value" part of the simple pair. This could be a suggested value
default = "$USER"
```

## template directory

This directory contains any and all files that are part of the template.

Everything in this directory (including paths and file names!) will be executed as a [Go template](https://golang.org/pkg/text/template/).

See the [documentation](https://golang.org/pkg/text/template/) for every available possibility, but some basic examples are...

* A variable defined in template.toml (tmpl allows for keys to be called as a func or variable, whichever you prefer!)
   * `{{project}}` or `{{.project}}`
   * `{{author}}` or `{{.author}}`
* Conditionally including something
   * `{{if eq project ""}} something... {{end}}`

### template helpers

For a full list, see [helper.go](registry/helper.go)

|Helper|Example|Output|
|-----|-----|-----|
|upper|`{{upper project}}`|`MY-PROJECT`|
|lower|`{{lower project}}`|`my-project`|
|title|`{{title project}}`|`My-Project`|
|snake|`{{snake project}}`|`my_project`|
|kebab|`{{kebab project}}`|`my-project`|
|pascal|`{{pascal project}}`|`MyProject`|
|camel|`{{camel project}}`|`myProject`|
|env|`{{env "USER"}}`|The current user|
|sep|`{{sep}}`|Filepath separator for current OS|
|time}|`{{time "01/02/2006"}}`|`11/21/2020` - The time according to the given [format](https://flaviocopes.com/go-date-time-format/)|

## Sources

tmpl was designed to work with any local or git-based template. Unfortunately, in contrast to boilr, this means 
it cannot be used with `user/repo` notation out of the box. 

However, you _can_ set up a source (and subsequent env variable) to make it easier to use your preferred source while
still allowing for others.

### Setting up a source

Let's set up a source for [Gitea](https://gitea.com)

```
tmpl source add https://gitea.com gitea
```

To use it, either pass it in with the `--source` flag

```
tmpl --source gitea download jolheiser/tmpls tmpls
```

Or set it as the env variable `TMPL_SOURCE`

## Using a different branch

By default, tmpl will want to use a branch called `main` in your repository.

If you are using another branch as your default, you can set it as the env variable `TMPL_BRANCH`

Alternatively, you can specify on the command-line with the `--branch` flag of the `download` command

```
tmpl --source gitea download --branch license jolheiser/tmpls license
```
The above command would download the [license](https://gitea.com/jolheiser/tmpls/src/branch/license) template from `jolheiser/tmpls`

## Putting it all together

I realize that many users will be using GitHub, and most will likely still be using the `master` branch.

1. Set up a source for GitHub
   1. `tmpl source add https://github.com github`
   2. Set the env variable `TMPL_SOURCE` to `github`
2. Set the env variable `TMPL_BRANCH` to `master`
3. Happy templating! `tmpl download user/repo repo`