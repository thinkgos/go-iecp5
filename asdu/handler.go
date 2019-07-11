package asdu

// 服务器端解析
func ServerHandler(c Connect, frame []byte) error {
	asdu := NewEmptyASDU(c.Params())
	if err := asdu.UnmarshalBinary(frame); err != nil {
		return err
	}

	identifier := asdu.Identifier

	switch identifier.Type {

	}
	return nil
}
