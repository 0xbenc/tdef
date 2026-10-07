package game

import "testing"

func TestTrainingRecruitmentAndSafePlanning(t *testing.T) {
	for _, stage := range []int{1, 2} {
		s := NewState(loadTestMap(t))
		s.ConfigureTraining(stage, 1<<TowerGunner)
		if stage == 1 && (s.TowerAvailable(TowerFrost) || s.TowerAvailable(TowerCannon)) {
			t.Fatal("first defense did not start with only Gunners")
		}
		if stage == 2 && s.UnlockedTowers&MainTowersMask != MainTowersMask {
			t.Fatal("second defense relocked main roster")
		}
		s.LessonPending = false
		for _, r := range RecruitmentOrder(stage) {
			if stage == 2 && r.Kind == TowerRuneforge {
				continue
			}
			s.Wave = r.AfterWave
			s.WaveActive = true
			s.recruitDue()
			if s.TowerAvailable(r.Kind) {
				t.Fatal("recruited during active wave")
			}
			s.WaveActive = false
			before := s.Gold
			s.recruitDue()
			if !s.LessonPending || s.Lesson != r.Kind || !s.TowerAvailable(r.Kind) {
				t.Fatal("missing scheduled recruit")
			}
			if s.Gold-before != TowerSpecs[r.Kind].Cost[0] {
				t.Fatal("recruit not funded")
			}
			clock, gold := s.Time, s.Gold
			s.Step(100)
			s.StartWave()
			if s.Time != clock || s.Gold != gold || s.WaveActive {
				t.Fatal("lesson allowed time or wave advancement")
			}
			s.LessonPending = false
			s.recruitDue()
			if s.Gold != gold {
				t.Fatal("recruit awarded twice")
			}
		}
		s.LessonPending = false
		s.NextWaveAt = 0
		s.Step(.01)
		if !s.WaveActive {
			t.Fatal("training wave did not start automatically")
		}
	}
}

func TestTrainingLockEnforcedByBuildAndRetry(t *testing.T) {
	s := NewState(loadTestMap(t))
	s.Gold = 100000
	s.ConfigureTraining(1, 1<<TowerGunner)
	v := Vec{4, 1}
	if s.Build(v, TowerFrost) != nil {
		t.Fatal("locked unit built")
	}
	if s.Build(v, TowerGunner) == nil {
		t.Fatal("initial unit unavailable")
	}
	s.LessonPending = false
	s.Wave = 1
	s.recruitDue()
	retry := NewState(s.Map)
	retry.ConfigureTraining(1, s.UnlockedTowers)
	if !retry.TowerAvailable(TowerFrost) || retry.TowerAvailable(TowerCannon) || retry.LessonPending {
		t.Fatal("retry lost recruitment progress")
	}
	if NewState(s.Map).TowerAvailable(TowerWitch) == false {
		t.Fatal("ordinary games restricted")
	}
}
