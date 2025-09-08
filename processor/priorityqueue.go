package processor

import "AdvanceTradeEngine/models"

type PriorityQueue []*models.Order

func (prique PriorityQueue) Len() int {
	return len(prique)
}

func (prique PriorityQueue) Less(i, j int) bool {
	if prique[i].Side == "BUY" {
		if prique[i].Price == prique[j].Price {
			return prique[i].Timestamp.Before(prique[j].Timestamp)
		}
		return prique[i].Price > prique[j].Price
	}
	if prique[i].Price == prique[j].Price {
		return prique[i].Timestamp.Before(prique[j].Timestamp)
	}
	return prique[i].Price < prique[j].Price
}

func (prique PriorityQueue) Swap(i, j int) {
	prique[i], prique[j] = prique[j], prique[i]
}

func (prique *PriorityQueue) Push(x interface{}) {
	order := x.(*models.Order)
	*prique = append(*prique, order)
}

func (prique *PriorityQueue) Pop() interface{} {
	lOld := *prique
	lLen := len(lOld)
	lOrder := lOld[lLen-1]
	*prique = lOld[:lLen-1]
	return lOrder
}
