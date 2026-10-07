package render

import "testing"

func BenchmarkEarthCutscenes(b *testing.B) {
	for _, shot := range []int{7, 8} {
		b.Run(FilmShots(FilmEnding)[shot].Name, func(b *testing.B) {
			st := CutsceneState{Film: FilmEnding, Shot: shot, Frame: 120, Revealed: true}
			RenderCutscene(160, 54, st)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				st.Frame = 48 + i%150
				RenderCutscene(160, 54, st)
			}
		})
	}
}
