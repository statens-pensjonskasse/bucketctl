# Settings

## Template

```yaml
<PROJECT_KEY>:
  restrictions:
    <MATCHER_TYPE>:
      <branch>:
        <restriction>:
          exempt-users: [ ]
          exempt-groups: [ ]
  repositories:
    <repo-slug>:
      restrictions:
        <MATCHER_TYPE>:
          <branch>:
            <restriction>:
              exempt-users: [ ]
              exempt-groups: [ ]
```

### Keys

#### `<PROJECT_KEY>`

* Prosjektnøkkel i Bitbucket

#### `<MATCHER_TYPE>`

* `BRANCH` – Navn på branch
* `MODEL_BRANCH` – Type branch (`production`/`development`)

#### `<branch>`

* Branch navn eller type

#### `<restriction>`

* `fast-forward-only` – Hindre omskriving av Git-loggen
* `no-deletes` – Hindre sletting av branch
* `pull-request-only` – Hindre endring av branch utenom PR
* `read-only` – Hindre alle endringer

#### `<repo-slug>`

* Reponavn i Bitbucket

## Branch restrictions

Branch restriction unntak er _permissive_.
Det vil si at dersom skriv til en branch uten PR (`pull-request-only`) er satt på prosjektnivå uten unntak, 
men også er anngitt på repo-nivå med et unntak for en gruppe, så vil dette unntaket gjelde for den gruppa i det gitte 
repoet.

Dersom samme branch dekkes av to forskjellige restrictions gjelder samme regler.
I eksempelet under er `main` er angitt som `development`-branch i `foo`- og `bar`-repoene.
Her vil `ci`-brukeren ha lov til å pushe endringer direkte til `main`-branch i både `foo` og `bar`-repoene uten en PR,
mens andre brukere vil bli stoppet. `ci`-brukeren har også mulighet til å skrive om Git-loggen (e.g. ved force-push)
i `main`-branch i `bar`-repoet, men blir stoppet i `foo`-repoet.

```yaml
PROJ:
  restrictions:
    MODEL_BRANCH:
      development:
        fast-forward-only: {}
        pull-request-only:
          exempt-users:
            - ci
  repositories:
    foo:
      restrictions:
        BRANCH:
          heads/refs/main:
            pull-request-only: {}
    bar:
      restrictions:
        BRANCH:
          heads/refs/main:
            fast-forward-only:
              exempt-users:
                - ci
```

