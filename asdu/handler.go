package asdu

// 服务器端解析
func ServerHandler(c Connect, frame []byte) error {
	adu := NewEmptyASDU(c.Params())
	if err := adu.UnmarshalBinary(frame); err != nil {
		return err
	}

	return nil
}
