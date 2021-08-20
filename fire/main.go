package fire

import (
	"fmt"
	_"fmt"
	"encoding/json"
	"io/ioutil"
	"time"
)

type Poker struct {
	Alice  string `json:"alice"`
	Bob    string `json:"bob"`
	Result int    `json:"result"`
}

type Match struct {
	Matches []Poker `json:"matches"`
}


type CardCom struct {
	cardSizeMap1 map[byte]int
	cardSizeMap2 map[byte]int
}


func SizeTranByte(card byte)(res byte) {

	switch card {
	case 50:
		// 2
		res = 0x02
	case 51:
		res = 0x03
	case 52:
		res = 0x04
	case 53:
		res = 0x05
	case 54:
		res = 0x06
	case 55:
		res = 0x07
	case 56:
		res = 0x08
	case 57:
		res = 0x09
	case 84:
		res = 0x0A
	case 74:
		res = 0x0B
	case 81:
		res = 0x0C
	case 75:
		res = 0x0D
	case 65:
		res = 0x0E
	}
	return
}


// 把数据读取出来 分别放在切片中
func ReadFile(filename string) (alices, bobs []string, results []int) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var matches Match
	err = json.Unmarshal(buf, &matches)
	if err != nil {
		panic(err)
	}

	alices = make([]string, len(matches.Matches))
	bobs = make([]string, len(matches.Matches))
	results = make([]int, len(matches.Matches))

	for k, v := range matches.Matches {

		//fmt.Printf("k=%#v,v=%#v\n",k,v)
		alices[k] = v.Alice
		bobs[k] = v.Bob
		results[k] = v.Result
	}
	return
}



type CardType int

func JudgMentGroupNew(card []byte)(cardType CardType,cardSizeMap map[byte]int,resMax byte) {
	cardColorMap := make(map[byte]int,5)
	cardSizeMap = make(map[byte]int,5)
	// 扫描牌 分别放好大小，花色
	for i, v := range card {
		if i%2 ==0 {
			// 大小
			if _,ok := cardSizeMap[v]; ok {
				cardSizeMap[v] +=1
			}else {
				cardSizeMap[v] = 1
			}
			// 颜色
		}else {
			if _,ok := cardColorMap[v]; ok {
				cardColorMap[v] +=1
			}else {
				cardColorMap[v] = 1
			}
		}
	}
	// 获取map的长度
	sizeLen := len(cardSizeMap)
	colorLen := len(cardColorMap)
	//fmt.Println("colorLen=",colorLen)

	if colorLen > 1 {
		// 非同花
		switch sizeLen {
		case 4:
			// 一对
			cardType = 9
			return
		case 2 : // 3带2  或是 4带1
			// 遍历map value
			for _,v:=range cardSizeMap {
				if v == 4{
					cardType = 3
					return
				}
			}
			cardType = 4
			return
		case 3 :
			// 3条 或是 两对
			for _,v :=range cardSizeMap {
				if v == 3 {
					cardType = 7
					return
				}
			}
			cardType = 8
			return
		case 5:
			// 单牌或是顺子
			isShun ,max:= IsShunZiNew(card)
			if isShun {
				resMax = max
				cardType = 6
				return
			}else{
				cardType = 10
				return
			}

		}

	}else {
		// 同花 或是 同花顺
		//fmt.Println("同花")
		isShun,max := IsShunZiNew(card)
		if isShun {
			resMax = max
			cardType = 1
		}else {
			cardType = 5
		}

	}

	return

}

func IsShunZiNew(card []byte) (shunZi bool, max byte) {
	shunZi = false
	saves := make([]byte, 14)
	for i, v := range card {
		if i %2 ==0 {
			switch v {
			case 50:
				saves[1] = v
			case 51:
				saves[2] = v
			case 52:
				saves[3] = v
			case 53:
				saves[4] = v
			case 54:
				saves[5] = v
			case 55:
				saves[6] = v
			case 56:
				saves[7] = v
			case 57:
				saves[8] = v
			case 84:
				saves[9] = v
			case 74:
				saves[10] = v
			case 81:
				saves[11] = v
			case 75:
				saves[12] = v
			case 65:
				saves[13] = v
				saves[0] = v
			default:
				fmt.Println("无法解析的扑克牌","card --v=",v)
			}
		}

	}
	// 判断数组是否连续
	sum := 0
	for i:=len(saves)-1;i>=0;i-- {
		if saves[i] != 0x00 {
			sum++
		} else {
			sum = 0
		}
		if sum >= 5 {
			// break
			shunZi = true
			max = saves[i+4]
			//fmt.Printf("saves[%#v]=%#v\n",i,saves[i])
			return
		}
	}
	return
}

