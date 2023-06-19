package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// The type of element in the array in the "builds.json" file.
type Build struct {
	Name             string            `json:"name"`
	Category         string            `json:"category"`
	Collectibles     []Collectible     `json:"collectibles"`
	BannedCharacters []BannedCharacter `json:"bannedCharacters"`
}

type Collectible struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BannedCharacter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

var (
	builds = make([]Build, 0)
)

func loadAllBuilds() {
	libPath := getLibraryPath()
	jsonFilePath := path.Join(libPath, "builds.json")
	var jsonFile []byte
	if v, err := ioutil.ReadFile(jsonFilePath); err != nil {
		log.Fatal("Failed to open \""+jsonFilePath+"\":", err)
	} else {
		jsonFile = v
	}

	if err := json.Unmarshal(jsonFile, &builds); err != nil {
		log.Fatal("Failed to unmarshal the builds:", err)
	}
}

func getLibraryPath() string {
	libPath := path.Join(projectPath, "lib", "node_modules", "isaac-racing-common", "src")
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		log.Fatal("The library path at \"" + libPath + "\" does not exist. Did you forget to run \"npm install\" in the \"lib\" subdirectory?")
	} else if err != nil {
		log.Fatal("Failed to check if the \""+libPath+"\" file exists:", err)
	}

	return libPath
}

func getBuildObjectFromBuildName(name string) Build {
	for _, build := range builds {
		if build.Name == name {
			return build
		}
	}

	log.Fatal("Failed to find a build matching the build name of: " + name)
	return builds[0]
}
