/*
   file created by Junlin Chen in 2022

*/

package client

import (
	"context"
	"net/url"
	"testing"
)

func TestStarlightProxy_Ping(t *testing.T) {
	t.Skip("use your own server")

	proxy := NewStarlightProxy(context.TODO(), "https", "test.yuri.moe")
	proxy.auth = *url.UserPassword("username", "password")

	if err := proxy.Ping(); err != nil {
		t.Error(err)
	}
}
