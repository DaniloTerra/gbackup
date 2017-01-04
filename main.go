package main

import "fmt"
import "os"
import "path/filepath"
import "compress/gzip"
import "io"
import "time"

func main() {
    fmt.Println("Testing compressing function")
    compress("/usr/src/mygoapp/toCompress.json", "/usr/src/mygoapp", true)

    time.Sleep(10 * time.Second)

    fmt.Println("Testing decompressing function")
    decompress("/usr/src/mygoapp/toCompress.json.gz", "/usr/src/mygoapp", true)
}

// source = absolute filepath
// target = absolute path. The filename is going to be the same from source
// delete = the source file should be deleted?
func compress(source, target string, delete bool) (written int64, err error) {
    written = 0

    reader, err := os.Open(source)
    if err != nil {
        return
    }
    
    filename := filepath.Base(source)
    target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
    writer, err := os.Create(target)
    if err != nil {
        return
    }
    defer writer.Close()

    archiver := gzip.NewWriter(writer)
    archiver.Name = filename
    defer archiver.Close()

    _, err = io.Copy(archiver, reader)

    if delete == true {
        err = os.Remove(source)
    }

    return
}

// source = absolute filepath
// target = absolute path. The filename is going to be the same from source
// delete = the source file should be deleted?
func decompress(source, target string, delete bool) (written int64, err error) {
    written = 0

    reader, err := os.Open(source)
    if err != nil {
        return
    }
    defer reader.Close()

    archive, err := gzip.NewReader(reader)
    if err != nil {
        return
    }
    defer archive.Close()

    target = filepath.Join(target, archive.Name)
    writer, err := os.Create(target)
    if err != nil {
        return
    }
    defer writer.Close()

    written, err = io.Copy(writer, archive)

    if delete {
        err = os.Remove(source)
    }

    return
}

// func archive() (written int64, err error) {}

// func unarchive() (written int64, err error) {}
