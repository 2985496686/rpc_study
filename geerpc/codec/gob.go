package codec

import (
	"bufio"
	"encoding/gob"
	"geerpc/log"
	"io"
)

type GobCodec struct {
	conn io.ReadWriteCloser //网络连接实例
	buf  *bufio.Writer      //写缓冲
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// ReadHeader 将conn中的数据解码到 *Header中
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Error("rpc codec: gob error encoding header:", err)
		return err
	}
	if err = c.enc.Encode(body); err != nil {
		log.Error("rpc codec: gob error encoding body:", err)
		return err
	}
	return
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
