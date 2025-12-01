# KCST

Temporary file hosting.

## Retention Policy

```
min_age  = 1 hour
max_age  = 28 days
max_size = 100 MiB

retention = min_age + (max_age - min_age) * (1 - sqrt(size/max_size))

   days
     28 |.
        | ..
        |   ...
        |      ....
        |          .....
        |               ......
        |                     .......
        |                            ........
      1 |                                    ................
        +-------------------------------------------------->
        0                    50                          100
                                                         MiB
```

Smaller files are retained longer. A 100 MiB file lives ~1 hour, while tiny files can stay up to 28 days.

## Uploading Files

Send a `POST` request with `multipart/form-data` containing a `file` field.

| Field  | Description                  |
|--------|------------------------------|
| `file` | The file to upload (max 100 MiB) |

## cURL Examples

```bash
# Upload a file
curl -F 'file=@yourfile.png' https://example.com

# Upload from stdin
echo "hello world" | curl -F 'file=@-;filename=hello.txt' https://example.com

# Upload with a custom filename
curl -F 'file=@localfile.bin;filename=custom.bin' https://example.com
```

## Example TTLs

| File Size | Retention |
|-----------|-----------|
| 100 MiB   | ~1 hour   |
| 50 MiB    | ~9 days   |
| 25 MiB    | ~14 days  |
| 10 MiB    | ~19 days  |
| 1 MiB     | ~25 days  |
| <1 KiB    | ~28 days  |

## Running

```bash
go run cmd/kcst/main.go
```

The server starts on `:8080` by default. Files are stored in `./uploads/` and metadata in `./data/kcst.db`.