func QuickSortByte(bs []byte) []byte {
	if len(bs) <= 1 {
		return bs
	}
	splitdata := bs[0]           // 第一个数据
	low := make([]byte, 0, 0)   // 比我小的数据
	hight := make([]byte, 0, 0) // 比我大的数据
	mid := make([]byte, 0, 0)   // 与我一样大的数据
	mid = append(mid, splitdata)   // 加入一个
	for i := 1; i < len(bs); i++ {
		if bs[i] > splitdata {
			low = append(low, bs[i])
		} else if bs[i] < splitdata {
			hight = append(hight, bs[i])
		} else {
			mid = append(mid, bs[i])
		}
	}
	low, hight = QuickSortByte(low), QuickSortByte(hight)
	myarr := append(append(low, mid...), hight...)
	return myarr
}

// 同类型单牌
func (this *CardCom) SingleCardCompareSizeNew()(result int) {

	cardSizeSlice1 := make([]byte,len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte,len(this.cardSizeMap1))

	i:=0
	for k,_:= range this.cardSizeMap1 {
		cardSizeSlice1[i] = SizeTranByte(k)
		i++
	}
	i = 0
	for k,_:= range this.cardSizeMap2 {
		cardSizeSlice2[i] = SizeTranByte(k)
		i++
	}

	result = SingleCardSizeCom(5,cardSizeSlice1,cardSizeSlice2)

	return
}

// 对比单牌 大小
func SingleCardSizeCom(comLen int, cardSizeSlice1,cardSizeSlice2 []byte)(result int) {
	cardSizeSlice1 = QuickSortByte(cardSizeSlice1)
	cardSizeSlice2 = QuickSortByte(cardSizeSlice2)

	// 一个个对比
	for i:=0;i<comLen;i++ {
		if cardSizeSlice1[i] > cardSizeSlice2[i]{
			return 1
		}else if cardSizeSlice1[i] < cardSizeSlice2[i]{
			return 2
		}
	}
	return 0
}

func (this *CardCom) aPairComNew()(result int) {
	cardSizeSlice1 := make([]byte,len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte,len(this.cardSizeMap1))

	var pair1 byte
	var pair2 byte
	i:=0
	for k,v:= range this.cardSizeMap1 {
		if v == 2 {
			pair1 = SizeTranByte(k)

		}else {
			cardSizeSlice1[i] = SizeTranByte(k)
			i++
		}
	}
	i = 0
	for k,v:= range this.cardSizeMap2 {
		if v == 2 {
			pair2 = SizeTranByte(k)

		}else {
			cardSizeSlice2[i] = SizeTranByte(k)
			i++
		}
	}
	//fmt.Println("pair1=",pair1,"  pair2=",pair2)
	if pair1 > pair2 {
		return 1
	} else if pair1 <pair2 {
		return  2
	}else {
		//fmt.Println("cardSizeSlice1=",cardSizeSlice1,"  cardSizeSlice2=",cardSizeSlice2)
		result = SingleCardSizeCom(3,cardSizeSlice1,cardSizeSlice2)
		return
	}
	return
}

func (this *CardCom) twoPairComNew()(result int){
	pairs1 :=make([]byte,2)
	pairs2 :=make([]byte,2)
	var val1 byte
	var val2 byte

	i:=0
	for k,v := range this.cardSizeMap1 {
		if v == 2{
			pairs1[i] = SizeTranByte(k)
			i++
		}else {
			val1 = SizeTranByte(k)
		}

	}
	i = 0
	for k,v := range this.cardSizeMap2 {
		if v == 2{
			pairs2[i] = SizeTranByte(k)
			i++
		}else{
			val2 = SizeTranByte(k)
		}

	}
	result = SingleCardSizeCom(2,pairs1,pairs2)
	if result != 0{
		return
	}else{

		if val1  >val2  {
			return 1
		}else if val1 < val2 {
			return 2
		}else {
			return 0
		}
	}
	return
}

func (this *CardCom)onlyThreeComNew()(result int) {
	cardSizeSlice1 := make([]byte,len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte,len(this.cardSizeMap1))

	var three1 byte
	var three2 byte
	i:=0
	for k,v:= range this.cardSizeMap1 {
		cardSizeSlice1[i] = SizeTranByte(k)
		if v == 3 {
			three1 = SizeTranByte(k)
		}else {
			i++
		}
	}
	i = 0
	for k,v:= range this.cardSizeMap2 {
		cardSizeSlice2[i] = SizeTranByte(k)

		if v == 3 {
			three2 = SizeTranByte(k)
		}else{
			i++
		}
	}
	//fmt.Println("three1=",three1,"  three2=",three2)
	if three1 > three2 {
		return 1
	}else if three1 < three2 {
		return 2
	}else {
		result = SingleCardSizeCom(2,cardSizeSlice1,cardSizeSlice2)
		return
	}
}


