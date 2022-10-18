package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/** this function is based on the premise that files will be iterated in lexographical order */
func run(root string, maxDepth int, out io.Writer) error {
	prevFile := struct {
		dirString string
		dirEntry fs.DirEntry
	} { "", nil };
	prevFilePrefix := "";
	sep := "-_ ";
	return filepath.WalkDir(root, func(rd string, d fs.DirEntry, err error) error {
		tempd, _ := filepath.Split(rd);
		if err != nil {
			return err;  // general approach here
		}
		if maxDepth > 0 && d.Type().IsDir() && strings.Count(rd, string(os.PathSeparator)) > maxDepth {
			return fs.SkipDir;
		}
		if !d.Type().IsDir() {
			if prevFilePrefix == "" || !strings.HasPrefix(d.Name(), prevFilePrefix) {  // used d.Name() to just get only last file (or folder) name quickly
				if prevFilePrefix != "" {
					fmt.Println("\twill move: ", prevFile.dirString);  // copy the already traversed file (the first of the group) at last
					// fmt.Printf("\tdetails: %s, %s\n", prevFile.dirString, filepath.Join(tempd, prevFilePrefix, prevFile.dirEntry.Name()));
					if err := os.Rename(prevFile.dirString, filepath.Join(tempd, prevFilePrefix, prevFile.dirEntry.Name())); err != nil {
						if !os.IsNotExist(err) {  // do a better implementation for this, basically needed if the directory is having the unique file
							log.Println(err);
							return err;
						}
					};
				}
				if idx := strings.LastIndexAny(d.Name(), sep); idx != -1 {
					prevFilePrefix = d.Name()[:idx];
					fmt.Printf("\t-> %s\n", d.Name());
					prevFile.dirString = rd;
					prevFile.dirEntry = d;
				}
			} else  {
				fmt.Printf("creating directory: %s, %s, %s\n", tempd, prevFilePrefix, filepath.Join(tempd, prevFilePrefix));
				if err := os.MkdirAll(filepath.Join(tempd, prevFilePrefix), 0755); err != nil {
					if !os.IsExist(err) {
						log.Println(err);
						return err;
					}
				}
				fmt.Printf("\t\twill copy: %s\n", d.Name());
				// fmt.Printf("\tdetails: %s, %s\n", rd, filepath.Join(tempd, prevFilePrefix, d.Name()));
				if err := os.Rename(rd, filepath.Join(tempd, prevFilePrefix, d.Name())); err != nil {
					log.Println(err);
					return err;
				}
			}
			// fmt.Printf("%s, %s\n", rd, prevFile);
		} else {
			fmt.Printf("%s\n", rd);
		}
		return nil;
	});
}

// this is just a test function to test the expected behavior before implementation
func filenameSplitTest() {
	sep := "-_ ";
	string := [3]string{"first thing_1.txt", "second-file 2.txt", "third_string-3.txt"};
	for _, str := range string {
		if idx := strings.LastIndexAny(str, sep); idx != -1 {
			fmt.Println(str[:idx]);
		}
	}
}

func main() {
	root := flag.String("root", ".", "specify the root directory for traversal");
	maxDepth := flag.Int("maxdepth", 0, "specify maximum depth to traverse");
	flag.Parse();
	// fmt.Println("root:", *root, string(os.PathSeparator));
	log.SetFlags(log.LstdFlags | log.Lshortfile);
	if err := run(*root, *maxDepth, os.Stdout); err != nil {
		log.Fatal(err);
	}

	// filenameSplitTest();
}
