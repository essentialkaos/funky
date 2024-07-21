<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/funky/ci"><img src="https://kaos.sh/w/funky/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/funky/codeql"><img src="https://kaos.sh/w/funky/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<br/>

`funky` is a simple Yandex.Cloud function for transforming [timer triggers](https://yandex.cloud/en/docs/serverless-containers/concepts/trigger/timer) into HTTP requests.

### Configuration

Entrypoint: `ycfunc.Timer`

#### Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `ASYNC`  | No       | _Send [async requests](https://yandex.cloud/en/docs/functions/concepts/function-invoke-async)_ |

#### Required service roles

- `serverless.containers.invoker`

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
