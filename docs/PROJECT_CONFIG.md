## ProjectConfig

```yaml
kind: ProjectConfig
apiVersion: bucketctl.spk.no/v1alpha1
metadata:
  name: <name>
spec:
  projectKey: <PROJECT_KEY>
  public: <bool>
  defaultPermission: <PROJECT_WRITE|PROJECT_READ|REPO_CREATE>
  permissions: <[]Permission>
  branchRestrictions: <[]BranchRestriction>
  webhooks: <[]Webhook>
  repositories: <[]Repository>
```

## Permission

```yaml
name: <ProjectPermission|RepoPermission>
groups: <[]string>
users: <[]string>
```

### ProjectPermission

* `PROJECT_ADMIN`
* `PROJECT_WRITE`
* `PROJECT_READ`

### RepoPermission

* `REPO_ADMIN`
* `REPO_WRITE`
* `REPO_READ`

## BranchRestriction

```yaml
type: <MatcherType>
matchers:
  - matching: <string>
    restrictions:
      - type: <Restriction>
        exempt-groups: <[]string>
        exempt-users: <[]string>
```

Branch restriction exemptions are permissive.
This means that if a branch is write-protected (`pull-request-only`) without exemptions on project level,
but with an exemption on repository level, then that exemption holds true.
If a branch is covered by two rules and an exemption is given in one of them,
that exemption also holds true.

### MatcherType

* `BRANCH` – Branch name
* `MODEL_BRANCH` – Type of branch (`production`/`development`)
* `MODEL_CATEGORY` – Branch category (`FEATURE`/`BUGFIX`/`HOTFIX`/`RELEASE`)
* `PATTERN` – Pattern,
  see [Bitbucket docs](https://confluence.atlassian.com/bitbucketserver/branch-permission-patterns-776639814.html) for
  details

### Restriction

* `fast-forward-only` – Prevent rewriting Git history
* `no-deletes` – Prevent branch deletion
* `pull-request-only` – Prevent branch changes, except through PRs
* `read-only` – Prevent all changes

## Webhook

```yaml
name: <string>
events: <[]WebhookEvent>
configuration: { }
url: <string>
active: <bool>
scopeType: <project|repository>
sslVerificationRequired: <bool>
```

### WebhookEvent

* pr:merged
* pr:reviewer:updated
* pr:modified
* pr:opened
* repo:refs_changed
* pr:declined
* pr:deleted
* pr:from_ref_updated

## Repository

```yaml
name: <string>
permissions: <[]Permission>
branchRestrictions: <[]BranchRestriction>
webhooks: <[]Webhook>
```
