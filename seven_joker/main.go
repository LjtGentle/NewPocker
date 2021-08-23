package main

import (
	"fmt"
	"ljtTest/NewPocker/fire"
	"time"
)

type cardType uint8

type comSevenJoker struct {
	cardSizeMap1, cardSizeMap2 map[byte]int
	resCard1, resCard2         byte
	joker1, joker2             int // 记录癞子数量
}

// 找到同花对应的面值序列
func SameFlowerSeq(cardColorMap map[byte]int, card []byte, joker int) (sizeSlice []byte) {
	// 1.找到同花的花色
	var color byte

	// fmt.Println("onlySameFlower-->cardColorMap1=",cardColorMap1,"  cardColorMap2=",cardColorMap2)
	sliceLen := 0
	for k, v := range cardColorMap {
		// fmt.Println("v1=",v)
		if v+joker >= 5 {
			color = k
			sliceLen = v
			break
		}
	}
	// fmt.Println("sliceLen=",sliceLen)
	sizeSlice = make([]byte, sliceLen+joker)
	// fmt.Println("color1=",color1)
	j := 0
	for i := 1; i < len(card); i += 2 {
		if card[i] == color {
			if card[i-1] == 88 {

			} else {
				sizeSlice[j] = fire.SizeTranByte(card[i-1])
				j++
			}

		}
	}
	// fmt.Println("sizeSlice=",sizeSlice)
	// 将所有jock都变A
	for joker > 0 && sliceLen != 0 {
		sizeSlice[j] = 0x0E
		j++
		joker--
	}
	// fmt.Println("joker sizeSlice=",sizeSlice)
	return

}

// 传进来的已经转译好的面值
func IsShunZi(seq []byte, joker int) (shunZi bool, max byte) {

	saves := make([]byte, 15)

	for _, v := range seq {

		switch v {
		case 0x02:
			saves[1] = v
		case 0x03:
			saves[2] = v
		case 0x04:
			saves[3] = v
		case 0x05:
			saves[4] = v
		case 0x06:
			saves[5] = v
		case 0x07:
			saves[6] = v
		case 0x08:
			saves[7] = v
		case 0x09:
			saves[8] = v
		case 0x0A:
			saves[9] = v
		case 0x0B:
			saves[10] = v
		case 0x0C:
			saves[11] = v
		case 0x0D:
			saves[12] = v
		case 0x0E:
			saves[13] = v
			saves[0] = v
		case 0x10:
			saves[14] = v

		default:
			fmt.Println("IsShunZi say 无法解析的扑克牌", "card --v=", v)
		}

	}
	sum := 0
	// 判断数组是否连续
	if joker < 1 {
		// 没有癞子的顺子
		for i := len(saves) - 1; i >= 0; i-- {
			if saves[i] != 0x00 {
				sum++
			} else {
				sum = 0
			}
			if sum >= 5 {
				// break
				shunZi = true
				max = saves[i+4]
				// fmt.Printf("saves[%#v]=%#v\n",i,saves[i])
				return
			}
		}
	} else {
		tmp := joker
		// // 有癞子的顺子  逆序
		// for i := len(saves)-1; i>=4; i-- {
		// 	for j := i; j>=0; j-- {
		// 		if saves[j] != 0x00 {
		// 			sum++
		// 		}else if joker >0{
		// 			joker --
		// 			sum++
		// 		}else {
		// 			joker = tmp
		// 			sum =0
		// 		}
		// 		if sum >= 5 {
		// 			fmt.Println("save=",saves[j:j+4])
		// 			shunZi = true
		// 			max = saves[j+4]
		// 			if max == 0 {
		// 				//fmt.Println("---max=0")
		// 				max = IndexFindByte(j+4)
		// 			}
		// 			return
		// 		}
		// 	}
		// 	sum = 0
		// 	joker = tmp
		// }

		sum = 0
		// fmt.Println("saves=",saves)
		for i := len(saves) - 1; i >= 0; i-- {
			if saves[i] != 0 {
				// fmt.Println(" 0 i=",i)
				sum++
			} else if joker > 0 {
				joker--
				sum++
			} else {
				// 这里是回退到一开始的下一个
				// fmt.Println("before i=",i)
				i = i + sum - joker
				// fmt.Println("after i=",i)
				joker = tmp
				sum = 0
			}
			if sum >= 5 {
				// fmt.Println("顺子啊")
				// fmt.Println("slice=",saves[i:i+5],"i=",i)
				max = saves[i+4]
				if max == 0 {
					// fmt.Println("---max=0,i=",i)
					max = IndexTranByte(i + 4)
				}
				shunZi = true
				return
			}
		}
	}
	return
}

