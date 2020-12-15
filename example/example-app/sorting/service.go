package sorting

type SortingService interface {
	BubbleSort(numbers []int)
}

type sortingService struct{}

func NewService() *sortingService {
	return &sortingService{}
}

func (s *sortingService) BubbleSort(numbers []int) {
	for i := 0; i < len(numbers); i++ {
		for j := 0; j < len(numbers)-1-i; j++ {
			if numbers[j] > numbers[j+1] {
				numbers[j], numbers[j+1] = numbers[j+1], numbers[j]
			}
		}
	}
}
