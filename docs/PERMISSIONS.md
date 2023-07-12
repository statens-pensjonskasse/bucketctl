# Permissions

## Template

```yaml
<PROJECT_KEY>:
  default-permission: <DEFAULT_PROJECT_PERMISSION>
  permissions:
    <PROJECT_PERMISSION>:
      groups: [ ]
      users: [ ]
  repositories:
    <repo-slug>:
      permissions:
        <REPO_PERMISSION>:
          groups: [ ]
          users: [ ]
```

### Keys

#### `<PROJECT_KEY>`

* Prosjektnøkkel i Bitbucket

#### `<DEFAULT_PROJECT_PERMISSION>`

Standardtilgang for innloggede brukere

* `PROJECT_WRITE`
* `PROJECT_READ`

#### `<PROJECT_PERMISSION>`

* `PROJECT_ADMIN`
* `PROJECT_WRITE`
* `PROJECT_READ`

#### `<repo-slug>`

* Reponavn i Bitbucket

#### `<REPO_PERMISSION>`

* `REPO_ADMIN`
* `REPO_WRITE`
* `REPO_READ`
