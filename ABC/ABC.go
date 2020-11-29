package main

import (
	"math"
	"math/rand"
	"time"
)

var (
	NP         = 40     //种群规模，采蜜蜂+观察蜂
	FoodNumber = NP / 2 // 食物的数量，为采蜜蜂的数量
	limit      = 20     // 限度，超过这个限度没有更新，采蜜蜂变为侦查蜂
	maxCycle   = 10000  // 循环次数

	// 函数参数
	D          = 2    // 函数参数的个数
	lb float64 = -100 // 函数的下界
	ub float64 = 100  // 函数的上界
)

type BeeGroup struct {
	code     [2]float64
	trueFit  float64
	fitness  float64
	rfitness float64
	trail    int
}

var NectarSource [20]BeeGroup
var EmployedBee [20]BeeGroup
var OnLooker [20]BeeGroup
var BestSource BeeGroup

func main() {
	initilize()
	MemorizeBestSource() // 保存最好蜜源
	for gen := 0; gen < maxCycle; gen++ {
		sendEmployedBees()
		CalculateProbabilities()
		sendOnlookerBees()
		MemorizeBestSource()
		sendScoutBees()
		MemorizeBestSource()
		println(BestSource.trueFit)
	}
}

func randFloats(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	value := min + rand.Float64()*(max-min)
	return value
}

func initilize() {
	for i := 0; i < FoodNumber; i++ {
		for j := 0; j < D; j++ {
			NectarSource[i].code[j] = randFloats(lb, ub)
			EmployedBee[i].code[j] = NectarSource[i].code[j]
			OnLooker[i].code[j] = NectarSource[i].code[j]
			BestSource.code[j] = NectarSource[0].code[j]
		}
		// 蜜源初始化
		NectarSource[i].trueFit = calculationTruefit(NectarSource[i])
		NectarSource[i].fitness = calculationFitness(NectarSource[i].trueFit)
		NectarSource[i].rfitness = 0
		NectarSource[i].trail = 0
		// 采蜜蜂初始化
		EmployedBee[i].trueFit = NectarSource[i].trueFit
		EmployedBee[i].fitness = NectarSource[i].fitness
		EmployedBee[i].rfitness = NectarSource[i].rfitness
		EmployedBee[i].trail = NectarSource[i].trail
		// 观察蜂初始化
		OnLooker[i].trueFit = NectarSource[i].trueFit
		OnLooker[i].fitness = NectarSource[i].fitness
		OnLooker[i].rfitness = NectarSource[i].rfitness
		OnLooker[i].trail = NectarSource[i].trail
	}
	// 最优蜜源初始化
	BestSource.trueFit = NectarSource[0].trueFit
	BestSource.fitness = NectarSource[0].fitness
	BestSource.rfitness = NectarSource[0].rfitness
	BestSource.trail = NectarSource[0].trail
}

func calculationTruefit(bee BeeGroup) float64 {
	var truefit float64 = 0
	/******测试函数1******/
	truefit = 0.5 + (math.Sin(math.Sqrt(bee.code[0]*bee.code[0]+bee.code[1]*bee.code[1]))*math.Sin(math.Sqrt(bee.code[0]*bee.code[0]+bee.code[1]*bee.code[1]))-0.5)/((1+0.001*(bee.code[0]*bee.code[0]+bee.code[1]*bee.code[1]))*(1+0.001*(bee.code[0]*bee.code[0]+bee.code[1]*bee.code[1])))

	/******测试函数2******/
	return truefit
}

func calculationFitness(trueFit float64) float64 {
	var fitnessResult float64 = 0
	if trueFit >= 0.0 {
		fitnessResult = 1 / (trueFit + 1)
	} else {
		fitnessResult = 1 + math.Abs(trueFit)
	}
	return fitnessResult
}

func MemorizeBestSource() {
	for i := 0; i < FoodNumber; i++ {
		if NectarSource[i].trueFit < BestSource.trueFit {
			for j := 0; j < D; j++ {
				BestSource.code[j] = NectarSource[i].code[j]
			}
			BestSource.trueFit = NectarSource[i].trueFit
		}
	}
}

func CalculateProbabilities() { //计算轮盘赌的选择概率
	maxfit := NectarSource[0].fitness
	for i := 1; i < FoodNumber; i++ {
		if NectarSource[i].fitness > maxfit {
			maxfit = NectarSource[i].fitness
		}
	}
	for i := 0; i < FoodNumber; i++ {
		NectarSource[i].rfitness = (0.9 * (NectarSource[i].fitness / maxfit)) + 0.1
	}
}

