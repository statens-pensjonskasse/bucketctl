## ProjectConfig

```yaml
kind: ProjectConfig
apiVersion: bucketctl.spk.no/v1alpha1
metadata:
  name: <name>
spec:
  projectKey: <PROJECT_KEY>
  access:
    public: <bool>
    defaultPermission: <PROJECT_WRITE|PROJECT_READ>
    permissions: <[]Permission>
  branchRestrictions: <[]BranchRestriction>
  webhooks: <[]Webhook>
  repositories: <[]Repository>
```

## Repository

```yaml
name: <string>
permissions: <[]Permission>
branchRestrictions: <[]BranchRestriction>
webhooks: <[]Webhook>
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

Branch restriction unntak er _permissive_.
Det vil si at dersom skriv til en branch uten PR (`pull-request-only`) er satt på prosjektnivå uten unntak,
men også er anngitt på repo-nivå med et unntak for en gruppe, så vil dette unntaket gjelde for den gruppa i det gitte
repoet.

Dersom samme branch dekkes av to forskjellige restrictions gjelder samme regler.
I eksempelet under er `main` er angitt som `development`-branch i `foo`- og `bar`-repoene.
Her vil `ci`-brukeren ha lov til å pushe endringer direkte til `main`-branch i både `foo` og `bar`-repoene uten en PR,
mens andre brukere vil bli stoppet. `ci`-brukeren har også mulighet til å skrive om Git-loggen (e.g. ved force-push)
i `main`-branch i `bar`-repoet, men blir stoppet i `foo`-repoet.

### MatcherType

* `BRANCH` – Navn på branch
* `MODEL_BRANCH` – Type branch (`production`/`development`)
* `MODEL_CATEGORY` – Branch-kategori (`FEATURE`/`BUGFIX`/`HOTFIX`/`RELEASE`)
* `PATTERN` – Match på pattern (
  se [Bitbucket docs](https://confluence.atlassian.com/bitbucketserver/branch-permission-patterns-776639814.html))

### Restriction

* `fast-forward-only` – Hindre omskriving av Git-loggen
* `no-deletes` – Hindre sletting av branch
* `pull-request-only` – Hindre endring av branch utenom PR
* `read-only` – Hindre alle endringer

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
