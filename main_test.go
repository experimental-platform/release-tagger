package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/experimental-platform/release-tagger/git"

	"gopkg.in/stretchr/testify.v1/assert"
)

var testOldJSON = `[
  {
    "build": 213455,
    "codename": "Development Alpha",
    "url": "foobar",
    "published_at": "2016-08-24T14:02:38Z",
    "images": {
      "quay.io/experimentalplatform/afpd": "2016-08-24-1402",
      "quay.io/experimentalplatform/app-manager": "2016-08-24-1402",
      "quay.io/experimentalplatform/central-gateway": "2016-08-24-1402",
      "quay.io/experimentalplatform/collectd": "2016-08-24-1402",
      "quay.io/experimentalplatform/configure": "2016-08-29-1453",
      "quay.io/experimentalplatform/dnsmasq": "2016-08-24-1402",
      "quay.io/experimentalplatform/dokku": "2016-08-24-1402",
      "quay.io/experimentalplatform/elasticsearch": "2016-08-24-1402",
      "quay.io/experimentalplatform/frontend": "2016-08-24-1402",
      "quay.io/experimentalplatform/haproxy": "2016-08-24-1402",
      "quay.io/experimentalplatform/hardware": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostapd": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostname-avahi": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostname-smb": "2016-08-24-1402",
      "quay.io/experimentalplatform/http-proxy": "2016-08-24-1402",
      "quay.io/experimentalplatform/ldap": "2016-08-24-1402",
      "quay.io/experimentalplatform/monitoring": "2016-08-24-1402",
      "quay.io/experimentalplatform/mysql": "2016-08-24-1402",
      "quay.io/experimentalplatform/ptw": "2016-08-24-1402",
      "quay.io/experimentalplatform/pulseaudio": "2016-08-24-1402",
      "quay.io/experimentalplatform/rabbitmq": "2016-08-24-1402",
      "quay.io/experimentalplatform/redis": "2016-08-24-1402",
      "quay.io/experimentalplatform/skvs": "2016-08-24-1402",
      "quay.io/experimentalplatform/smb": "2016-08-24-1402",
      "quay.io/experimentalplatform/systemd-proxy": "2016-08-24-1402",
      "quay.io/protonetinc/german-shepherd": "2016-08-24-1402",
      "quay.io/protonetinc/soul-backup": "2016-08-24-1402",
      "quay.io/protonetinc/soul-nginx": "2016-08-24-1402",
      "quay.io/protonetinc/soul-owner": "2016-08-24-1402",
      "quay.io/protonetinc/soul-protosync": "2016-08-24-1402",
      "quay.io/protonetinc/soul-smb": "2016-08-24-1402"
    }
  }
]`

// TestRenamedImages tests whether the image list
// contains the same images with altered tags
func TestRenamedImages(t *testing.T) {
	repo, err := git.PrepareRepo("libgit")
	assert.Nil(t, err)
	defer repo.Close()

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	ioutil.WriteFile(srcJSONPath, []byte(testOldJSON), 0644)

	opts := taggerOptions{
		Commit:   false,
		Build:    0,
		URL:      "",
		Codename: "",
		Args: taggerOptionsArgs{
			Action:        "create",
			SourceChannel: "source",
			TargetChannel: "tgt",
		},
	}
	tagTimestamp := "tag-timestamp #124124"
	isoTimestamp := "wtf_timestamp %3215123"
	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	assert.Nil(t, err)

	var expectedJSON = `[
  {
    "build": 1,
    "codename": "Development Alpha",
    "url": "foobar",
    "published_at": "wtf_timestamp %3215123",
    "images": {
      "quay.io/experimentalplatform/afpd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/app-manager": "tag-timestamp #124124",
      "quay.io/experimentalplatform/central-gateway": "tag-timestamp #124124",
      "quay.io/experimentalplatform/collectd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/configure": "tag-timestamp #124124",
      "quay.io/experimentalplatform/dnsmasq": "tag-timestamp #124124",
      "quay.io/experimentalplatform/dokku": "tag-timestamp #124124",
      "quay.io/experimentalplatform/elasticsearch": "tag-timestamp #124124",
      "quay.io/experimentalplatform/frontend": "tag-timestamp #124124",
      "quay.io/experimentalplatform/haproxy": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hardware": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostapd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostname-avahi": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostname-smb": "tag-timestamp #124124",
      "quay.io/experimentalplatform/http-proxy": "tag-timestamp #124124",
      "quay.io/experimentalplatform/ldap": "tag-timestamp #124124",
      "quay.io/experimentalplatform/monitoring": "tag-timestamp #124124",
      "quay.io/experimentalplatform/mysql": "tag-timestamp #124124",
      "quay.io/experimentalplatform/ptw": "tag-timestamp #124124",
      "quay.io/experimentalplatform/pulseaudio": "tag-timestamp #124124",
      "quay.io/experimentalplatform/rabbitmq": "tag-timestamp #124124",
      "quay.io/experimentalplatform/redis": "tag-timestamp #124124",
      "quay.io/experimentalplatform/skvs": "tag-timestamp #124124",
      "quay.io/experimentalplatform/smb": "tag-timestamp #124124",
      "quay.io/experimentalplatform/systemd-proxy": "tag-timestamp #124124",
      "quay.io/protonetinc/german-shepherd": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-backup": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-nginx": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-owner": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-protosync": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-smb": "tag-timestamp #124124"
    }
  }
]`

	actualJSON, err := repo.DumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}

