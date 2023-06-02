# Cross functional calls

## api/delete
```
GET /api/delete?file=path/to/file.txt
```
## api/patch
```
POST /api/patch

Body:
{
  "patches": [
    {
      "file": "path/to/file.txt",
      "diff": "..."
    },
    {
      "file": "path/to/another-file.txt",
      "diff": "..."
    }
  ]
}
```

## api/create
```
POST /api/create

Body:
{
  "file": "path/to/new-file.txt",
  "contents": "..."
}
```

## api/install
```
POST /api/install

Body:
{
  "package": "package-name"
}
```

## api/operation
```
POST /api/operation

Body:
{
  "action": "start"
}
```