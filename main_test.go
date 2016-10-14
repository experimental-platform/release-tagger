package main

import (
	"io/ioutil"
	"path"
	"testing"

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

func TestRenamedImages(t *testing.T) {
	repo, err := prepareRepo()
	assert.Nil(t, err)
	defer repo.Close()

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	ioutil.WriteFile(srcJSONPath, []byte(testOldJSON), 0644)

	opts := taggerOptions{
		Commit:        false,
		Build:         666,
		SourceChannel: "source",
		TargetChannel: "tgt",
		URL:           "",
		Codename:      "",
	}
	tagTimestamp := "tag-timestamp #124124"
	isoTimestamp := "wtf_timestamp %3215123"
	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	assert.Nil(t, err)

	var expectedJSON = `[
  {
    "build": 666,
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

	actualJSON, err := repo.dumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}

func TestRenamedImages2(t *testing.T) {
	repo, err := prepareRepo()
	assert.Nil(t, err)
	defer repo.Close()

	srcJSONPath := path.Join(repo.GetDirectory(), "source.json")
	ioutil.WriteFile(srcJSONPath, []byte(testOldJSON), 0644)

	opts := taggerOptions{
		Commit:        false,
		Build:         9875,
		SourceChannel: "source",
		TargetChannel: "tgt",
		URL:           "https://www.example.com/",
		Codename:      "Zeitgeist",
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

	actualJSON, err := repo.dumpChannel("tgt")
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}
