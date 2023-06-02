# Cross functional calls

## api/checkStatus
```
GET /api/checkStatus
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
## api/delete
```
POST /api/delete

Body:
{
  "files": ["path/to/file1.txt", "path/to/file2.txt", "path/to/file3.txt"]
}
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

## api/install
```
POST /api/install

Body:
{
  "package": ["package1", "package2", "package3"]
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