func IndexTranByte(index int) (b byte) {
	switch index {
	case 1:
		b = 0x02
	case 2:
		b = 0x03
	case 3:
		b = 0x04
	case 4:
		b = 0x05
	case 5:
		b = 0x06
	case 6:
		b = 0x07
	case 7:
		b = 0x08
	case 8:
		b = 0x09
	case 9:
		b = 0x0A
	case 10:
		b = 0x0B
	case 11:
		b = 0x0C
	case 12:
		b = 0x0D
	case 13:
		fallthrough
	case 0:
		b = 0x0E

	case 14:
		b = 0x10

	default:
		fmt.Println("IsShunZi say 无法解析的扑克牌", "card --b=", b)
	}
	return
}

// 传进来的card还没有转译
func IsShunZiNoTran(card []byte, joker int) (shunZi bool, max byte) {
	shunZi = false
	saves := make([]byte, 14)
	for i, v := range card {
		if i%2 == 0 {
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
			case 88:
				// continue
				// fmt.Println("88")
			default:
				fmt.Println("无法解析的扑克牌", "card --v=", v)
			}
		}

	}
	// 判断数组是否连续
	sum := 0

	if joker < 1 {
		// 没有癞子的顺子
		for i := len(saves) - 1; i >= 0; i-- {
			if saves[i] != 0x00 {
				sum++
			} else {
				sum = 0
			}
			if sum >= 5 {
				// break
				shunZi = true
				max = saves[i+4]
				// fmt.Printf("saves[%#v]=%#v\n",i,saves[i])
				return
			}
		}
	} else {
		tmp := joker
		// // 有癞子的顺子  逆序
		// for i := len(saves)-1; i>=4; i-- {
		// 	for j := i; j>=0; j-- {
		// 		if saves[j] != 0x00 {
		// 			sum++
		// 		}else if joker >0{
		// 			joker --
		// 			sum++
		// 		}else {
		// 			joker = tmp
		// 			sum =0
		// 		}
		// 		if sum >= 5 {
		// 			fmt.Println("save=",saves[j:j+4])
		// 			shunZi = true
		// 			max = saves[j+4]
		// 			if max == 0 {
		// 				//fmt.Println("---max=0")
		// 				max = IndexFindByte(j+4)
		// 			}
		// 			return
		// 		}
		// 	}
		// 	sum = 0
		// 	joker = tmp
		// }

		sum = 0
		// fmt.Println("saves=",saves)
		for i := len(saves) - 1; i >= 0; i-- {
			if saves[i] != 0 {
				// fmt.Println(" 0 i=",i)
				sum++
			} else if joker > 0 {
				joker--
				sum++
			} else {
				// 这里是回退到一开始的下一个
				// fmt.Println("before i=",i)
				i = i + sum - joker
				// fmt.Println("after i=",i)
				joker = tmp
				sum = 0
			}
			if sum >= 5 {
				// fmt.Println("顺子啊")
				// fmt.Println("slice=",saves[i:i+5],"i=",i)
				max = saves[i+4]
				if max == 0 {
					// fmt.Println("---max=0,i=",i)
					max = IndexFindByte(i + 4)
				}
				shunZi = true
				return
			}
		}
	}

	// fmt.Println("sum=",sum)
	return
}

// 下标转成byte
func IndexFindByte(index int) (b byte) {
	switch index {

	case 5:
		b = 54
	case 6:
		b = 55
	case 7:
		b = 56
	case 8:
		b = 57
	case 9:
		b = 84
	case 10:
		b = 74
	case 11:
		b = 81
	case 12:
		b = 75
	case 13:
		b = 65
	case 0:
		b = 65
	case 88:
	case 1:
		b = 50
	case 2:
		b = 51
	case 3:
		b = 52
	case 4:
		b = 53
		// continue
		// fmt.Println("88")
	default:
		fmt.Println("无法解析的下标", "card --index=", index)

	}
	return
}

