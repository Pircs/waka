package database

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

const (
	supervisor1 = 0.3
	supervisor2 = 0.09
	supervisor3 = 0.06
)

// 交易记录
type Transaction int32

// 交易记录数据
type TransactionData struct {
	// 主键
	Id Transaction `gorm:"index;unique;primary_key;AUTO_INCREMENT"`
	// 记录所属玩家
	Player Player
	// 交易对象
	Target Player
	// 数额
	Number int32
	// 类型
	// 0 未知
	// 1 支付
	// 2 收款
	Type int32
	// 原因
	Reason string
	// 一级代理
	Agent1 Player
	// 二级代理
	Agent2 Player
	// 三级代理
	Agent3 Player
	// 一级流水
	Money1 int32
	// 二级流水
	Money2 int32
	// 三级流水
	Money3 int32
	// 创建时间
	CreatedAt time.Time
	// 绑定id
	BindID string
}

func (TransactionData) TableName() string {
	return "transactions"
}

// 交易
type playerTransaction struct {
	// 原因
	Reason string
	// 支付者
	Payer Player
	// 收款人
	Payee Player
	// 交易数额
	Number int32
	// 折损率
	Loss float64
	// 是否支付提成
	EnableTip bool
}

func buildTransaction(modifies []*modifyMoneyAction, transaction *playerTransaction) []*modifyMoneyAction {
	if transaction.Payer < DefaultSupervisor ||
		transaction.Payee < DefaultSupervisor ||
		transaction.Number <= 0 ||
		(transaction.EnableTip && (transaction.Loss < 0 || transaction.Loss > 1)) {
		return modifies
	}
	uuid, _ := uuid.NewV4()
	modifies = append(modifies, &modifyMoneyAction{
		Player: transaction.Payer,
		Number: transaction.Number * (-1),
		After: func(ts *gorm.DB, self *modifyMoneyAction) error {
			if err := ts.Create(&TransactionData{
				Player:    transaction.Payer,
				Target:    transaction.Payee,
				Number:    transaction.Number,
				Type:      1,
				Reason:    transaction.Reason + ".pay",
				BindID:    uuid.String(),
				CreatedAt: time.Now(),
			}).Error; err != nil {
				return err
			}
			return nil
		},
	})
	if transaction.EnableTip {

		supervisorPlayer1 := transaction.Payer.PlayerData().Supervisor
		supervisorPlayer2 := supervisorPlayer1.PlayerData().Supervisor
		supervisorPlayer3 := supervisorPlayer2.PlayerData().Supervisor

		number := int32(float64(transaction.Number)*(1-transaction.Loss) + 0.5)
		supervisorNumber1 := int32(float64(transaction.Number-number)*supervisor1 + 0.5)
		supervisorNumber2 := int32(float64(transaction.Number-number)*supervisor2 + 0.5)
		supervisorNumber3 := int32(float64(transaction.Number-number)*supervisor3 + 0.5)
		systemNumber := transaction.Number - number - supervisorNumber1 - supervisorNumber2 - supervisorNumber3

		modifies = append(modifies, &modifyMoneyAction{
			Player: transaction.Payee,
			Number: number,
		})
		modifies = append(modifies, &modifyMoneyAction{
			Player: supervisorPlayer1,
			Number: supervisorNumber1,
			After: func(ts *gorm.DB, self *modifyMoneyAction) error {
				if err := ts.Create(&TransactionData{
					Player:    transaction.Payer,
					Number:    supervisorNumber1,
					Type:      2,
					Agent1:    supervisorPlayer1,
					Agent2:    supervisorPlayer2,
					Agent3:    supervisorPlayer3,
					Money1:    supervisorNumber1,
					Money2:    supervisorNumber2,
					Money3:    supervisorNumber3,
					BindID:    uuid.String(),
					Reason:    transaction.Reason + ".tip",
					CreatedAt: time.Now(),
				}).Error; err != nil {
					return err
				}
				return nil
			},
		})
		modifies = append(modifies, &modifyMoneyAction{
			Player: DefaultSupervisor,
			Number: systemNumber,
		})
	} else {
		modifies = append(modifies, &modifyMoneyAction{
			Player: transaction.Payee,
			Number: transaction.Number,
			After: func(ts *gorm.DB, self *modifyMoneyAction) error {
				if err := ts.Create(&TransactionData{
					Player:    transaction.Payee,
					Target:    transaction.Payer,
					Number:    transaction.Number,
					Type:      2,
					Reason:    transaction.Reason + ".income",
					CreatedAt: time.Now(),
				}).Error; err != nil {
					return err
				}
				return nil
			},
		})
	}

	return modifies
}
