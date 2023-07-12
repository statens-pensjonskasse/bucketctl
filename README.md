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
For enkelhestsskyld kan dette lagres i config sammen med base url

```shell
bucketctl config set --token <token> --base-url 
```

## Eksempler

Hent alle prosjekter

```shell
bucketctl project list 
```

Hent alle repositories under `PROJ`-prosjektet

```shell
bucketctl repo list -k PROJ
```

Hent alle tilganger for `PROJ`-prosjektet med alle repos og skriv det til en `.yaml` fil

```shell
bucketctl permission list -k PROJ --include-repos -o yaml > PROJ.yaml
```

Hent abosolutt alle tilganger i Bitbucket og i `.json`-format

```shell
bucketctl permission all -o json --limit 9001
```

Sett tilganger ut fra en fil (`.json` eller `.yaml`)

```shell
bucketctl permission apply -f permissions.yaml --include-repos
```

**NB:** Tilganger til repositories fra fil vil kun bli brukt når `--include-repos` angis.
Repositories som ev. ikke er inkludert i lista vil miste all tilgangsstyring da denne er antatt satt på prosjektnivå.

Hent webhooks for `bucketctl` repoet i `PROJ`-prosjektet

```shell
bucketctl webhook list -k PROJ -r bucketctl
```