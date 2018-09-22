// Simple file server for recursive image crawling

package viewer

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type FileData []byte
type DirEntrySet mapset.Set
type DirContents []fuse.DirEntry

type ImageFs struct {
	pathfs.FileSystem
	Root    string
	BaseUrl string

	// mapping from full path of a file to its fuse.Attr, including dir
	Attrs map[string]fuse.Attr

	// mapping from full path of a file to its data, excluding dir
	Contents map[string]FileData

	// mapping from full path of a directory to all file entries under it
	Entries map[string]DirEntrySet

	// mapping from dir name to real url
	Urls map[string]string
}

func ToDirEntries(set DirEntrySet) []fuse.DirEntry {
	it := set.Iterator()
	result := make([]fuse.DirEntry, 0)
	for elem := range it.C {
		result = append(result, elem.(fuse.DirEntry))
	}
	return result
}

func (fs *ImageFs) fullpath(src string, base string) string {
	if base == "" && src == fs.Root {
		return "/"
	} else {
		return filepath.Join(base, src)
	}
}

// getData accesses to given url and returns images data and all hrefs
func (fs *ImageFs) getData(link string, base string) (DirContents, error) {
	crawlData, err := Crawl(link)
	if err != nil {
		return nil, err
	}
	fixBase := base
	if base == "" {
		fixBase = "/"
	}
	fs.Entries[fixBase] = mapset.NewSet()
	for _, data := range crawlData {
		fullpath := fs.fullpath(data.Name, base)
		if data.Type == Image {
			fs.Attrs[fullpath] = fuse.Attr{
				Mode:  fuse.S_IFREG | 0644,
				Size:  uint64(len(data.Data)),
				Atime: uint64(time.Now().Unix()),
				Mtime: uint64(time.Now().Unix()),
				Ctime: uint64(time.Now().Unix()),
			}
			fs.Contents[fullpath] = data.Data
			fs.Entries[fixBase].Add(
				fuse.DirEntry{Name: data.Name, Mode: fuse.S_IFREG})
		} else if data.Type == Href {
			fs.Attrs[fullpath] = fuse.Attr{
				Mode:  fuse.S_IFDIR | 0755,
				Atime: uint64(time.Now().Unix()),
				Mtime: uint64(time.Now().Unix()),
				Ctime: uint64(time.Now().Unix()),
			}
			fs.Entries[fixBase].Add(
				fuse.DirEntry{Name: data.Name, Mode: fuse.S_IFDIR})
			fs.Urls[data.Name] = data.Url
		}
	}
	return ToDirEntries(fs.Entries[fixBase]), nil
}

func (fs *ImageFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr name: %s", name)
	if name == "" {
		name = "/"
	}
	if attr, ok := fs.Attrs[name]; ok {
		return &attr, fuse.OK
	} else if name == "/" {
		attr := fuse.Attr{
			Mode:  fuse.S_IFDIR | 0755,
			Ctime: uint64(time.Now().Unix()),
		}
		fs.Attrs[name] = attr
		return &attr, fuse.OK
	} else {
		return nil, fuse.ENOENT
	}
}

func (fs *ImageFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir name: %s", name)
	fixName := name
	if name == "" {
		fixName = "/"
	}
	if entry, ok := fs.Entries[fixName]; ok {
		return ToDirEntries(entry), fuse.OK
	} else {
		var link string
		if name == "" {
			link = fs.BaseUrl
		} else {
			fields := strings.Split(name, string(os.PathSeparator))
			dirname := fields[len(fields)-1]
			var tok bool
			if link, tok = fs.Urls[dirname]; !tok {
				// should not here
				log.Printf("url not found for %s", dirname)
				return nil, fuse.ENOENT
			}
		}
		entries, err := fs.getData(link, name)
		if err != nil {
			log.Printf("get data from src with error: %s", err)
		} else {
			return entries, fuse.OK
		}
	}
	return nil, fuse.ENOENT
}

func (fs *ImageFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("Open name: %s", name)
	if data, ok := fs.Contents[name]; ok {
		return nodefs.NewDataFile(data), fuse.OK
	} else {
		return nil, fuse.ENOENT
	}
}

func Serve(root string, baseUrl string) {
	fs := ImageFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Root:       root,
		BaseUrl:    baseUrl,
		Attrs:      make(map[string]fuse.Attr),
		Contents:   make(map[string]FileData),
		Entries:    make(map[string]DirEntrySet),
		Urls:       make(map[string]string),
	}
	nfs := pathfs.NewPathNodeFs(&fs, nil)
	server, _, err := nodefs.MountRoot(root, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