// 判断牌型  传进来的参数：面值+花色
func judgmentGroup(card [] byte) (cardType cardType, cardSizeMap, cardColorMap map[byte]int, resCard byte, joker int) {
	//  对传进来的牌分拣成  面值 花色
	cardColorMap = make(map[byte]int, 7)
	cardSizeMap = make(map[byte]int, 7)
	// 扫描牌 分别放好大小，花色
	for i, v := range card {
		if i%2 == 0 {
			// 大小
			// 判断是否有癞子
			if v == 88 {
				joker++ // 记录多少个癞子但是不放入map中
				continue
			}
			if _, ok := cardSizeMap[v]; ok {
				cardSizeMap[v] += 1
			} else {
				cardSizeMap[v] = 1
			}
			// 颜色
		} else {
			if v == 110 { // 癞子的花色
				continue
			}
			if _, ok := cardColorMap[v]; ok {
				cardColorMap[v] += 1
			} else {
				cardColorMap[v] = 1
			}
		}
	}

	sizeLen := len(cardSizeMap)
	// colorLen := len(cardColorMap)
	// fmt.Println("sizeLen=",sizeLen)
	flag := false
	for _, v := range cardColorMap {
		if v+joker >= 5 {
			flag = true
			break
		}
	}

	if flag {
		// 同花
		// 判断是不是顺子
		seq := SameFlowerSeq(cardColorMap, card, joker)
		// fmt.Println("seq=",seq)
		isShun, max := IsShunZi(seq[0:len(seq)-joker], joker)
		if isShun {
			// 同花顺
			// fmt.Println("seq=",seq)
			resCard = max
			cardType = 1
			return
		}
		if joker == 0 {
			cardType = 5
			return
		} else {
			// 有癞子的情况下，是同花，也有可能是变成4条 或是 三带2
			i := 0
			for _, v := range cardSizeMap {
				if v+joker == 4 {
					// 4 条
					// resCard = k
					cardType = 3
					return
				} else if v+joker == 3 {
					i++
				}
			}
			if i == 2 {
				// 3带2
				cardType = 4
				return
			}
			cardType = 5
			return
		}

	} else {
		// 不是同花
		switch sizeLen {
		case 7:
			// 单牌
			if isShun, max := IsShunZiNoTran(card, joker); isShun {
				cardType = 6
				resCard = max
				return
			}
			cardType = 10
			return
		case 6:
			// fmt.Println("66666666")
			// 1对
			// fmt.Println("len---->",6,"  card=",string(card))
			if isShun, max := IsShunZiNoTran(card, joker); isShun {
				resCard = max
				cardType = 6
				return
			}
			// 就算有癞子也是一对， 因为癞子不放入map，map len为6 如果有癞子的话，其他牌是单牌
			cardType = 9
			return

		case 5:
			// fmt.Println("555555555")
			if isShun, max := IsShunZiNoTran(card, joker); isShun {
				resCard = max
				cardType = 6
				return
			}
			// 有癞子的五张牌一定是3条
			if joker > 0 {
				cardType = 7
				return
			}
			for _, v := range cardSizeMap {
				if v == 3 {
					flag = true
					break
				}
			}
			if flag {
				// 3条
				cardType = 7
				return
			} else {
				// 2对
				cardType = 8
				return
			}
		case 4:
			i := 0
			j := 0
			for _, v := range cardSizeMap {
				if v+joker == 4 {
					// 4条
					cardType = 3
					return
				} else if v+joker == 3 {
					// 3条
					i++

				} else if v == 2 {
					j++
				}
			}
			if j < 2 {
				cardType = 4
				return
			}
			cardType = 8
			return
		case 3:
			for _, v := range cardSizeMap {
				if v+joker == 4 {
					cardType = 3
					return
				}
			}
			cardType = 4
			return
		case 2:
			cardType = 3
			return

		}
	}
	return
}

// 同类型同花顺比较  --1 1
func (this *comSevenJoker) straightFlush() (result int) {
	if this.resCard1 > this.resCard2 {
		return 1
	} else if this.resCard1 < this.resCard2 {
		return 2
	} else {
		return 0
	}
}

// 同类型四条比较  --3
func (this *comSevenJoker) fourCom(joker1, joker2 int) (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap2))
	var four1 byte
	var four2 byte

	i := 0
	for k, v := range this.cardSizeMap1 {
		k = fire.SizeTranByte(k)
		if v+joker1 == 4 {
			if four1 < k {
				four1 = k
			}
		} else {
			cardSizeSlice1[i] = k
			i++
		}
	}
	i = 0
	for k, v := range this.cardSizeMap2 {
		k = fire.SizeTranByte(k)
		if v+joker2 == 4 {
			if four2 < k {
				four2 = k
			}
		} else {
			cardSizeSlice2[i] = k
			i++
		}
	}
	// fmt.Println("four1=",four1,"  four2=",four2)
	// fmt.Println("cardSizeSlice1=",cardSizeSlice1,"  cardSizeSlice2=",cardSizeSlice2)
	if four1 > four2 {
		return 1
	} else if four1 < four2 {
		return 2
	} else {
		result = fire.SingleCardSizeCom(1, cardSizeSlice1, cardSizeSlice2)
		return
	}
}

