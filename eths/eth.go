package eths

import (
	"fmt"
	"math/big"

	"go-copyright-p2/configs"
	"go-copyright-p2/utils"

	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//创建一个账户
func NewAcc(pass, connstr string) (string, error) {
	//连接到geth
	client, err := rpc.Dial(connstr)
	if err != nil {
		fmt.Println("failed to connect to geth", err)
		return "", err
	}

	//创建账户
	var account string
	err = client.Call(&account, "personal_newAccount", pass)
	if err != nil {
		fmt.Println("failed to create acc", err)
		return "", err
	}
	fmt.Println("acc==", account)
	return account, nil
}

func UploadPic(from, pass, hash, data string) error {
	///1. dial
	cli, err := ethclient.Dial(configs.Config.Eth.Connstr)
	if err != nil {
		fmt.Println("failed to conn to geth", err)
		return err
	}
	//2. 入口
	instance, err := NewPxa(common.HexToAddress(configs.Config.Eth.PxaAddr), cli)
	if err != nil {
		fmt.Println("failed to cNewCallabi", err)
		return err
	}
	//3. 设置签名 -- 需要owner的keystore文件
	//需要获得文件名字
	fileName, err := utils.GetFileName(string([]rune(from)[2:]), configs.Config.Eth.Keydir)
	if err != nil {
		fmt.Println("failed to GetFileName", err)
		return err
	}

	file, err := os.Open(configs.Config.Eth.Keydir + "/" + fileName)
	if err != nil {
		fmt.Println("failed to open file ", err)
		return err
	}
	auth, err := bind.NewTransactor(file, pass)
	if err != nil {
		fmt.Println("failed to NewTransactor  ", err)
		return err
	}
	//4. 调用
	//mint(string _hash, uint256 _price, uint256 _weight, string _data)
	_, err = instance.Mint(auth, hash, big.NewInt(100), big.NewInt(100), data)
	if err != nil {
		fmt.Println("failed to Mint  ", err)
		return err
	}
	return nil
}
