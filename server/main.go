package main

import (
	"fmt"
	"ljtTest/NewPocker/fire"
	"time"
)

type CardType uint8

type Seven struct {
	cardSizeMap1, cardSizeMap2 map[byte]int
}

// 传进来的已经转译好的面值
func IsShunZi(seq []byte) (flag bool, max byte) {

	flag = false
	saves := make([]byte, 14)

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
		default:
			fmt.Println("无法解析的扑克牌","card --v=",v)
		}

	}
	// 判断数组是否连续
	sum := 0
	for i := len(saves) - 1; i >= 0; i-- {
		if saves[i] != 0x00 {
			sum++
		} else {
			sum = 0
		}
		if sum >= 5 {
			// break
			flag = true
			max = saves[i+4]
			// fmt.Printf("saves[%#v]=%#v\n",i,saves[i])
			return
		}
	}
	return
}

// 判断牌型
func JudgMentGroup(card []byte) (cardType CardType, cardSizeMap, cardColorMap map[byte]int, resMax byte) {
	cardColorMap = make(map[byte]int, 7)
	cardSizeMap = make(map[byte]int, 7)
	// 扫描牌 分别放好大小，花色
	for i, v := range card {
		if i%2 == 0 {
			// 大小
			if _, ok := cardSizeMap[v]; ok {
				cardSizeMap[v] += 1
			} else {
				cardSizeMap[v] = 1
			}
			// 颜色
		} else {
			if _, ok := cardColorMap[v]; ok {
				cardColorMap[v] += 1
			} else {
				cardColorMap[v] = 1
			}
		}
	}
	// 获取map的长度
	sizeLen := len(cardSizeMap)
	// colorLen := len(cardColorMap)
	// fmt.Println("sizeLen=",sizeLen)
	flag := false
	for _, v := range cardColorMap {
		if v >= 5 {
			flag = true
			break
		}
	}
	if flag {
		// 同花
		// 判断是不是顺子
		seq := SameFlowerSeq(cardColorMap, card)
		//fmt.Println("seq=",seq)
		isShun, max := IsShunZi(seq)
		if isShun {
			// 同花顺
			resMax = max
			cardType = 1
			return
		}
		cardType = 5
		return
	} else {
		switch sizeLen {
		case 7:
			// 单牌
			if isShun, max := fire.IsShunZiNew(card); isShun {
				cardType = 6
				resMax = max
				return
			}
			cardType = 10
			return
		case 6:
			// 1对
			// fmt.Println("len---->",6,"  card=",string(card))
			if isShun, max := fire.IsShunZiNew(card); isShun {
				resMax = max
				cardType = 6
				return
			}
			cardType = 9
			return
		case 5:
			if isShun, max := fire.IsShunZiNew(card); isShun {
				resMax = max
				cardType = 6
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
				if v == 4 {
					// 4条
					cardType = 3
					return
				} else if v == 3 {
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
				if v == 4 {
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

// 单牌
func (this *Seven) SingleCard() (result int) {
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

// 一对
func (this *Seven) APair() (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap1))
	var val1 byte
	var val2 byte
	i := 0
	for k, v := range this.cardSizeMap1 {
		if v == 2 {
			val1 = fire.SizeTranByte(k)
			continue
		}
		cardSizeSlice1[i] = fire.SizeTranByte(k)
		i++
	}
	i = 0
	for k, v := range this.cardSizeMap2 {
		if v == 2 {
			val2 = fire.SizeTranByte(k)
			continue
		}

		cardSizeSlice2[i] = fire.SizeTranByte(k)
		i++
	}
	if val1 > val2 {
		return 1
	} else if val1 < val2 {
		return 2
	}
	result = fire.SingleCardSizeCom(3, cardSizeSlice1, cardSizeSlice2)
	return
}

// 两对
func (this *Seven) TwoPair() (result int) {

	pairs1 := make([]byte, 3)
	pairs2 := make([]byte, 3)
	vals1 := make([]byte, 3)
	vals2 := make([]byte, 3)

	j := 0
	i := 0
	for k, v := range this.cardSizeMap1 {
		k =  fire.SizeTranByte(k)
		if v == 2 {
			pairs1[i] =k
			i++
		} else {
			vals1[j] = k
			j++
		}
	}
	i = 0
	j = 0
	for k, v := range this.cardSizeMap2 {
		k =  fire.SizeTranByte(k)
		if v == 2 {
			pairs2[i] =k
			i++
		} else {
			vals2[j] = k
			j++
		}
	}
	// 对对子序列排序  逆序
	pairs1 = fire.QuickSortByte(pairs1)
	pairs2 = fire.QuickSortByte(pairs2)

	for i :=0; i<2; i++ {
		if pairs1[i] > pairs2[i]{
			return 1
		}else if pairs1[i] < pairs2[i]{
			return 2
		}
	}
	// 对剩余的单牌排序
	vals1 = fire.QuickSortByte(vals1)
	vals2 = fire.QuickSortByte(vals2)
	if vals1[0] <pairs1[2] {
		vals1[0] =  pairs1[2]
	}
	if vals2[0] <pairs2[2] {
		vals2[0] =  pairs2[2]
	}
	if vals1[0] >vals2[0] {
		return 1
	}else if vals1[0] <vals2[0] {
		return 2
	}else {
		return 0
	}
}

// 3条
func (this *Seven) OnlyThree() (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap2))
	var three1 byte
	var three2 byte

	i := 0

	for k, v := range this.cardSizeMap1 {
		if v == 3 {
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
		if v == 3 {
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

// 顺子
func (this *Seven) OnlyShunZi() (result int) {
	return this.SingleCard()
}

// 找到同花对应的面值序列
func SameFlowerSeq(cardColorMap map[byte]int, card []byte) (sizeSlice []byte) {
	// 1.找到同花的花色
	var color byte


	// fmt.Println("onlySameFlower-->cardColorMap1=",cardColorMap1,"  cardColorMap2=",cardColorMap2)
	sliceLen := 0
	for k, v := range cardColorMap {
		// fmt.Println("v1=",v)
		if v >= 5 {
			color = k
			sliceLen = v
			break
		}
	}
	sizeSlice = make([]byte, sliceLen)
	// fmt.Println("color1=",color1)
	j := 0
	for i := 1; i < len(card); i+=2 {
		if card[i] == color {
			sizeSlice[j] = fire.SizeTranByte(card[i-1])
			j++
		}
	}
	//fmt.Println("--->sizeSlice=",sizeSlice)
	return

}

// 同花
func (this *Seven) onlySameFlower(cardColorMap1, cardColorMap2 map[byte]int, card1, card2 []byte) (result int) {
	sizeSlice1 := SameFlowerSeq(cardColorMap1, card1)
	sizeSlice2 := SameFlowerSeq(cardColorMap2, card2)
	result = fire.SingleCardSizeCom(5, sizeSlice1, sizeSlice2)
	return
}

// 3带2
func (this *Seven) ThreeAndTwo() (result int) {

	threes1  := make([]byte,2)
	threes2 := make([]byte,2)
	twos1 := make([]byte,2)
	twos2  := make([]byte,2)
	// fmt.Println("three1=",three1)
	i := 0
	j :=0

	for k, v := range this.cardSizeMap1 {
		// fmt.Println("1---->k=",k,"  v=",v)
		if v == 3 {
			threes1[i] = fire.SizeTranByte(k)
			i++
		} else if v == 2 {
			// fmt.Println("twoooooooooo")
			twos1[j] = fire.SizeTranByte(k)
			//fmt.Println("2222222222")
			//fmt.Printf("twos1[%#v]=%#v\n",j,twos1[j])
			j++
		}
	}
	i = 0
	j = 0
	for k, v := range this.cardSizeMap2 {
		// fmt.Println("2---->k=",k,"  v=",v)
		if v == 3 {
			// fmt.Println("333333333333")
			threes2[i] = fire.SizeTranByte(k)
			i++

		} else if v == 2 {
			// fmt.Println("222222222")
			twos2[j] = fire.SizeTranByte(k)
			//fmt.Println("twos2[j]=",twos2[j])
			//fmt.Println("2222-----2222222")
			j++
		}
	}

	//fmt.Println("425---->twos1[0]=",twos1[0])
	if threes1[1] >0 {
		if threes1[0] > threes1[1]{
			twos1[0] = threes1[1]
		}else {
			twos1[0] = threes1[0]
			threes1[0]= threes1[1]

		}
	}
	if threes2[1] >0 {
		if threes2[0] > threes2[1]{
			twos2[0] = threes2[1]
		}else {
			twos2[0] = threes2[0]
			threes2[0]= threes2[1]

		}
	}
	if twos1[1] >0 {
		if twos1[0]<twos1[1]{
			twos1[0] = twos1[1]
		}
	}
	if twos2[1] >0 {
		if twos2[0]<twos2[1]{
			twos2[0] = twos2[1]
		}
	}
	//fmt.Println("449---->twos1[0]=",twos1[0])
	if len(twos2) >1 {
		if twos2[0]<twos2[1]{
			twos2[0] = twos2[1]
		}
	}

	//fmt.Println("  three1=",threes1,"  two1=",twos1,"  three2=",threes2,"  two2=",twos2)

	if threes1[0] >threes2[0] {
		return 1
	}else if  threes1[0] <threes2[0] {
		return 2
	}else {
		if twos1[0] >twos2[0] {
			return 1
		}else if twos1[0] <twos2[0] {
			return 2
		}else {
			return 0
		}
	}

	return
}

// 4条
func (this *Seven) FourCom() (result int) {
	cardSizeSlice1 := make([]byte, len(this.cardSizeMap1))
	cardSizeSlice2 := make([]byte, len(this.cardSizeMap2))
	var four1 byte
	var four2 byte

	i := 0
	for k, v := range this.cardSizeMap1 {
		if v == 4 {
			four1 = fire.SizeTranByte(k)
		} else {
			cardSizeSlice1[i] = fire.SizeTranByte(k)
			i++
		}
	}
	i = 0
	for k, v := range this.cardSizeMap2 {
		if v == 4 {
			four2 = fire.SizeTranByte(k)
		} else {
			cardSizeSlice2[i] = fire.SizeTranByte(k)
			i++
		}
	}
	if four1 > four2 {
		return 1
	} else if four1 < four2 {
		return 2
	} else {
		result = fire.SingleCardSizeCom(5, cardSizeSlice1, cardSizeSlice2)
		return
	}
	return
}

// 同花顺
func (this *Seven) straightFlush() (result int) {
	// 同花判断是否有癞子

	return this.SingleCard()
}

func PokerMan() {
	//fmt.Println("come in PokerMan()")
	//defer fmt.Println("out of PokerMan")
	file := "/home/weilijie/chromeDown/seven_cards_with_ghost.json"
	alices := make([]string, 1024)
	bobs := make([]string, 1024)
	results := make([]int, 1024)
	alices, bobs, results = fire.ReadFile(file)
	t1 := time.Now()

	j := 0
	k := 0
	for i := 0; i < len(alices); i++ {
		// fmt.Println("i=",i)
		result := -1
		val1, cardSizesMap1, cardColorMap1, max1 := JudgMentGroup([]byte(alices[i]))
		val2, cardSizesMap2, cardColorMap2, max2 := JudgMentGroup([]byte(bobs[i]))
		seven := &Seven{
			cardSizeMap1: cardSizesMap1,
			cardSizeMap2: cardSizesMap2,
		}

		if val1 < val2 {
			result = 1
		} else if val1 > val2 {
			result = 2
		} else {
			switch val1 {
			case 1:
				// 同花顺
				if fire.SizeTranByte(max1) > fire.SizeTranByte(max2) {
					result = 1
				} else if fire.SizeTranByte(max1) < fire.SizeTranByte(max2) {
					result = 2
				} else {
					result = 0
				}
			case 3:
				// 四条
				result = seven.FourCom()
			case 4:
				// 3带2
				result = seven.ThreeAndTwo()
			case 5:
				// 同花
				// fmt.Println("cardSizesMap1=",cardSizesMap1,"  cardSizesMap2=",cardSizesMap2)
				result = seven.onlySameFlower(cardColorMap1, cardColorMap2, []byte(alices[i]), []byte(bobs[i]))
			case 6:
				// 顺子
				if fire.SizeTranByte(max1) > fire.SizeTranByte(max2) {
					result = 1
				} else if fire.SizeTranByte(max1) < fire.SizeTranByte(max2) {
					result = 2
				} else {
					result = 0
				}
			case 7:
				// 3条
				result = seven.OnlyThree()
			case 8:
				// 2对
				result = seven.TwoPair()
			case 9:
				// 一对
				result = seven.APair()
			case 10:
				// 单牌
				result = seven.SingleCard()
			}
		}

		if result == results[i] {
			j++
			// fmt.Println("判断正确", j)
		} else if results[i] != result {
			fmt.Printf("[%#v]判断错误--->alice:%#v,bob:%#v<----- ===>文档的结果：%#v, 我的结果:%#v <==\n",
				k, alices[i], bobs[i], results[i], result)
			k++
		}
	}

	t2 := time.Now()
	fmt.Println("time----->>>", t2.Sub(t1))
}

func main() {
	 PokerMan()
	//
	// // res1 ,Map1:= JudgMentGroup([]byte("9dTdKc2h6h7sQh"))
	// // res2 ,Map2:= JudgMentGroup([]byte("Kc2h6h7sQh3c4c"))
	// // fmt.Println("res1=",res1," res2=",res2," Map1=",Map1," Map2=",Map2)
	//
	// val1, cardSizesMap1, cardColorsMap1, max1 := JudgMentGroup([]byte("As9cAdTsTh3h6s"))
	// val2, cardSizesMap2, cardColorsMap2, max2 := JudgMentGroup([]byte("9s4sAs9cAdTsTh"))
	// fmt.Println("val1=", val1, " Map:", cardSizesMap1, " max1=", max1, "  cardColorsMap1=", cardColorsMap1)
	// fmt.Println("val2=", val2, " Map:", cardSizesMap2, " max2=", max2, "  cardColorsMap2=", cardColorsMap2)
	// seven := &Seven{
	// 	cardSizeMap1: cardSizesMap1,
	// 	cardSizeMap2: cardSizesMap2,
	// }
	// result := seven.TwoPair()
	// //  result := seven.ThreeAndTwo()
	// // // result := seven.onlySameFlower(cardColorsMap1, cardColorsMap1, []byte("6s5s2sAd8s9cAs"), []byte("7cJs6s5s2sAd8s"))
	//  fmt.Println("result=", result)
	//
	// // b, max := fire.IsShunZiNew([]byte("5h4h3d2sAdKh7d"))
	// // fmt.Println("b=",b," max=",max)


}
