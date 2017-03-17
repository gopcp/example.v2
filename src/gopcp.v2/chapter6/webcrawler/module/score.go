package module

// CalculateScore 代表用于计算组件评分的函数类型。
type CalculateScore func(counts Counts) uint64

// CalculateScoreSimple 代表简易的组件评分计算函数。
func CalculateScoreSimple(counts Counts) uint64 {
	return counts.CalledCount +
		counts.AcceptedCount<<1 +
		counts.CompletedCount<<2 +
		counts.HandlingNumber<<4
}

// SetScore 用于设置给定组件的评分。
// 结果值代表是否更新了评分。
func SetScore(module Module) bool {
	calculator := module.ScoreCalculator()
	if calculator == nil {
		calculator = CalculateScoreSimple
	}
	newScore := calculator(module.Counts())
	if newScore == module.Score() {
		return false
	}
	module.SetScore(newScore)
	return true
}