func (this *CardCom)onlyShunZiNew()(result int) {
	result = this.SingleCardCompareSizeNew()
	return
}

func (this *CardCom) onlySameFlowerNew()(result int) {
	result = this.SingleCardCompareSizeNew()
	return
}
// 同类型同花顺
func (this *CardCom) straightFlushNew()(result int){
	result = this.SingleCardCompareSizeNew()
	return
}

func (this *CardCom) fourComNew()(result int){
	var four1 byte
	var four2 byte
	var val1 byte
	var val2 byte

	for k,v:=range this.cardSizeMap1 {
		if v == 4 {
			four1 = k
		}else{
			val1 = k
		}
	}
	for k,v:=range this.cardSizeMap2 {
		if v == 4 {
			four2 = k
		}else{
			val2 = k
		}
	}

	if four1 >four2 {
		return 1
	}else if four1 <four2 {
		return 2
	}else {
		if val1 > val2 {
			return 1
		}else if val1<val2 {
			return 2
		}else  {
			return 0
		}
	}
}

func (this *CardCom) threeAndTwoNew()(result int){
	var three1 byte
	var three2 byte
	var two1 byte
	var two2 byte
	for k ,v := range this.cardSizeMap1 {
		if v == 3 {
			three1 = SizeTranByte(k)
		}else {
			two1 = SizeTranByte(k)
		}
	}
	for k ,v := range this.cardSizeMap2 {
		if v == 3 {
			three2 = SizeTranByte(k)
		}else {
			two2 = SizeTranByte(k)
		}
	}
	if three1 > three2 {
		return 1
	}else if three1 < three2 {
		return 2
	}else {
		if two1 > two2 {
			return 1
		}else if two1 <two2 {
			return 2
		}else {
			return  0
		}
	}

}

func PokerMan() {
	file := "/home/weilijie/chromeDown/match_result.json"
	alices := make([]string, 1024)
	bobs := make([]string, 1024)
	results := make([]int, 1024)
	alices, bobs, results = ReadFile(file)
	t1 := time.Now()
	k := 0
	for i := 0; i < len(alices); i++ {
		result := -1
		val1, cardSizesMap1,max1:= JudgMentGroupNew([]byte(alices[i]))
		val2, cardSizesMap2,max2:= JudgMentGroupNew([]byte(bobs[i]))
		fmt.Println("max1=",max1," max2=",max2)
		if val1 < val2 {
			result = 1
		} else if val1 > val2 {
			result = 2
		} else {
			// 牌型处理相同的情况
			// ...
			cardCom := CardCom {
				cardSizeMap1:cardSizesMap1,
				cardSizeMap2:cardSizesMap2,
			}
			switch val1 {
			case 10:
				// 同类型下的单张大牌比较
				result = cardCom.SingleCardCompareSizeNew()
			case 9:
				// 同类型的一对
				result = cardCom.aPairComNew()
			case 8:
				// 同类型两对
				result = cardCom.twoPairComNew()
			case 7:
				// 同类型三条
				result = cardCom.onlyThreeComNew()
			case 6:
				// 同类型顺子
				result = cardCom.onlyShunZiNew()
			case 5:
				// 同类型同花
				result = cardCom.onlySameFlowerNew()
			case 4 :
				// 同类型3带2
				result = cardCom.threeAndTwoNew()
			case 3:
				// 同类型四条
				result = cardCom.fourComNew()
			case 1: // 同类型同花顺
				result = cardCom.straightFlushNew()
			}

			// 最后比较结果
		}

		if result != results[i] {
			k++
			fmt.Printf("[%#v]判断错误--->alice:%#v,bob:%#v<----- ===>文档的结果：%#v, 我的结果:%#v <==\n",k, alices[i], bobs[i],results[i],result)
		} else {
			//fmt.Println("判断正确222222")
		}
	}
	t2 := time.Now()
	fmt.Println("timetime--->",t2.Sub(t1))

}


// func main() {
// 	// file := "/home/weilijie/chromeDown/match_result.json"
// 	// ReadFile(file)
// 	PokerMan()
//
// 	 // res1 ,Map1 := JudgMentGroupNew([]byte("Kc9d8h6dKs"))
// 	 // res2,Map2 := JudgMentGroupNew([]byte("Kd4d5cKhJc"))
// 	 // fmt.Println("res1=",res1,"  Map1=",Map1)
// 	 // fmt.Println("res2=",res2,"  Map2=",Map2)
// 	 // res := aPairComNew(Map1,Map2)
// 	 // fmt.Println("res=",res)
//
// 	 // type1 ,_:=JudgMentGroupNew([]byte("6h6s6cJhAd"))
// 	 // fmt.Println("type1=",type1)
// 	 //3hAh2cAsAc",bob:"6h6s6cJhAd"
// }