func round(x float64) int {
	//return int(math.Floor(x + 0.5))
	return int(math.Floor(x))
}

func sendEmployedBees() {
	var k int
	var param2change int //需要改变的维数
	var Rij float64      //[-1,1]之间的随机数
	for i := 0; i < FoodNumber; i++ {
		param2change = round(randFloats(0, float64(D))) //随机选取需要改变的维数
		/******选取不等于i的k********/
		for {
			k := round(randFloats(0, float64(FoodNumber)))
			if k != i {
				break
			}
		}
		for j := 0; j < D; j++ {
			EmployedBee[i].code[j] = NectarSource[i].code[j]
		}
		/*******采蜜蜂去更新信息*******/
		Rij = randFloats(-1, 1)
		EmployedBee[i].code[param2change] = NectarSource[i].code[param2change] + Rij*(NectarSource[i].code[param2change]-NectarSource[k].code[param2change])
		/*******判断是否越界********/
		if EmployedBee[i].code[param2change] > ub {
			EmployedBee[i].code[param2change] = ub
		}
		if EmployedBee[i].code[param2change] < lb {
			EmployedBee[i].code[param2change] = lb
		}
		EmployedBee[i].trueFit = calculationTruefit(EmployedBee[i])
		EmployedBee[i].fitness = calculationFitness(EmployedBee[i].trueFit)

		/******贪婪选择策略*******/
		if EmployedBee[i].trueFit < NectarSource[i].trueFit {
			for j := 0; j < D; j++ {
				NectarSource[i].code[j] = EmployedBee[i].code[j]
			}
			NectarSource[i].trail = 0
			NectarSource[i].trueFit = EmployedBee[i].trueFit
			NectarSource[i].fitness = EmployedBee[i].fitness
		} else {
			NectarSource[i].trail++
		}
	}
}

func sendOnlookerBees() { // 采蜜蜂与观察蜂交流信息，观察蜂更改信息
	var k int
	var R_choosed float64 // 被选中的概率
	var param2change int  //需要改变的维数
	var Rij float64       //[-1,1]之间的随机数
	i := 0
	for t := 0; t < FoodNumber; {
		R_choosed = randFloats(0, 1)
		if R_choosed < NectarSource[i].rfitness {
			t++
			param2change = round(randFloats(0, float64(D))) //随机选取需要改变的维数
			/******选取不等于i的k********/
			for {
				k := round(randFloats(0, float64(FoodNumber)))
				if k != i {
					break
				}
			}
			for j := 0; j < D; j++ {
				OnLooker[i].code[j] = NectarSource[i].code[j]
			}
			/*******采蜜蜂去更新信息*******/
			Rij = randFloats(-1, 1)
			OnLooker[i].code[param2change] = NectarSource[i].code[param2change] + Rij*(NectarSource[i].code[param2change]-NectarSource[k].code[param2change])
			/*******判断是否越界********/
			if OnLooker[i].code[param2change] > ub {
				OnLooker[i].code[param2change] = ub
			}
			if OnLooker[i].code[param2change] < lb {
				OnLooker[i].code[param2change] = lb
			}
			OnLooker[i].trueFit = calculationTruefit(OnLooker[i])
			OnLooker[i].fitness = calculationFitness(OnLooker[i].trueFit)

			/******贪婪选择策略*******/
			if OnLooker[i].trueFit < NectarSource[i].trueFit {
				for j := 0; j < D; j++ {
					NectarSource[i].code[j] = OnLooker[i].code[j]
				}
				NectarSource[i].trail = 0
				NectarSource[i].trueFit = OnLooker[i].trueFit
				NectarSource[i].fitness = OnLooker[i].fitness
			} else {
				NectarSource[i].trail++
			}
		}
		i++
		if i == FoodNumber {
			i = 0
		}
	}
}

func sendScoutBees() {
	maxtrialindex := 0
	var R float64 //[0,1]之间的随机数
	maxtrialindex = 0
	for i := 1; i < FoodNumber; i++ {
		if NectarSource[i].trail > NectarSource[maxtrialindex].trail {
			maxtrialindex = i
		}
	}
	if NectarSource[maxtrialindex].trail >= limit {
		/*******重新初始化*********/
		for j := 0; j < D; j++ {
			R = randFloats(0, 1)
			NectarSource[maxtrialindex].code[j] = lb + R*(ub-lb)
		}
		NectarSource[maxtrialindex].trail = 0
		NectarSource[maxtrialindex].trueFit = calculationTruefit(NectarSource[maxtrialindex])
		NectarSource[maxtrialindex].fitness = calculationFitness(NectarSource[maxtrialindex].trueFit)
	}
}
