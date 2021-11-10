package test

import "context"

type AuthToekn struct {
	Token string
}

func (c AuthToekn) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": c.Token,
	}, nil
}

func (c AuthToekn) RequireTransportSecurity() bool {
	return false
}
