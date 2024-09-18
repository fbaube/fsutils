package fsutils

// https://robthorne-26852.medium.com/a-tale-of-two-file-systems-in-go-b749038c7373

// TODO: Add three slices of strings to specify the filters.

import(
	"io/fs"
	FP "path/filepath"
)

type FilteringFS struct {
    fs fs.FS
}
// And make the wrapper into an fs.FS by implementing its
// interface.
func (wrapper FilteringFS) Open(name string) (fs.File, error) {
    f, err := wrapper.fs.Open(name)
    if err != nil {
        return nil, err
    }
    s, err := f.Stat()
    if err != nil {
        return nil, err
    }
    if s.IsDir() {
        // Have an index file or go home!
        index := FP.Join(name, "index.html")
        if _, err := wrapper.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }
            return nil, err
        }
    }
    return f, nil
}

