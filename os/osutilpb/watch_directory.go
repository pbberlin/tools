package osutilpb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	// _ "gopkg.in/fsnotify.v1"
)

var WorkingDir = "c:\\temp\\"
var FilePathSeparator = "\\"
var serverPort = 1313
var liveReload = false

var SourceFs afero.Fs = new(afero.OsFs)
var DestinationFS afero.Fs = new(afero.OsFs)
var OsFs afero.Fs = new(afero.OsFs)

func init() {
	// WatchDir()
}

func getDirList() []string {
	return []string{"ch1", "ch2"}

}

func WatchDir() {

	// Watch runs its own server as part of the routine
	watched := getDirList()
	workingDir := AbsPathify(WorkingDir)
	for i, dir := range watched {
		watched[i], _ = GetRelativePath(dir, workingDir)
	}
	unique := strings.Join(RemoveSubpaths(watched), ",")

	fmt.Printf("Watching for changes in %s/{%s}\n", workingDir, unique)
	err := NewWatcher(serverPort)
	if err != nil {
		fmt.Println(err)
	}

}

func AbsPathify(inPath string) string {
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	// todo consider move workingDir to argument list
	return filepath.Clean(filepath.Join(WorkingDir, inPath))
}

func MakeStaticPathRelative(inPath string) (string, error) {
	staticDir := GetStaticDirPath()
	themeStaticDir := GetStaticDirPath()

	return MakePathRelative(inPath, staticDir, themeStaticDir)
}

func MakePathRelative(inPath string, possibleDirectories ...string) (string, error) {

	for _, currentPath := range possibleDirectories {
		if strings.HasPrefix(inPath, currentPath) {
			return strings.TrimPrefix(inPath, currentPath), nil
		}
	}
	return inPath, errors.New("Can't extract relative path, unknown prefix")
}

func GetStaticDirPath() string {
	return AbsPathify("static")
}

// GetRelativePath returns the relative path of a given path.
func GetRelativePath(path, base string) (final string, err error) {
	if filepath.IsAbs(path) && base == "" {
		return "", errors.New("source: missing base directory")
	}
	name := filepath.Clean(path)
	base = filepath.Clean(base)

	name, err = filepath.Rel(base, name)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(filepath.FromSlash(path), FilePathSeparator) && !strings.HasSuffix(name, FilePathSeparator) {
		name += FilePathSeparator
	}
	return name, nil
}

// RemoveSubpaths takes a list of paths and removes everything that
// contains another path in the list as a prefix. Ignores any empty
// strings. Used mostly for logging.
//
// e.g. ["hello/world", "hello", "foo/bar", ""] -> ["hello", "foo/bar"]
func RemoveSubpaths(paths []string) []string {
	a := make([]string, 0)
	for _, cur := range paths {
		// ignore trivial case
		if cur == "" {
			continue
		}

		isDupe := false
		for i, old := range a {
			if strings.HasPrefix(cur, old) {
				isDupe = true
				break
			} else if strings.HasPrefix(old, cur) {
				a[i] = cur
				isDupe = true
				break
			}
		}

		if !isDupe {
			a = append(a, cur)
		}
	}

	return a
}

// NewWatcher creates a new watcher to watch filesystem events.
func NewWatcher(port int) error {

	if runtime.GOOS == "darwin" {
		// tweakLimit()
	}

	watcher, err := NewWatcher2(1 * time.Second)
	var wg sync.WaitGroup

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer watcher.Close()

	wg.Add(1)

	for _, d := range getDirList() {
		if d != "" {
			_ = watcher.Add(d)
		}
	}

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				log.Println("File System Event:", evs)

				staticChanged := false
				dynamicChanged := false
				staticFilesChanged := make(map[string]bool)

				for _, ev := range evs {
					ext := filepath.Ext(ev.Name)
					istemp := strings.HasSuffix(ext, "~") || (ext == ".swp") || (ext == ".swx") || (ext == ".tmp") || strings.HasPrefix(ext, ".goutputstream")
					if istemp {
						continue
					}
					// renames are always followed with Create/Modify
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						continue
					}

					isstatic := true
					if isstatic {
						if staticPath, err := MakeStaticPathRelative(ev.Name); err == nil {
							staticFilesChanged[staticPath] = true
						}
					}

					// add new directory to watch list
					if s, err := os.Stat(ev.Name); err == nil && s.Mode().IsDir() {
						if ev.Op&fsnotify.Create == fsnotify.Create {
							watcher.Add(ev.Name)
						}
					}
				}

				if staticChanged {
					log.Printf("Static file changed, syncing\n\n")
					copyStatic()

					// livereload.ForceRefresh()

				}

				if dynamicChanged {
					fmt.Print("\nChange detected, rebuilding site\n")
				}
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	if port > 0 {
		if liveReload {
		}

		// go serve(port)
	}

	wg.Wait()
	return nil
}

// ==================0

type Batcher struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}

	Events chan []fsnotify.Event // Events are returned on this channel
}

func NewWatcher2(interval time.Duration) (*Batcher, error) {
	watcher, err := fsnotify.NewWatcher()

	batcher := &Batcher{}
	batcher.Watcher = watcher
	batcher.interval = interval
	batcher.done = make(chan struct{}, 1)
	batcher.Events = make(chan []fsnotify.Event, 1)

	if err == nil {
		go batcher.run()
	}

	return batcher, err
}

func (b *Batcher) run() {
	tick := time.Tick(b.interval)
	evs := make([]fsnotify.Event, 0)
OuterLoop:
	for {
		select {
		case ev := <-b.Watcher.Events:
			evs = append(evs, ev)
		case <-tick:
			if len(evs) == 0 {
				continue
			}
			b.Events <- evs
			evs = make([]fsnotify.Event, 0)
		case <-b.done:
			break OuterLoop
		}
	}
	close(b.done)
}

func (b *Batcher) Close() {
	b.done <- struct{}{}
	b.Watcher.Close()
}

func copyStatic() error {
	staticDir := AbsPathify("static") + "/"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Println("Unable to find Static Directory:", staticDir)
		return nil
	}

	publishDir := AbsPathify("public") + "/"

	syncer := fsync.NewSyncer()
	syncer.NoTimes = true
	syncer.SrcFs = SourceFs
	syncer.DestFs = DestinationFS

	// Copy Static to Destination
	log.Println("syncing from", staticDir, "to", publishDir)
	return syncer.Sync(publishDir, staticDir)
}
