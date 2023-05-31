# GoBit
Et CLI-verktøy for BitBucket-APIet skrevet i Go

## Installasjon
```shell
make build && make install
```

## Hjelp
For å komme i gang skriv
```shell 
gobit help
```

For å kunne bruke noen av funksjonene trenger du et access-token,
dette kan lages under [profilen din](https://git.spk.no/account) i Bitbucket.
For enkelhestsskyld kan dette lagres i en config-fil med kommandoen

```shell
gobit config set --token <token>
```

### Eksempler
Hent alle prosjekter

```shell
gobit project list 
```

Hent alle repositories under `INFRA`-prosjektet

```shell
gobit repo list -k INFRA
```

Hent alle tilganger for `INFRA`-prosjektet med alle repos og skriv det til en `.yaml` fil

```shell
gobit permissions list -k --include-repos -o yaml > INFRA.yaml
```

Hent abosolutt alle tilganger i Bitbucket og i `.json`-format

```shell
gobit permissions all -o json --limit 9001
```

Sett tilganger ut fra en fil (`.json` eller `.yaml`)

```shell
gobit permissions apply -f <fil> --include-repos
```

**Obs:** Tilganger til repositories fra fil vil kun bli brukt når `--include-repos` anngis.
Repositories som ev. ikke er inkludert i liste vil miste all tilgangsstyring da denne er antatt satt på prosjektnivå.