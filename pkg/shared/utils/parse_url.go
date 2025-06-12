package utils

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type URL struct {
	Subdomain, Domain, TLD, Port string
	ICANN                        bool
	*url.URL
}

func ParseUrl(s string) (*URL, error) {
	urlParsed, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	// working with url without scheme
	// example: www.example.com
	if urlParsed.Host == "" {
		urlParsed, err = url.Parse("https://" + s)
		if err != nil {
			return nil, err
		}
		if urlParsed.Host == "" {
			return &URL{URL: urlParsed}, nil
		}
	}
	dom, port := domainPort(urlParsed.Host)
	etld1, err := publicsuffix.EffectiveTLDPlusOne(dom)
	_, icann := publicsuffix.PublicSuffix(strings.ToLower(dom))
	if err != nil {
		return nil, err
	}
	// convert to domain name, and tld
	i := strings.Index(etld1, ".")
	if i < 0 {
		return nil, fmt.Errorf("tld: failed parsing %q", s)
	}
	domName := etld1[0:i]
	tld := etld1[i+1:]
	// and subdomain
	sub := ""
	if rest := strings.TrimSuffix(dom, "."+etld1); rest != dom {
		sub = rest
	}
	return &URL{
		Subdomain: sub,
		Domain:    domName,
		TLD:       tld,
		Port:      port,
		ICANN:     icann,
		URL:       urlParsed,
	}, nil
}

func domainPort(host string) (string, string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			return host, ""
		}
	}
	return host, ""
}

// IsNotCompatible comment.
// en: return true if url is not compatible with target url else false.
func (u *URL) IsNotCompatible(target *URL) bool {
	if target == nil || u == nil {
		return true
	}
	return u.Domain != target.Domain || u.TLD != target.TLD || !u.ICANN || !target.ICANN
}
