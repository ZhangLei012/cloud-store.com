package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	"cloud-store.com/service/dbproxy/mapper"
	"cloud-store.com/service/dbproxy/orm"
	dbproto "cloud-store.com/service/dbproxy/proto"
)

type DBProxy struct{}

func (db *DBProxy) ExecuteAction(ctx context.Context, req *dbproto.ReqExec, res *dbproto.RespExec) error {
	resList := make([]orm.ExecResult, len(req.Action))
	//TODO 检查req.Sequence req.Transaction两个参数，执行不同的流程
	for i, singleAction := range req.Action {
		params := []interface{}{}
		dec := json.NewDecoder(bytes.NewReader(singleAction.Params))
		dec.UseNumber()
		if err := dec.Decode(&params); err != nil {
			resList[i] = orm.ExecResult{
				Success: false,
				Msg:     "request parameters exist problem(s)",
			}
			continue
		}

		for k, v := range params {
			if _, ok := v.(json.Number); ok {
				params[k], _ = v.(json.Number).Int64()
			}
		}

		//默认串行执行sql函数
		execRes, err := mapper.FuncCall(singleAction.Name, params...)
		if err != nil {
			resList[i] = orm.ExecResult{
				Success: false,
				Msg:     err.Error(),
			}
			continue
		}
		resList[i] = execRes[0].Interface().(orm.ExecResult)
	}

	//TODO 处理异常
	res.Data, _ = json.Marshal(resList)
	return nil
}
