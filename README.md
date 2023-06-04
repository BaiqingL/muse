# Cross functional calls

## api/coldStart
```
POST /api/coldStart

Body:
{
  "framework": "...",
  "useCase": "...",
  "apiKey": "...",
}

```

## api/getFile
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

## api/writeFile
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

## api/iterate
```
POST /api/iterate

Body:
{
  "html": "...",
  "prompt": "...",
}
```

## api/export
```
GET /api/export

Response:

Zip file data
```