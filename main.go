package main

import (
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

func localFeed(path string) (*gofeed.Feed, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return gofeed.NewParser().Parse(f)
}

func remoteFeed(url string) (*gofeed.Feed, error) {
	return gofeed.NewParser().ParseURL(url)
}

// ComparesFeeds compares GUIDs of local and remote feeds. It uses publish times to identify items.
func CompareFeeds(localPath, url string) error {
	local, err := localFeed(localPath)
	if err != nil {
		return err
	}

	remote, err := remoteFeed(url)
	if err != nil {
		return err
	}

	localGUIDs := map[time.Time]string{}
	for _, item := range local.Items {
		if item.PublishedParsed == nil {
			log.Warningf("Item \"%s\" in local feed has no published date. Skipping.", item.Title)
			continue
		}
		localGUIDs[*item.PublishedParsed] = item.GUID
	}
	if len(localGUIDs) != len(local.Items) {
		log.Errorf("Local feed has %d items, but only %d unique publish times were found.", len(local.Items), len(localGUIDs))
	}

	remoteGUIDs := map[time.Time]string{}
	for _, item := range remote.Items {
		if item.PublishedParsed == nil {
			log.Warningf("Item \"%s\" in remote feed has no published date. Skipping.", item.Title)
			continue
		}
		remoteGUIDs[*item.PublishedParsed] = item.GUID
	}
	if len(remoteGUIDs) != len(remote.Items) {
		log.Errorf("Remote feed has %d items, but only %d unique publish times were found.", len(remote.Items), len(remoteGUIDs))
	}

	combinedGUIDs := map[time.Time]string{}
	for k, v := range localGUIDs {
		combinedGUIDs[k] = v
	}
	for k, v := range remoteGUIDs {
		combinedGUIDs[k] = v
	}

	for date := range combinedGUIDs {
		var localGUID, remoteGUID string
		var ok bool
		if localGUID, ok = localGUIDs[date]; !ok {
			log.Warningf("Item with date %s is not in local feed. GUID of remote item is \"\".", date, remoteGUIDs[date])
			continue
		}
		if remoteGUID, ok = remoteGUIDs[date]; !ok {
			log.Warningf("Item with date %s is not in remote feed. GUID of local item is \"\".", date, localGUIDs[date])
			continue
		}

		if localGUID != remoteGUID {
			log.Errorf("Item with date %s has different GUIDs in local and remote feeds. Local: \"%s\", remote: \"%s\".", date, localGUID, remoteGUID)
		}
	}

	return nil
}
