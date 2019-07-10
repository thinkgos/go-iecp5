package asdu

//
//func (d *Delegate) singleCmd(req *asdu.ASDU, c *Caller) {
//	addr := req.DecodeInfoObjAddr(req.InfoObj)
//	f, ok := d.SingleCmds[addr]
//	if !ok {
//		f, ok = d.SingleCmds[asdu.InfoObjIrrelevantAddr]
//	}
//	if !ok {
//		c.Send(req.Reply(asdu.UnkInfo))
//		return
//	}
//
//	cmd := asdu.Command(req.InfoObj[req.InfoObjAddrSize])
//	attrs := ExecAttrs{
//		InfoObjAddr:  addr,
//		CmdQualifier: cmd.Qual(),
//		InSelect:     !cmd.Exec(),
//	}
//	if req.Type == asdu.C_SC_TA_1 {
//		// TODO(pascaldekloe): time tag
//	}
//
//	req = selectProc(req, req.InfoObjAddrSize, c)
//	if req == nil {
//		return
//	}
//
//	d.Add(1)
//	go func() {
//		defer d.Done()
//
//		terminate := func() {
//			c.Send(req.Reply(asdu.Actterm))
//		}
//		p := asdu.SinglePoint(cmd & 1)
//		if f(req.Identifier, p, attrs, terminate) {
//			c.Send(req.Reply(asdu.Actcon))
//		} else {
//			c.Send(req.Reply(asdu.Actcon | asdu.NegFlag))
//		}
//	}()
//}
//
//func (d *Delegate) doubleCmd(req *asdu.ASDU, c *Caller) {
//	addr := req.DecodeInfoObjAddr(req.InfoObj)
//	f, ok := d.DoubleCmds[addr]
//	if !ok {
//		f, ok = d.DoubleCmds[asdu.InfoObjIrrelevantAddr]
//	}
//	if !ok {
//		c.Send(req.Reply(asdu.UnkInfo))
//		return
//	}
//
//	cmd := asdu.Command(req.InfoObj[req.InfoObjAddrSize])
//	attrs := ExecAttrs{
//		InfoObjAddr:  addr,
//		CmdQualifier: cmd.Qual(),
//		InSelect:     !cmd.Exec(),
//	}
//	if req.Type == asdu.C_DC_TA_1 {
//		// TODO(pascaldekloe): time tag
//	}
//
//	req = selectProc(req, req.InfoObjAddrSize, c)
//	if req == nil {
//		return
//	}
//
//	d.Add(1)
//	go func() {
//		defer d.Done()
//
//		terminate := func() {
//			c.Send(req.Reply(asdu.Actterm))
//		}
//		p := asdu.DoublePoint(cmd & 3)
//		if f(req.Identifier, p, attrs, terminate) {
//			c.Send(req.Reply(asdu.Actcon))
//		} else {
//			c.Send(req.Reply(asdu.Actcon | asdu.NegFlag))
//		}
//	}()
//}
//
//func (d *Delegate) floatSetpoint(req *asdu.ASDU, c *Caller) {
//	addr := req.DecodeInfoObjAddr(req.InfoObj)
//	f, ok := d.FloatSetpoints[addr]
//	if !ok {
//		f, ok = d.FloatSetpoints[asdu.InfoObjIrrelevantAddr]
//	}
//	if !ok {
//		c.Send(req.Reply(asdu.UnkInfo))
//		return
//	}
//
//	cmd := asdu.SetPointCmd(req.InfoObj[req.InfoObjAddrSize+4])
//	attrs := ExecAttrs{
//		InfoObjAddr:  addr,
//		CmdQualifier: cmd.Qual(),
//		InSelect:     !cmd.Exec(),
//	}
//	if req.Type == asdu.C_SE_TC_1 {
//		// TODO(pascaldekloe): time tag
//	}
//
//	req = selectProc(req, req.InfoObjAddrSize+4, c)
//	if req == nil {
//		return
//	}
//
//	d.Add(1)
//	go func() {
//		defer d.Done()
//
//		terminate := func() {
//			c.Send(req.Reply(asdu.Actterm))
//		}
//		p := math.Float32frombits(binary.LittleEndian.Uint32(req.InfoObj[req.InfoObjAddrSize:]))
//		if f(req.Identifier, p, attrs, terminate) {
//			c.Send(req.Reply(asdu.Actcon))
//		} else {
//			c.Send(req.Reply(asdu.Actcon | asdu.NegFlag))
//		}
//	}()
//}
