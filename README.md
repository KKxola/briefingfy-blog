# Briefingfy Blog

Pequeno blog server-side rendered em Go. Os posts ficam em Markdown e sao carregados do diretorio `content/posts`.

## Como rodar

```sh
go run ./cmd/web
```

O servidor usa a porta `8080` por padrao. Para alterar:

```sh
PORT=3000 go run ./cmd/web
```

No PowerShell:

```powershell
$env:PORT = "3000"; go run ./cmd/web
```

## Como testar

```sh
go test ./...
```

## Como criar um post

Crie um arquivo `.md` em `content/posts` com front matter simples:

```markdown
---
title: Titulo do post
slug: titulo-do-post
date: 2026-05-21
summary: Resumo curto do post.
---
# Titulo do post

Conteudo em Markdown.
```

Campos obrigatorios:

- `title`
- `slug`
- `date`, no formato `YYYY-MM-DD`
- `summary`

## Escopo inicial

Esta versao usa apenas `net/http`, `html/template`, `testing`, `httptest` e `github.com/yuin/goldmark`. Nao ha banco de dados, painel admin, autenticacao ou framework web externo.

## Deploy gratuito no Render

O projeto pode ser publicado como um Web Service gratuito no Render, sem Docker.

1. Suba este repositorio para o GitHub.
2. Crie uma conta em https://render.com.
3. No dashboard do Render, escolha **New > Web Service**.
4. Conecte o repositorio do GitHub.
5. Configure o servico com:
   - Build Command: `go build -tags netgo -ldflags '-s -w' -o app ./cmd/web`
   - Start Command: `./app`
6. Mantenha a porta configurada pela variavel de ambiente `PORT`. O app ja le `PORT` e usa `8080` como fallback local.
7. Aguarde o deploy terminar e acesse a URL `*.onrender.com` gerada pelo Render.

Limitacoes do plano gratuito:

- O servico pode dormir depois de um periodo sem acessos.
- O primeiro acesso depois do servico dormir pode demorar.
- Alteracoes feitas no filesystem do servidor nao sao persistentes.
- Posts devem ser editados nos arquivos Markdown em `content/posts` e publicados por commit no repositorio.