// TestRenamedImages2 tests whether the Codename and URL have been updated
func TestRenamedImages2(t *testing.T) {
	repo, err := git.PrepareRepo("libgit")
	assert.Nil(t, err)
	defer repo.Close()

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	ioutil.WriteFile(srcJSONPath, []byte(testOldJSON), 0644)

	opts := taggerOptions{
		Commit:   false,
		Build:    9875,
		URL:      "https://www.example.com/",
		Codename: "Zeitgeist",
		Args: taggerOptionsArgs{
			SourceChannel: "source",
			TargetChannel: "tgt",
			Action:        "create",
		},
	}
	tagTimestamp := "tag-timestamp #124124"
	isoTimestamp := "wtf_timestamp %3215123"
	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	assert.Nil(t, err)

	var expectedJSON = `[
  {
    "build": 9875,
    "codename": "Zeitgeist",
    "url": "https://www.example.com/",
    "published_at": "wtf_timestamp %3215123",
    "images": {
      "quay.io/experimentalplatform/afpd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/app-manager": "tag-timestamp #124124",
      "quay.io/experimentalplatform/central-gateway": "tag-timestamp #124124",
      "quay.io/experimentalplatform/collectd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/configure": "tag-timestamp #124124",
      "quay.io/experimentalplatform/dnsmasq": "tag-timestamp #124124",
      "quay.io/experimentalplatform/dokku": "tag-timestamp #124124",
      "quay.io/experimentalplatform/elasticsearch": "tag-timestamp #124124",
      "quay.io/experimentalplatform/frontend": "tag-timestamp #124124",
      "quay.io/experimentalplatform/haproxy": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hardware": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostapd": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostname-avahi": "tag-timestamp #124124",
      "quay.io/experimentalplatform/hostname-smb": "tag-timestamp #124124",
      "quay.io/experimentalplatform/http-proxy": "tag-timestamp #124124",
      "quay.io/experimentalplatform/ldap": "tag-timestamp #124124",
      "quay.io/experimentalplatform/monitoring": "tag-timestamp #124124",
      "quay.io/experimentalplatform/mysql": "tag-timestamp #124124",
      "quay.io/experimentalplatform/ptw": "tag-timestamp #124124",
      "quay.io/experimentalplatform/pulseaudio": "tag-timestamp #124124",
      "quay.io/experimentalplatform/rabbitmq": "tag-timestamp #124124",
      "quay.io/experimentalplatform/redis": "tag-timestamp #124124",
      "quay.io/experimentalplatform/skvs": "tag-timestamp #124124",
      "quay.io/experimentalplatform/smb": "tag-timestamp #124124",
      "quay.io/experimentalplatform/systemd-proxy": "tag-timestamp #124124",
      "quay.io/protonetinc/german-shepherd": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-backup": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-nginx": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-owner": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-protosync": "tag-timestamp #124124",
      "quay.io/protonetinc/soul-smb": "tag-timestamp #124124"
    }
  }
]`

	actualJSON, err := repo.DumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}