// 同类型3带2比较  --4
func (this *comSevenJoker) threeAndTwo() (result int) {
	threes1 := make([]byte, 3)
	threes2 := make([]byte, 3)
	twos1 := make([]byte, 3)
	twos2 := make([]byte, 3)
	// fmt.Println("three1=",three1)
	i := 0
	j := 0

	for k, v := range this.cardSizeMap1 {
		// fmt.Println("1---->k=",k,"  v=",v)
		k = fire.SizeTranByte(k)
		if v+this.joker1 == 3 {
			threes1[i] = k
			i++
		} else if v == 2 {
			// fmt.Println("twoooooooooo")
			twos1[j] = k
			// fmt.Println("2222222222")
			// fmt.Printf("twos1[%#v]=%#v\n",j,twos1[j])
			j++
		}
	}
	i = 0
	j = 0
	for k, v := range this.cardSizeMap2 {
		// fmt.Println("2---->k=",k,"  v=",v)
		k = fire.SizeTranByte(k)
		if v+this.joker2 == 3 {
			// fmt.Println("333333333333")
			threes2[i] = k
			i++

		} else if v == 2 {
			// fmt.Println("222222222")
			twos2[j] = k
			// fmt.Println("twos2[j]=",twos2[j])
			// fmt.Println("2222-----2222222")
			j++
		}
	}

	// fmt.Println("threes1 = ",threes1)
	// fmt.Println("threes2 = ",threes2)
	// fmt.Println("twos1 = ",twos1)
	// fmt.Println("twos2 = ",twos2)

	// 对3条排序
	threes1 = fire.QuickSortByte(threes1)
	threes2 = fire.QuickSortByte(threes2)

	if threes1[0] > threes2[0] {
		return 1
	} else if threes1[0] < threes2[0] {
		return 2
	}

	// 对 对子排序
	twos1 = fire.QuickSortByte(twos1)
	twos2 = fire.QuickSortByte(twos2)

	if twos1[0] < threes1[1] {
		twos1[0] = threes1[1]
	}
	if twos2[0] < threes2[1] {
		twos2[0] = threes2[1]
	}
	if twos1[0] > twos2[0] {
		return 1
	} else if twos1[0] < twos2[0] {
		return 2
	} else {
		return 0
	}

}

// 同类型同花比较  --5  1
func (this *comSevenJoker) onlyFlush(cardColorMap1, cardColorMap2 map[byte]int, card1, card2 []byte) (result int) {
	sizeSlice1 := SameFlowerSeq(cardColorMap1, card1, this.joker1)
	sizeSlice2 := SameFlowerSeq(cardColorMap2, card2, this.joker2)
	result = fire.SingleCardSizeCom(5, sizeSlice1, sizeSlice2)
	return
}

// 同类型顺子比较 --6
func (this *comSevenJoker) OnlyShunZi() (result int) {
	v1 := fire.SizeTranByte(this.resCard1)
	v2 := fire.SizeTranByte(this.resCard2)
	if v1 > v2 {
		return 1
	} else if v1 < v2 {
		return 2
	} else {
		return 0
	}
}

// 同类型3条比较  --7
func (this *comSevenJoker) onlyThree(joker1, joker2 int) (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap2))
	var three1 byte
	var three2 byte

	i := 0

	for k, v := range this.cardSizeMap1 {
		if v+joker1 == 3 {
			if three1 < fire.SizeTranByte(k) {
				cardSizeSlice1[i] = three1
				i++
				three1 = fire.SizeTranByte(k)
			}
		} else {
			cardSizeSlice1[i] = fire.SizeTranByte(k)
			i++
		}
	}
	i = 0
	for k, v := range this.cardSizeMap2 {
		if v+joker2 == 3 {
			if three2 < fire.SizeTranByte(k) {
				cardSizeSlice2[i] = three2
				i++
				three2 = fire.SizeTranByte(k)
			}
		} else {
			cardSizeSlice2[i] = fire.SizeTranByte(k)
			i++
		}
	}
	if three1 > three2 {
		return 1
	} else if three1 < three2 {
		return 2
	} else {
		result = fire.SingleCardSizeCom(2, cardSizeSlice1, cardSizeSlice2)
		return
	}
}

