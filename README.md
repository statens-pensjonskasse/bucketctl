# bucketctl

Et CLI-verktøy for BitBucket-APIet skrevet i Go

## Installasjon

```shell
make install
```

## Hjelp

For å komme i gang skriv

```shell 
bucketctl help
```

For å kunne bruke noen av funksjonene trenger du et access-token,
dette kan lages under [profilen din](https://git.spk.no/account) i Bitbucket.
For enkelhestsskyld kan dette lagres i en config-fil med kommandoen

```shell
bucketctl config set --token <token>
```

### Eksempler

Hent alle prosjekter

```shell
bucketctl project list 
```

Hent alle repositories under `INFRA`-prosjektet

```shell
bucketctl repo list -k INFRA
```

Hent alle tilganger for `INFRA`-prosjektet med alle repos og skriv det til en `.yaml` fil

```shell
bucketctl permission list -k --include-repos -o yaml > INFRA.yaml
```

Hent abosolutt alle tilganger i Bitbucket og i `.json`-format

```shell
bucketctl permission all -o json --limit 9001
```

Sett tilganger ut fra en fil (`.json` eller `.yaml`)

```shell
bucketctl permission apply -f permissions.yaml --include-repos
```

**Obs:** Tilganger til repositories fra fil vil kun bli brukt når `--include-repos` anngis.
Repositories som ev. ikke er inkludert i liste vil miste all tilgangsstyring da denne er antatt satt på prosjektnivå.

Hent webhooks for jenkins-pipeline-library repoet i INFRA-prosjektet

```shell
bucketctl webhook list -k INFRA -r jenkins-pipeline-library
```