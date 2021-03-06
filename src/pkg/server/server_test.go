package server

import (
	"bytes"
	"doozer/store"
	"github.com/bmizerany/assert"
	"goprotobuf.googlecode.com/hg/proto"
	"os"
	"testing"
)


var (
	fooPath = "/foo"
)


type bchan chan []byte


func (b bchan) Write(buf []byte) (int, os.Error) {
	b <- buf
	return len(buf), nil
}


func (b bchan) Read(buf []byte) (int, os.Error) {
	return 0, os.EOF // not implemented
}


func mustUnmarshal(b []byte) (r *response) {
	r = new(response)
	err := proto.Unmarshal(b, r)
	if err != nil {
		panic(err)
	}
	return
}


func assertResponseErrCode(t *testing.T, exp response_Err, c *conn) {
	b := c.c.(*bytes.Buffer).Bytes()
	assert.T(t, len(b) > 4, b)
	assert.Equal(t, &exp, mustUnmarshal(b[4:]).ErrCode)
}


func TestDelNilFields(t *testing.T) {
	c := &conn{
		c:        &bytes.Buffer{},
		canWrite: true,
		access:   true,
	}
	tx := &txn{
		c:   c,
		req: request{Tag: proto.Int32(1)},
	}
	tx.del()
	assertResponseErrCode(t, response_MISSING_ARG, c)
}


func TestSetNilFields(t *testing.T) {
	c := &conn{
		c:        &bytes.Buffer{},
		canWrite: true,
		access:   true,
	}
	tx := &txn{
		c:   c,
		req: request{Tag: proto.Int32(1)},
	}
	tx.set()
	assertResponseErrCode(t, response_MISSING_ARG, c)
}


func TestServerNoAccess(t *testing.T) {
	b := make(bchan, 2)
	c := &conn{
		c:        b,
		canWrite: true,
		st:       store.New(),
	}
	tx := &txn{
		c:   c,
		req: request{Tag: proto.Int32(1)},
	}

	for i, op := range ops {
		if i != request_ACCESS {
			op(tx)
			var exp response_Err = response_OTHER
			assert.Equal(t, 4, len(<-b), request_Verb_name[i])
			assert.Equal(t, &exp, mustUnmarshal(<-b).ErrCode, request_Verb_name[i])
		}
	}
}

func TestServerInvalidAccessToken(t *testing.T) {
	b := make(bchan, 2)
	c := &conn{
		c:        b,
		canWrite: true,
		st:       store.New(),
		secret:   "abc",
	}
	tx := &txn{
		c:   c,
		req: request{Tag: proto.Int32(1), Value: []byte("bad")},
	}

	tx.access()

	var exp response_Err = response_OTHER
	assert.Equal(t, 4, len(<-b))
	assert.Equal(t, &exp, mustUnmarshal(<-b).ErrCode)

	assert.T(t, !c.access)
}

func TestServerValidAccessToken(t *testing.T) {
	b := make(bchan, 2)
	c := &conn{
		c:        b,
		canWrite: true,
		st:       store.New(),
		secret:   "abc",
	}
	tx := &txn{
		c:   c,
		req: request{Tag: proto.Int32(1), Value: []byte(c.secret)},
	}

	tx.access()

	assert.Equal(t, 4, len(<-b))
	assert.Equal(t, (*response_Err)(nil), mustUnmarshal(<-b).ErrCode)

	assert.T(t, c.access)
}
