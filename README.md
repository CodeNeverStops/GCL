# GCL
Go Count Lines

# Usage
```
➜  gcl git:(master) ./gcl -help
usage: gcl [flags] [dir]
flags:
  -excludedir string
        specify the exclude dirs
  -filetype string
        specify the file type to count line
  -help
        show usage help and quit
  -sortasc
        sort files in ascending order
  -top int
        list top N files
  -version
        show version and quit
```

# Example
```
➜ ./gcl -filetype=".js|.php" -top=10 /your/dir1 /your/dir2
```
