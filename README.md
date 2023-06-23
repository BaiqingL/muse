<p align="center">
  <img src="resources/gradient.svg" alt="Logo" width="64">
</p>

# Muse

[Demo Video](https://www.youtube.com/watch?v=vA2KrzqN6Hw)

Muse is a Chrome extension that automatically generates Typescript front-end code given a prompt & UI framework. It can also automatically install any dependencies if needed. From there, users can select specific elements & prompt additional iterations via additional prompts.

## API Spec

### api/coldStart
```
POST /api/coldStart

Body:
{
  "framework": "...",
  "useCase": "...",
  "apiKey": "...",
}

```

### api/getFile
```
GET /api/getFile

Query:
{
  "filename": "...",
}

Response:
{
  "exist": true/false,
  "content": "...",
}

```

### api/writeFile
```
POST /api/writeFile

Body:
{
  "filename": "...",
  "content": "...",
}

Response:
{
  "success": true,
}

```

### api/iterate
```
POST /api/iterate

Body:
{
  "html": "...",
  "prompt": "...",
}
```

### api/export
```
GET /api/export

Response:

Zip file data
```
