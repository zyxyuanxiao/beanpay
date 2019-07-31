package account

import (
	"github.com/micro-plat/beanpay/beanpay/const/ecodes"
	"github.com/micro-plat/beanpay/beanpay/const/ttypes"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/db"
)

//Create 根据uid,name创建帐户,如果帐户存在直接返回帐户编号
func Create(db db.IDBExecuter, uid string, name string) (interface{}, error) {
	acc, err := GetAccount(db, uid)
	if err == nil {
		return context.NewResult(ecodes.HasExists, acc), nil
	}
	if context.GetCode(err) != ecodes.NotExists {
		return nil, err
	}
	if err = create(db, uid, name); err != nil {
		return nil, err
	}
	return GetAccount(db, uid)
}

//GetBalance 获取帐户余额
func GetBalance(db db.IDBExecuter, uid string) (int, error) {
	return getBalance(db, uid)
}

//GetAccount 根据uid获取帐户
func GetAccount(db db.IDBExecuter, uid string) (acc *Account, err error) {
	row, err := getAccount(db, uid)
	if err != nil {
		return nil, err
	}
	acc = &Account{}
	if err = row.ToStruct(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

//AddAmount 资金加款
func AddAmount(db db.IDBExecuter, uid string, tradeNo string, amount int) (*context.Result, error) {
	if amount <= 0 {
		return nil, context.NewErrorf(ecodes.AmountErr, "金额错误%d", amount)
	}
	acc, err := GetAccount(db, uid)
	if err != nil {
		return nil, err
	}

	b, err := exists(db, acc.ID, tradeNo, 0, ttypes.Add)
	if err != nil {
		return nil, err
	}
	if b {
		row, err := getRecordByTradeNo(db, acc.ID, tradeNo, ttypes.Add)
		if err != nil {
			return nil, context.NewError(ecodes.Failed, "暂时无法加款")
		}
		return context.NewResult(ecodes.HasExists, row), nil
	}
	row, err := change(db, acc.ID, tradeNo, ttypes.Add, amount)
	if err != nil {
		return nil, err
	}
	return context.NewResult(ecodes.Success, row), nil
}

//DeductAmount 资金扣款
func DeductAmount(db db.IDBExecuter, uid string, tradeNo string, amount int) (*context.Result, error) {
	if amount <= 0 {
		return nil, context.NewErrorf(ecodes.AmountErr, "金额错误%d", amount)
	}
	acc, err := GetAccount(db, uid)
	if err != nil {
		return nil, err
	}
	b, err := exists(db, acc.ID, tradeNo, 0, ttypes.Deduct)
	if err != nil {
		return nil, err
	}
	if b {
		row, err := getRecordByTradeNo(db, acc.ID, tradeNo, ttypes.Add)
		if err != nil {
			return nil, context.NewError(ecodes.Failed, "暂时无法扣款")
		}
		return context.NewResult(ecodes.HasExists, row), nil
	}
	row, err := change(db, acc.ID, tradeNo, ttypes.Deduct, -amount)
	if err != nil {
		return nil, err
	}
	return context.NewResult(ecodes.Success, row), nil
}

//RefundAmount 资金退款
func RefundAmount(db db.IDBExecuter, uid string, tradeNo string, amount int) (*context.Result, error) {
	if amount <= 0 {
		return nil, context.NewErrorf(ecodes.AmountErr, "金额错误%d", amount)
	}
	acc, err := GetAccount(db, uid)
	if err != nil {
		return nil, err
	}
	//检查是否已退款
	b, err := exists(db, acc.ID, tradeNo, amount, ttypes.Refund)
	if err != nil {
		return nil, err
	}
	if b {
		row, err := getRecordByTradeNo(db, acc.ID, tradeNo, ttypes.Refund)
		if err != nil {
			return nil, context.NewError(ecodes.Failed, "暂时无法退款")
		}
		return context.NewResult(ecodes.HasExists, row), nil
	}
	//检查是否存在加款记录
	b, err = exists(db, acc.ID, tradeNo, amount, ttypes.Add)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, context.NewError(ecodes.HasExists, "加款交易编号不存在")
	}
	row, err := change(db, acc.ID, tradeNo, ttypes.Refund, amount)
	if err != nil {
		return nil, err
	}
	return context.NewResult(ecodes.Success, row), nil
}

//Query 查询余额变动明细
func Query(db db.IDBExecuter, uid string, startTime string, endTime string, pi int, ps int) (db.QueryRows, error) {
	acc, err := GetAccount(db, uid)
	if err != nil {
		return nil, err
	}
	return query(db, acc.ID, startTime, endTime, pi, ps)
}
