# bucketctl

Et CLI-verktøy for BitBucket-APIet skrevet i Go

Config blir lagret under `$HOME/.config/bucketctl/config.yaml`

## Installasjon

```shell
make install
```

## Hjelp

`bucketctl` har en innebygd hjelpekommando

```shell 
bucketctl help
```

For å kunne bruke noen av funksjonene trenger du et access-token,
dette kan lages under profilen din i Bitbucket.
For enkelhestsskyld kan dette lagres i config sammen med base url og Git url

```shell
bucketctl config set --token <token> --base-url <url> --git-url <ssh>
```

Det er også mulig å lage forskjellige kontekster som bygger oppå basisconfig.
En ny kontekst kan f.eks. lages med

```shell
bucketctl config context create --context infra --key INFRA --include-repos
```

for å heller bruke kunne konteksten `-x infra` enn å skrive `--project INFRA --include-repos` hver gang.

## Eksempler

### Basic

Hent alle prosjekter

```shell
bucketctl get projects
```

Hent alle repositories under `PROJ`-prosjektet

```shell
bucketctl get repos -p PROJ
```

For å klone alle repositories i et prosjekt

```shell
bucketctl git clone -p PROJ
```

### Git

brukes.
Alle repositories blir da klonet inn under `PROJ`-mappa hvis ingen andre mapper er gitt.

For å synkronisere hovedbranchen i alle repositories i et prosjekt som ligger i mappa `prosjekt` kjøres kommandoen

```shell
bucketctl git clone -p PROJ --update
```

### Prosjektkonfigurasjon

`bucketctl` kan brukes til å endre prosjektkonfigurasjon

Sjekk gjeldende konfigurasjon for et prosjekt med

```shell
bucketctl get project-config -p PROJ
```

Etter å ha endret prosjektkonfigurasjonen kan forskjell fra gjeldende konfigurasjon sjekkes ved

```shell
bucketctl apply -f <PROJ>.yaml --dry-run
```

Se [dokumentasjon](./docs/PROJECT_CONFIG.md) for flere detaljer.