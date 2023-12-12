package internal

import (
	"fmt"
	"strings"
)

type Player struct {
	CurrentRoom     Room
	Inventory       []string
	WearingBackpack bool
	KeysTaken       bool
	NotesTaken      bool
}

type Room interface {
	Describe() string
	HandleAction(action string, player *Player) string
}

type BaseRoom struct {
	Description string
	Actions     map[string]func(player *Player) string
	Inventory   []string
}

func (r *BaseRoom) Describe() string {
	return r.Description
}

func containsItem(inventory []string, item string) bool {
	for _, i := range inventory {
		if i == item {
			return true
		}
	}
	return false
}

func (r *BaseRoom) HandleAction(action string, player *Player) string {
	parts := strings.Fields(action)

	if len(parts) > 1 {
		switch parts[0] {
		case "взять":
			return r.takeAction(parts[1], player)
		case "идти":
			return r.goAction(strings.Join(parts[1:], " "), player)
		case "применить":
			return r.useAction(parts[1], parts[2], player)
		}
	}

	if act, ok := r.Actions[action]; ok {
		return act(player)
	}

	return "неизвестная команда"
}

func (r *BaseRoom) takeAction(item string, player *Player) string {
	if !containsItem(r.Inventory, item) {
		return fmt.Sprintf("нет такого")
	}

	if takeAction, ok := r.Actions[fmt.Sprintf("взять %s", item)]; ok {
		return takeAction(player)
	} else {
		return "нет такого"
	}
}

func (r *BaseRoom) goAction(direction string, player *Player) string {
	if dirAction, ok := r.Actions["идти "+direction]; ok {
		return dirAction(player)
	} else {
		return fmt.Sprintf("нет пути в %s", direction)
	}
}

func (r *BaseRoom) useAction(item, target string, player *Player) string {
	if !containsItem(player.Inventory, item) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", item)
	}

	if useAction, ok := r.Actions[fmt.Sprintf("применить %s %s", item, target)]; ok {
		return useAction(player)
	} else {
		return "не к чему применить"
	}
}

func NewBaseRoom(description string, actions map[string]func(player *Player) string, inventory []string) *BaseRoom {
	return &BaseRoom{
		Description: description,
		Actions:     actions,
		Inventory:   inventory,
	}
}

func NewKitchen() Room {
	return NewBaseRoom("ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор",
		map[string]func(player *Player) string{
			"осмотреться": func(player *Player) string {
				if player.WearingBackpack {
					return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
				}
				return player.CurrentRoom.Describe()
			},
			"идти коридор": func(player *Player) string {
				player.CurrentRoom = NewCorridor()
				return player.CurrentRoom.Describe()
			},
		},
		[]string{})
}

func NewCorridor() Room {
	doorLocked := true
	return NewBaseRoom("ничего интересного. можно пройти - кухня, комната, улица",
		map[string]func(player *Player) string{
			"осмотреться": func(player *Player) string {
				return player.CurrentRoom.Describe()
			},
			"идти кухня": func(player *Player) string {
				player.CurrentRoom = NewKitchen()
				return "кухня, ничего интересного. можно пройти - коридор"
			},
			"идти комната": func(player *Player) string {
				player.CurrentRoom = NewRoom(player)
				return player.CurrentRoom.Describe()
			},
			"идти улица": func(player *Player) string {
				if doorLocked {
					return "дверь закрыта"
				}
				player.CurrentRoom = NewStreet()
				return player.CurrentRoom.Describe()
			},
			"применить ключи дверь": func(player *Player) string {
				if containsItem(player.Inventory, "ключи") {
					doorLocked = false
					return "дверь открыта"
				}
				return "у вас нет ключей"
			},
		},
		[]string{})
}

func NewRoom(player *Player) Room {
	return NewBaseRoom("ты в своей комнате. можно пройти - коридор",
		map[string]func(player *Player) string{
			"осмотреться": func(player *Player) string {
				if player.WearingBackpack && player.KeysTaken && player.NotesTaken {
					return "пустая комната. можно пройти - коридор"
				} else if player.WearingBackpack && player.KeysTaken {
					return "на столе: конспекты. можно пройти - коридор"
				} else if player.WearingBackpack && player.NotesTaken {
					return "на столе: ключи. можно пройти - коридор"
				} else if player.WearingBackpack {
					return "на столе: ключи, конспекты. можно пройти - коридор"
				} else if player.KeysTaken {
					return "на столе: конспекты, на стуле: рюкзак. можно пройти - коридор"
				} else if player.NotesTaken {
					return "на столе: ключи, на стуле: рюкзак. можно пройти - коридор"
				} else {
					return "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"
				}
			},
			"взять ключи": func(player *Player) string {
				if !player.WearingBackpack {
					return "некуда класть"
				}

				if player.KeysTaken {
					return "нет такого"
				}

				player.KeysTaken = true
				player.Inventory = append(player.Inventory, "ключи")
				return "предмет добавлен в инвентарь: ключи"
			},
			"надеть рюкзак": func(player *Player) string {
				if !player.WearingBackpack {
					player.WearingBackpack = true
					return "вы надели: рюкзак"
				}
				return "нет такого"
			},
			"взять конспекты": func(player *Player) string {
				if !player.WearingBackpack {
					return "некуда класть"
				}

				if player.NotesTaken {
					return "нет такого"
				}

				player.NotesTaken = true
				return "предмет добавлен в инвентарь: конспекты"
			},
			"идти коридор": func(player *Player) string {
				player.CurrentRoom = NewCorridor()
				return player.CurrentRoom.Describe()
			},
		},
		[]string{"рюкзак", "ключи", "конспекты"})
}

func NewStreet() Room {
	return NewBaseRoom("на улице весна. можно пройти - домой",
		map[string]func(player *Player) string{
			"осмотреться": func(player *Player) string {
				return player.CurrentRoom.Describe()
			},
			"идти домой": func(player *Player) string {
				player.CurrentRoom = NewRoom(player)
				return "вы идете домой"
			},
		},
		[]string{})
}

func InitGame() *Player {
	player := &Player{
		CurrentRoom: NewKitchen(),
		Inventory:   make([]string, 0),
	}
	return player
}

func HandleCommand(player *Player, command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	action := parts[0]
	if len(parts) > 1 {
		action += " " + strings.Join(parts[1:], " ")
	}

	result := player.CurrentRoom.HandleAction(action, player)

	return result
}
