package pool

type Updater struct {
	update func()
}

func NewUpdater(updateFunc func()) *Updater {
	return &Updater{
		update: updateFunc,
	}
}

func (u *Updater) Update() {
	u.update()
}
