# GCL
Go Count Lines

# Usage
```
➜ ./gcl -help
usage: gcl [flags] [dir]
flags:
  -filetype string
        Specify the file type to count line
  -help
        show usage help and quit
  -top int
        list top N files
  -version
        show version and quit
```

# Example
```
➜ ./gcl -filetype=".js|.php" -top=10 /your/dir1 /your/dir2
```

# TODO List
- [ ] Add exclude dirs flag