// TestRenamedImages3 tests whether the image list
// contains the same images with unchanged tags
func TestRenamedImages3(t *testing.T) {
	repo, err := git.PrepareRepo("libgit")
	assert.Nil(t, err)
	defer repo.Close()

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	ioutil.WriteFile(srcJSONPath, []byte(testOldJSON), 0644)

	opts := taggerOptions{
		Commit:   false,
		Build:    0,
		URL:      "",
		Codename: "",
		Args: taggerOptionsArgs{
			Action:        "copy",
			SourceChannel: "source",
			TargetChannel: "tgt",
		},
	}
	tagTimestamp := "tag-timestamp #124124"
	isoTimestamp := "wtf_timestamp %3215123"
	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	assert.Nil(t, err)

	var expectedJSON = `[
  {
    "build": 1,
    "codename": "Development Alpha",
    "url": "foobar",
    "published_at": "wtf_timestamp %3215123",
    "images": {
      "quay.io/experimentalplatform/afpd": "2016-08-24-1402",
      "quay.io/experimentalplatform/app-manager": "2016-08-24-1402",
      "quay.io/experimentalplatform/central-gateway": "2016-08-24-1402",
      "quay.io/experimentalplatform/collectd": "2016-08-24-1402",
      "quay.io/experimentalplatform/configure": "2016-08-29-1453",
      "quay.io/experimentalplatform/dnsmasq": "2016-08-24-1402",
      "quay.io/experimentalplatform/dokku": "2016-08-24-1402",
      "quay.io/experimentalplatform/elasticsearch": "2016-08-24-1402",
      "quay.io/experimentalplatform/frontend": "2016-08-24-1402",
      "quay.io/experimentalplatform/haproxy": "2016-08-24-1402",
      "quay.io/experimentalplatform/hardware": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostapd": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostname-avahi": "2016-08-24-1402",
      "quay.io/experimentalplatform/hostname-smb": "2016-08-24-1402",
      "quay.io/experimentalplatform/http-proxy": "2016-08-24-1402",
      "quay.io/experimentalplatform/ldap": "2016-08-24-1402",
      "quay.io/experimentalplatform/monitoring": "2016-08-24-1402",
      "quay.io/experimentalplatform/mysql": "2016-08-24-1402",
      "quay.io/experimentalplatform/ptw": "2016-08-24-1402",
      "quay.io/experimentalplatform/pulseaudio": "2016-08-24-1402",
      "quay.io/experimentalplatform/rabbitmq": "2016-08-24-1402",
      "quay.io/experimentalplatform/redis": "2016-08-24-1402",
      "quay.io/experimentalplatform/skvs": "2016-08-24-1402",
      "quay.io/experimentalplatform/smb": "2016-08-24-1402",
      "quay.io/experimentalplatform/systemd-proxy": "2016-08-24-1402",
      "quay.io/protonetinc/german-shepherd": "2016-08-24-1402",
      "quay.io/protonetinc/soul-backup": "2016-08-24-1402",
      "quay.io/protonetinc/soul-nginx": "2016-08-24-1402",
      "quay.io/protonetinc/soul-owner": "2016-08-24-1402",
      "quay.io/protonetinc/soul-protosync": "2016-08-24-1402",
      "quay.io/protonetinc/soul-smb": "2016-08-24-1402"
    }
  }
]`

	actualJSON, err := repo.DumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}

func TestRenamedImagesBuildIncrement(t *testing.T) {
	repo, err := git.PrepareRepo("libgit")
	assert.Nil(t, err)
	defer repo.Close()

	var oldJSON1 = `[
  {
    "build": 213455,
    "codename": "Development Alpha",
    "url": "foobar",
    "published_at": "2016-08-24T14:02:38Z",
    "images": {}
  }
]`

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	tgtJSONPath := path.Join(repo.GetDirectory(), "tgt.json")
	ioutil.WriteFile(srcJSONPath, []byte(oldJSON1), 0644)
	ioutil.WriteFile(tgtJSONPath, []byte(oldJSON1), 0644)

	opts := taggerOptions{
		Commit:   false,
		Build:    0,
		URL:      "https://www.example.com/",
		Codename: "Zeitgeist",
		Args: taggerOptionsArgs{
			SourceChannel: "source",
			TargetChannel: "tgt",
			Action:        "create",
		},
	}
	tagTimestamp := "tag-timestamp #124124"
	isoTimestamp := "wtf_timestamp %3215123"
	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	assert.Nil(t, err)

	var expectedJSON = `[
  {
    "build": 213456,
    "codename": "Zeitgeist",
    "url": "https://www.example.com/",
    "published_at": "wtf_timestamp %3215123",
    "images": {}
  }
]`

	actualJSON, err := repo.DumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}
