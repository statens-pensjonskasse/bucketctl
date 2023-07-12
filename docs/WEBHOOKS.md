# Webhooks

## Template

```yaml
<PROJECT_KEY>:
  webhooks:
    - name: <string>
      events: [ ]
      configuration: { }
      url: <string>
      active: <bool>
      scopeType: project
      sslVerificationRequired: <bool>
  repositories:
    <repo-slug>:
      webhooks:
        - name: <string>
          events: [ ]
          configuration: { }
          url: <string>
          active: <bool>
          scopeType: repository
          sslVerificationRequired: <bool>
```

### Keys

#### `<PROJECT_KEY>`

* Prosjektnøkkel i Bitbucket

#### `events`

Noen mulige events er

* pr:merged
* pr:reviewer:updated
* pr:modified
* pr:opened
* repo:refs_changed
* pr:declined
* pr:deleted
* pr:from_ref_updated
