package main

import "fmt"
import "os"
import "path/filepath"
import "compress/gzip"
import "io"
import "time"
import "archive/tar"
import "strings"

func main() {
    fmt.Println("Testing archiving function")
    archive("/usr/src/mygoapp/archive", "/usr/src/mygoapp", true)

    time.Sleep(3 * time.Second)

    fmt.Println("Testing compressing function")
    compress("/usr/src/mygoapp/archive.tar", "/usr/src/mygoapp", true)

    time.Sleep(3 * time.Second)

    fmt.Println("Testing decompressing function")
    decompress("/usr/src/mygoapp/archive.tar.gz", "/usr/src/mygoapp", true)

    time.Sleep(3 * time.Second)

    fmt.Println("Testing unarchiving function")
    unarchive("/usr/src/mygoapp/archive.tar", "/usr/src/mygoapp", true)
}

func compress(source, target string, delete bool) (written int64, err error) {
    written = 0

    reader, err := os.Open(source)
    if err != nil { return }
    
    filename := filepath.Base(source)
    target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
    writer, err := os.Create(target)
    if err != nil { return }
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

func decompress(source, target string, delete bool) (written int64, err error) {
    written = 0

    reader, err := os.Open(source)
    if err != nil { return }
    defer reader.Close()

    archive, err := gzip.NewReader(reader)
    if err != nil { return }
    defer archive.Close()

    target = filepath.Join(target, archive.Name)
    writer, err := os.Create(target)
    if err != nil { return }
    defer writer.Close()

    written, err = io.Copy(writer, archive)

    if delete {
        err = os.Remove(source)
    }

    return
}

func archive(source, target string, delete bool) (written int64, err error) {
    written = 0

    filename := filepath.Base(source)
    target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
    tarfile, err := os.Create(target)
    if err != nil { return }
    defer tarfile.Close()

    tarball := tar.NewWriter(tarfile)
    defer tarball.Close()

    info, err := os.Stat(source)
    if err != nil { return }

    var baseDir string
    if info.IsDir() {
        baseDir = filepath.Base(source)
    }

    err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }
        header, err := tar.FileInfoHeader(info, info.Name())
        if err != nil { return err }

        if baseDir != "" {
            header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
        }

        err = tarball.WriteHeader(header)
        if err != nil { return err }

        if info.IsDir() { return nil }

        file, err := os.Open(path)
        if err != nil { return err }
        defer file.Close()
        writtenFile, err := io.Copy(tarball, file)

        written += writtenFile

        // if delete {
        //     fmt.Println("É para excluir")
        //     err = os.Remove(path)
        // }

        return err
    })

    if err != nil { return }

    if delete {
        err = os.RemoveAll(source)
    }

    return
}

func unarchive(source, target string, delete bool) (written int64, err error) {
    written = 0
    reader, err := os.Open(source)
    if err != nil { return }
    defer reader.Close()
    
    tarReader := tar.NewReader(reader)
    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return written, err
        }

        path := filepath.Join(target, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            err = os.MkdirAll(path, info.Mode())
            if err != nil { return written, err }
            continue
        }

        file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil { return written, err }
        defer file.Close()

        writtenFile, err := io.Copy(file, tarReader)
        if err != nil { return written, err }

        written += writtenFile
    }

    if delete {
        err = os.Remove(source)
        if err != nil { return }
    }

    return
}