// 同类型两对比较 --8
func (this *comSevenJoker) TwoPair(joker1, joker2 int) (result int) {

	pairs1 := make([]byte, 3)
	pairs2 := make([]byte, 3)
	vals1 := make([]byte, 3)
	vals2 := make([]byte, 3)

	j := 0
	i := 0
	for k, v := range this.cardSizeMap1 {
		k = fire.SizeTranByte(k)
		if v+joker1 == 2 {
			pairs1[i] = k
			i++
		} else {
			vals1[j] = k
			j++
		}
	}
	i = 0
	j = 0
	for k, v := range this.cardSizeMap2 {
		k = fire.SizeTranByte(k)
		if v+joker2 == 2 {
			pairs2[i] = k
			i++
		} else {
			vals2[j] = k
			j++
		}
	}
	// fmt.Println("pairs1=",pairs1,"  pairs2=",pairs2)
	// 对 对子 序列排序
	pairs1 = fire.QuickSortByte(pairs1)
	pairs2 = fire.QuickSortByte(pairs2)

	for i := 0; i < 2; i++ {
		if pairs1[i] > pairs2[i] {
			return 1
		} else if pairs1[i] < pairs2[i] {
			return 2
		}
	}

	// 对单牌序列排序
	vals1 = fire.QuickSortByte(vals1)
	vals2 = fire.QuickSortByte(vals2)
	if vals1[0] < pairs1[2] {
		vals1[0] = pairs1[2]
	}
	if vals2[0] < pairs2[2] {
		vals2[0] = pairs2[2]
	}
	// fmt.Println("vals1=",vals1,"  vals2=",vals2)

	if vals1[0] > vals2[0] {
		return 1
	} else if vals1[0] < vals2[0] {
		return 2
	}
	return 0
}

// 同类型一对比较 --9
func (this *comSevenJoker) OnePair() (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap1))
	var val1 byte
	var val2 byte
	i := 0

	for k, v := range this.cardSizeMap1 {
		k = fire.SizeTranByte(k)
		if v == 2 {
			val1 = k
			continue
		}
		cardSizeSlice1[i] = k
		i++
	}
	i = 0
	for k, v := range this.cardSizeMap2 {
		k = fire.SizeTranByte(k)
		if v == 2 {
			val2 = k
			continue
		}

		cardSizeSlice2[i] = k
		i++
	}

	cardSizeSlice1 = fire.QuickSortByte(cardSizeSlice1)
	cardSizeSlice2 = fire.QuickSortByte(cardSizeSlice2)

	if val1 == 0 {
		val1 = cardSizeSlice1[0]
	}
	if val2 == 0 {
		val2 = cardSizeSlice2[0]
		// fmt.Println("val22222222")
	}

	if val1 > val2 {
		return 1
	} else if val1 < val2 {
		return 2
	}

	comLen := 3
	// 一个个对比
	for i := 0; i < comLen; i++ {
		if cardSizeSlice1[i+this.joker1] > cardSizeSlice2[i+this.joker2] {
			return 1
		} else if cardSizeSlice1[i+this.joker1] < cardSizeSlice2[i+this.joker2] {
			return 2
		}
	}

	return 0
}

// 同类型单牌比较  --10
func (this *comSevenJoker) SingleCard() (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap2))
	i := 0
	for k, _ := range this.cardSizeMap1 {
		cardSizeSlice1[i] = fire.SizeTranByte(k)
		i++
	}
	i = 0
	for k, _ := range this.cardSizeMap2 {
		cardSizeSlice2[i] = fire.SizeTranByte(k)
		i++
	}
	result = fire.SingleCardSizeCom(5, cardSizeSlice1, cardSizeSlice2)
	return
}

