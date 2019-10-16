// +build oracle

package sql

//CreateAccount 创建帐户信息
const CreateAccount = `INSERT INTO beanpay_account_info(account_id,account_name,
	eid,account_type,balance,credit,status,create_time)values(seq_account_info_id.nextval,
	@name,@eid,@account_type,0,0,0,sysdate)`

//GetAccountByeid 根据eid查询帐户编号
const GetAccountByeid = `select t.account_id,t.account_name,t.eid,t.balance,t.credit from beanpay_account_info t where t.eid=@eid`

//ChangeAmount 帐户加款
const ChangeAmount = `update beanpay_account_info t set t.balance=t.balance + @amount where t.account_id=@account_id
and t.balance + @amount >= 0`

//ExistsBalanceRecord 查询交易变动记录是否已存在
const ExistsBalanceRecord = `select count(0) from beanpay_account_record t 
where t.trade_no=@trade_no 
and t.account_id=@account_id
and t.change_type=@tp
and abs(t.amount) >= @max_amount
`

//GetBalanceRecord 查询交易变动记录是否已存在
const GetBalanceRecord = `select  t.record_id,t.account_id,
t.trade_no,t.change_type,t.amount,t.balance,to_char(t.create_time, 'yyyymmddhh24miss') create_time from beanpay_account_record t 
where t.record_id=@record_id`

//GetBalanceRecordByTradeNo 查询交易变动记录是否已存在
const GetBalanceRecordByTradeNo = `select t.record_id,t.account_id,
t.trade_no,t.change_type,t.amount,t.balance,to_char(t.create_time, 'yyyymmddhh24miss') create_time
 from beanpay_account_record t 
where t.trade_no=@trade_no and t.account_id=@account_id
and t.change_type=@tp`

//AddBalanceRecord 添加资金变动
const AddBalanceRecord = `insert into beanpay_account_record
(record_id,account_id,trade_no,change_type,amount,balance,create_time)
select seq_account_record_id.nextval,@account_id,@trade_no,@tp,@amount,t.balance,sysdate
 from beanpay_account_info t where t.account_id=@account_id`

//QueryBalanceRecord 查询余额资金变动信息
const QueryBalanceRecord = `select 
l2.record_id,l2.account_id,
l2.trade_no,l2.change_type,l2.amount,l2.balance,l2.create_time
from(select rownum rn,l1.* from(select t.record_id,t.account_id,
t.trade_no,t.change_type,t.amount,t.balance,to_char(t.create_time, 'yyyymmddhh24miss') create_time from  
beanpay_account_record t where t.account_id = @account_id and 
t.create_time >= to_date(@start,'yyyymmdd')
and t.create_time < to_date(@end,'yyyymmdd')+1
order by t.record_id desc) l1
where rownum <= (@pi+1) * @ps) l2 
where l2.rn > (@pi) * @ps`
