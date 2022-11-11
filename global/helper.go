package global

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	cryptoEth "github.com/ethereum/go-ethereum/crypto"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ipfsGateway []string
var ipfsUploadMutex sync.Mutex

func FormatEthSignMsg(msg string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len([]byte(msg)), msg)
}

func EthSignHash(msg string) []byte {
	msg = FormatEthSignMsg(msg)
	return cryptoEth.Keccak256([]byte(msg))
}

func ReloadIpfsGateway() error {
	r, err := http.Get("https://raw.githubusercontent.com/ipfs/public-gateway-checker/master/src/gateways.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &ipfsGateway)
}

func GetIpfsInfo(h string) ([]byte, error) {
	var err error
	if len(ipfsGateway) == 0 {
		if err = ReloadIpfsGateway(); err != nil {
			return nil, err
		}
	}
	hc := &http.Client{
		Timeout: 2 * time.Second,
	}
	for _, gateway := range ipfsGateway {
		r, err2 := hc.Get(strings.ReplaceAll(gateway, ":hash", h))
		if err2 != nil {
			err = err2
			continue
		}
		b, err2 := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		if err2 != nil {
			err = err2
			continue
		}
		return b, nil
	}
	return nil, err
}

func UploadIpfs(data interface{}) (string, error) {
	log.Println("upload ipfs start .....")
	defer log.Println("upload ipfs ok .....")
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	ipfsUploadMutex.Lock()
	defer ipfsUploadMutex.Unlock()
	sh := shell.NewShell("https://cdn.ipfsscan.io")
	return sh.Add(bytes.NewReader(b))
}

func UploadIpfsReader(r *bytes.Reader) (string, error) {
	ipfsUploadMutex.Lock()
	defer ipfsUploadMutex.Unlock()
	sh := shell.NewShell("https://cdn.ipfsscan.io")
	return sh.Add(r)
}

func IsLocalIp(ip string) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		intf, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if net.ParseIP(ip).Equal(intf) {
			return true
		}
	}
	return false
}

func StrLen(str string) int {
	return strings.Count(str, "") - 1
}

func PwdPlaintext2CipherText(pwd string, salt string) string {
	pwd = salt + "{_}" + pwd + "{_}" + salt
	has := md5.Sum([]byte(pwd))
	return fmt.Sprintf("%x", has)
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func LibP2pPriToEthPri(pri crypto.PrivKey) (*ecdsa.PrivateKey, error) {
	priBytes, err := pri.Raw()
	if err != nil {
		return nil, err
	}
	return cryptoEth.HexToECDSA(hex.EncodeToString(priBytes))
}

func EthPriToLibP2pPri(pri *ecdsa.PrivateKey) (crypto.PrivKey, error) {
	pri0, _, err := crypto.ECDSAKeyPairFromKey(pri)
	return pri0, err
}

func Base58ToPubKey(str string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := base58.Decode(str)
	if err != nil {
		return nil, err
	}
	return cryptoEth.UnmarshalPubkey(pubKeyBytes)
}

func LibP2pPriToAddress(pri crypto.PrivKey) (common.Address, error) {
	ethPri, err := LibP2pPriToEthPri(pri)
	if err != nil {
		return common.Address{}, err
	}
	return cryptoEth.PubkeyToAddress(ethPri.PublicKey), nil
}

func GenerateRangeNum(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max-min) + min
	return randNum
}

func Hour2Unix(hour string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormatStr, time.Now().Format(DateFormatStr)+" "+hour, time.Local)
}

func Md5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func Json2Map(j string) map[string]interface{} {
	r := make(map[string]interface{})
	_ = json.Unmarshal([]byte(j), &r)
	return r
}

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func RandFloats(min, max float64, n int) float64 {
	rand.Seed(time.Now().UnixNano())
	res := min + rand.Float64()*(max-min)
	res, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(n)+"f", res), 64)
	return res
}

func RemoveDuplicationByMap(arr []string, before func(string2 *string)) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		if before != nil {
			before(&v)
		}
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}
	return arr[:j]
}
