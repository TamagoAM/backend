package engine

import "math"

// ComputeSocialSatis computes social satisfaction.
// Formula: 80% inverted-boredom + 20% friend-score + modifier additive.
func ComputeSocialSatis(stats *TamaStats, mods *StatModifiers, friends *FriendContext) float64 {
	invertedBoredom := StatMax - stats.Boredom

	friendCoeff := 0.0
	if friends != nil {
		totalFriends := friends.AliveFriends + friends.DeadFriends
		if totalFriends == 0 {
			friendCoeff = NoFriendsMalus
		} else {
			alive := float64(friends.AliveFriends) * FriendAliveBonus
			if alive > FriendBonusCap {
				alive = FriendBonusCap
			}
			friendCoeff = alive - float64(friends.DeadFriends)*FriendDeadMalus
		}
	}

	friendScore := Clamp(50 + friendCoeff*2)
	raw := 0.80*invertedBoredom + 0.20*friendScore

	modAdj := GetAdditive(mods, "socialSatis")
	return Clamp(math.Round((raw+modAdj)*100) / 100)
}

// ComputeWorkSatis computes work satisfaction.
// Formula: 45% money-score + 30% work-recency + 25% safety-record.
func ComputeWorkSatis(stats *TamaStats, mods *StatModifiers) float64 {
	moneyRatio := math.Min(float64(stats.Money), WorkSatisMoneyCAP) / WorkSatisMoneyCAP
	moneyScore := math.Sqrt(moneyRatio) * StatMax

	workEffort := math.Min(100, float64(stats.Worked)*15)
	safetyRecord := Clamp(100 - float64(stats.WorkAccident)*15)

	raw := 0.45*moneyScore + 0.30*workEffort + 0.25*safetyRecord
	modAdj := GetAdditive(mods, "workSatis")
	return Clamp(math.Round((raw+modAdj)*100) / 100)
}

// ComputePersonalSatis computes personal satisfaction.
// Formula: 50% hunger + 50% hygiene + modifier additive.
func ComputePersonalSatis(stats *TamaStats, mods *StatModifiers) float64 {
	raw := 0.50*stats.Hunger + 0.50*stats.Hygiene
	modAdj := GetAdditive(mods, "personalSatis")
	return Clamp(math.Round((raw+modAdj)*100) / 100)
}

// ComputeAllSatisfaction computes all three satisfaction values.
func ComputeAllSatisfaction(stats *TamaStats, mods *StatModifiers, friends *FriendContext) (social, work, personal float64) {
	social = ComputeSocialSatis(stats, mods, friends)
	work = ComputeWorkSatis(stats, mods)
	personal = ComputePersonalSatis(stats, mods)
	return
}

// HappinessFromSatis derives happiness from the three satisfaction values.
func HappinessFromSatis(social, work, personal float64) float64 {
	return Clamp(math.Round(((social+work+personal)/3)*100) / 100)
}