func pockerMan() {
	filePath := "/home/weilijie/chromeDown/seven_cards_with_ghost.result.json"
	alices, bobs, results := fire.ReadFile(filePath)
	t1 := time.Now()
	k := 0
	for i := 0; i < len(alices); i++ {

		result := -1
		// 先判断各自的牌型
		cardType1, cardSizesMap1, cardColorMap1, max1, joker1 := judgmentGroup([]byte(alices[i]))
		cardType2, cardSizesMap2, cardColorMap2, max2, joker2 := judgmentGroup([]byte(bobs[i]))
		if cardType1 < cardType2 {
			result = 1
		} else if cardType1 > cardType2 {
			result = 2
		} else {

			csj := &comSevenJoker{
				cardSizeMap1: cardSizesMap1,
				cardSizeMap2: cardSizesMap2,
				resCard1:     max1,
				resCard2:     max2,
				joker1:       joker1,
				joker2:       joker2,
			}
			// 同类型比较
			switch cardType1 {
			case 1:
				// 同花顺
				result = csj.straightFlush()
			case 3:
				// 4条
				result = csj.fourCom(joker1, joker2)
			case 4:
				// 3带2
				result = csj.threeAndTwo()
			case 5:
				// 同花
				result = csj.onlyFlush(cardColorMap1, cardColorMap2, []byte(alices[i]), []byte(bobs[i]))
			case 6:
				// 顺子
				result = csj.OnlyShunZi()
			case 7:
				// 3条
				result = csj.onlyThree(joker1, joker2)
			case 8:
				// 两对
				result = csj.TwoPair(joker1, joker2)
			case 9:
				// 一对
				result = csj.OnePair()
			case 10:
				// 单牌
				result = csj.SingleCard()

			}
		}
		if result != results[i] {
			k++
			fmt.Println("[", k, "]"+"判断有误，alice=", alices[i], " bob=", bobs[i], "  我的结果:", result, "  文档的结果：", results[i])
		}
	}

	t2 := time.Now()
	fmt.Println("time--->", t2.Sub(t1))

}

func main() {

	pockerMan()

	// byte := "Xn"
	//
	// for i,v := range byte {
	// 	fmt.Printf("byte[%#v] = %#v\n",i,v)
	// }
	// //
	// cardType1, cardSizesMap1, cardColorMap1, max1,joker1 := judgmentGroup([]byte("Xn8d8s3h3c9hQd"))
	// cardType2, cardSizesMap2, cardColorMap2, max2,joker2 := judgmentGroup([]byte("KdKsXn8d8s3h3c"))
	// fmt.Println("cardType1=",cardType1,"  cardSizesMap1=",cardSizesMap1,"  cardColorMap1=",cardColorMap1, "  max1=",max1,"  joker1=",joker1)
	// fmt.Println("cardType2=",cardType2,"  cardSizesMap2=",cardSizesMap2,"  cardColorMap2=",cardColorMap2, "  max2=",max2,"  joker2=",joker2)
	// // // //
	// jocker :=&comSevenJoker{
	// 	joker1:joker1,
	// 	joker2:joker2,
	// 	cardSizeMap1:cardSizesMap1,
	// 	cardSizeMap2:cardSizesMap2,
	// }
	// // // // result := jocker.OnePair()
	// // // result := jocker.fourCom(0,0)
	// //
	// //  result := jocker.TwoPair(0,0)
	// result := jocker.threeAndTwo()
	//  fmt.Println("result=",result)
	//
	//
	//  isShun ,max:=  IsShunZiNoTran([]byte("Xn7h8s3sTs6d2c"),1)
	//  fmt.Println("isShun=",isShun,"  max=",max)
	//
	// isShun2 ,max2:=  IsShunZiNoTran([]byte("9hJhXn8s7c3d6c"),1)
	// fmt.Println("isShun=",isShun2,"  max=",max2)
	//
	// isShun3 ,max3:=  IsShunZiNoTran([]byte("Qs4hXn5h6c7dKd"),1)
	// fmt.Println("isShun=",isShun3,"  max=",max3)
	//
	// isShun4 ,max4:=  IsShunZiNoTran([]byte("3c9sQs4hXn5h6c"),1)
	// fmt.Println("isShun=",isShun4,"  max=",max4)
	//
	// isShun5 ,max5:=  IsShunZiNoTran([]byte("TdJdXn3hQs9sTc"),1)
	// fmt.Println("isShun=",isShun5,"  max=",max5)
	//
	// isShun6 ,max6:=  IsShunZiNoTran([]byte("Th8cTdJdXn3hQs"),1)
	// fmt.Println("isShun=",isShun6,"  max=",max6)

	// b, max := IsShunZi([]byte{3, 4, 12, 2, 14},1)
	// fmt.Println("b=",b,"  max=",max)

}
