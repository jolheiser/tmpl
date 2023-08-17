# tmpl templates

This documentation aims to cover FAQs and setup.

## Setting up a template

A "valid" tmpl template only requires two things

1. A `tmpl.yaml` file in the root directory.
2. A `template` directory that serves as the "root" of the template.

## tmpl.yaml

**NOTE:** The tmpl.yaml file will be expanded, though not with the full power of the template itself.  
The tmpl.yaml file will only expand environment variables with syntax `$USER` or `${USER}`.  
For full documentation on the syntax, see [os.ExpandEnv](https://golang.org/pkg/os/#ExpandEnv).

When using the `--defaults` flag, no prompts will be shown and only default values will be used.  
As another alternative, any environment variable that matches a key will bypass the prompt.  
For example, `author` would have the corresponding environment variable `TMPL_VAR_AUTHOR`.

```yaml
# tmpl.yaml
# Write any template args here to prompt the user for, giving any defaults/options as applicable

prompts:
  - id: project                           # The unique ID for the prompt
    label: Project Name                   # The prompt message/label
    help: The name to use in the project  # Optional help message for the prompt
    default: tmpl                         # Prompt default
```

## template directory

This directory contains any and all files that are part of the template.

Everything in this directory (including paths and file names!) will be executed as a [Go template](https://golang.org/pkg/text/template/).

See the [documentation](https://golang.org/pkg/text/template/) for every available possibility, but some basic examples are...

* An id defined in `tmpl.yaml` (tmpl allows for keys to be called as a func or variable, whichever you prefer!)
   * `{{project}}` or `{{.project}}`
   * `{{author}}` or `{{.author}}`
* Conditionally including something
   * `{{if eq project ""}} something... {{end}}`

### template helpers

For a full list, see [helper.go](registry/helper.go)

| Helper      | Example                            | Output                                                                                                |
|-------------|------------------------------------|-------------------------------------------------------------------------------------------------------|
| upper       | `{{upper project}}`                | `MY-PROJECT`                                                                                          |
| lower       | `{{lower project}}`                | `my-project`                                                                                          |
| title       | `{{title project}}`                | `My-Project`                                                                                          |
| snake       | `{{snake project}}`                | `my_project`                                                                                          |
| kebab       | `{{kebab project}}`                | `my-project`                                                                                          |
| pascal      | `{{pascal project}}`               | `MyProject`                                                                                           |
| camel       | `{{camel project}}`                | `myProject`                                                                                           |
| env         | `{{env "USER"}}`                   | The current user                                                                                      |
| sep         | `{{sep}}`                          | Filepath separator for current OS                                                                     |
| time        | `{{time "01/02/2006"}}`            | `11/21/2020` - The time according to the given [format](https://flaviocopes.com/go-date-time-format/) |
| trim_prefix | `{{trim_prefix "foobar" "foo"}}`   | `bar`                                                                                                 |
| trim_suffix | `{{trim_suffix "foobar" "bar"}}`   | `foo`                                                                                                 |
| replace     | `{{replace "foobar" "bar" "baz"}}` | `foobaz`                                                                                              |

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
The above command would download the [license](https://git.jojodev.com/jolheiser/tmpls/src/branch/license) template from `jolheiser/tmpls`

## Putting it all together

I realize that many users will be using GitHub, and most will likely still be using the `master` branch.

1. Set up a source for GitHub
   1. `tmpl source add https://github.com github`
   2. Set the env variable `TMPL_SOURCE` to `github`
2. Set the env variable `TMPL_BRANCH` to `master`
3. Happy templating! `tmpl download user/repo repo`

## Backup and Restore

1. The simplest solution is to make a copy of your `registry.yaml` (default: `~/.tmpl/registry.yaml`).
   * Once in the new location, you will need to use `tmpl restore`.
   
2. Alternatively, you can copy/paste the entire registry (default: `~/.tmpl`) and skip the restore step.

## `.tmplkeep`

Perhaps you are familiar with `.gitkeep` and its unofficial status in git. Git does not like empty directories, so usually
a `.gitkeep` (or just `.keep`) file is added to retain the directory while keeping it effectively empty.

tmpl instead uses `.tmplkeep` files for this purpose. The difference is, tmpl will **not** create the `.tmplkeep` file
when the template is executed. This allows you to set up directory structures (for staging, examples, etc.) that
will *actually* be empty after execution.