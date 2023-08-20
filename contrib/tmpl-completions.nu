def _tmpl_env_keys [] {
  ["TMPL_SOURCE", "TMPL_REGISTRY", "TMPL_BRANCH"]
}

def _tmpl_source_list [] {
  ^tmpl source list --json | from json | each { |it| { value: $it.Name, description: $it.URL } }
}

def _tmpl_template_list [] {
  ^tmpl list --json | from json | each { |it| { value: $it.Name, description: (if $it.Path != "" { $it.Path } else { $"($it.Repository)@($it.Branch)" }) } }
}


# Template automation
export extern "tmpl" [
  --registry(-r): string # Registry directory of tmpl (default: ~/.tmpl) [$TMPL_REGISTRY]
  --source(-s): string   # Short-name source to use [$TMPL_SOURCE]
  --help(-h): bool       # Show help
  --version(-v): bool    # Show version
]

# Download a template
export extern "tmpl download" [
  repo_url: string     # Repository URL
  name: string         # Local name for template
  --branch(-b): string # Branch to clone (default: "main") [$TMPL_BRANCH]
  --help(-h): bool     # Show help
]

# Manage tmpl environment variables
export extern "tmpl env" [
  --help(-h): bool # Show help
]

# Set a tmpl environment variable
export extern "tmpl env set" [
  key: string@"_tmpl_env_keys" # Env key
  value: string                # Env value
  --help(-h): bool             # Show help
]

# Unset a tmpl environment variable
export extern "tmpl env unset" [
  key: string@"_tmpl_env_keys" # Env key
]

# Initialize a blank tmpl template
export extern "tmpl init" [
  --help(-h): bool # Show help
]

# List all templates in registry
export extern "tmpl list" [
  --json: bool     # Output in JSON
  --help(-h): bool # Show help
]

# Remove a template
export extern "tmpl remove" [
  name: string     # Name of the template to remove
  --help(-h): bool #Show help
]

# Restore templates present in the registry, but missing archives
export extern "tmpl restore" [
  --help(-h): # Show help
]

# Save a local template
export extern "tmpl save" [
  path: string     # Path to the local template
  name: string     # Name of the template
  --help(-h): bool # Show help
]

# Work with tmpl sources
export extern "tmpl source" [
  --help(-h): # Show help
]

# Add a tmpl source
export extern "tmpl source add" [
  base_url: string # Base URL
  name: string     # Name
  --help(-h): bool # Show help
]

# Remove a tmpl source
export extern "tmpl source remove" [
  name: string@"_tmpl_source_list" # Source to remove
  --help(-h): bool                 # Show help
]

# Test whether a directory is a valid tmpl template
export extern "tmpl test" [
  path?: string    # Path to test (default: ".")
  --help(-h): bool # Show help
]

# Update a template
export extern "tmpl update" [
  name: string@"_tmpl_template_list" # Template to update
  --help(-h):                        # Show help
]

# Use a template
export extern "tmpl use" [
  name: string@"_tmpl_template_list" # The template to execute
  dest?: path                       # Destination for the template (default: ".")
  --defaults: bool                  # Use template defaults
  --force: bool                     # Overwrite existing files
  --help(-h): bool                  # Show help
]
