package game

import (
	"math"
	"time"

	"hk4e/common/constant"
	"hk4e/common/mq"
	"hk4e/gdconf"
	"hk4e/gs/model"
	"hk4e/pkg/alg"
	"hk4e/pkg/logger"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"
)

// 世界管理器

type WorldManager struct {
	worldMap         map[uint32]*World
	snowflake        *alg.SnowflakeWorker
	aiWorld          *World                     // 本服的Ai玩家世界
	sceneBlockAoiMap map[uint32]*alg.AoiManager // 全局各场景地图的aoi管理器
}

func NewWorldManager(snowflake *alg.SnowflakeWorker) (r *WorldManager) {
	r = new(WorldManager)
	r.worldMap = make(map[uint32]*World)
	r.snowflake = snowflake
	r.sceneBlockAoiMap = make(map[uint32]*alg.AoiManager)
	for _, sceneConfig := range gdconf.CONF.SceneMap {
		minX := int16(0)
		maxX := int16(0)
		minZ := int16(0)
		maxZ := int16(0)
		blockXLen := int16(0)
		blockYLen := int16(0)
		blockZLen := int16(0)
		ok := true
		for _, blockConfig := range sceneConfig.BlockMap {
			if int16(blockConfig.BlockRange.Min.X) < minX {
				minX = int16(blockConfig.BlockRange.Min.X)
			}
			if int16(blockConfig.BlockRange.Max.X) > maxX {
				maxX = int16(blockConfig.BlockRange.Max.X)
			}
			if int16(blockConfig.BlockRange.Min.Z) < minZ {
				minZ = int16(blockConfig.BlockRange.Min.Z)
			}
			if int16(blockConfig.BlockRange.Max.Z) > maxZ {
				maxZ = int16(blockConfig.BlockRange.Max.Z)
			}
			xLen := int16(blockConfig.BlockRange.Max.X - blockConfig.BlockRange.Min.X)
			yLen := int16(blockConfig.BlockRange.Max.Y - blockConfig.BlockRange.Min.Y)
			zLen := int16(blockConfig.BlockRange.Max.Z - blockConfig.BlockRange.Min.Z)
			if blockXLen == 0 {
				blockXLen = xLen
			} else {
				if blockXLen != xLen {
					ok = false
					break
				}
			}
			if blockYLen == 0 {
				blockYLen = yLen
			} else {
				if blockYLen != yLen {
					ok = false
					break
				}
			}
			if blockZLen == 0 {
				blockZLen = zLen
			} else {
				if blockZLen != zLen {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		numX := int16(0)
		if blockXLen != 0 {
			if blockXLen > 32 {
				blockXLen = 32
			}
			numX = (maxX - minX) / blockXLen
		} else {
			numX = 1
		}
		if numX == 0 {
			numX = 1
		}
		numZ := int16(0)
		if blockZLen != 0 {
			if blockZLen > 32 {
				blockZLen = 32
			}
			numZ = (maxZ - minZ) / blockZLen
		} else {
			numZ = 1
		}
		if numZ == 0 {
			numZ = 1
		}
		aoiManager := alg.NewAoiManager()
		aoiManager.SetAoiRange(minX, maxX, -1.0, 1.0, minZ, maxZ)
		aoiManager.Init3DRectAoiManager(numX, 1, numZ)
		for _, blockConfig := range sceneConfig.BlockMap {
			for _, groupConfig := range blockConfig.GroupMap {
				for _, monsterConfig := range groupConfig.MonsterList {
					aoiManager.AddObjectToGridByPos(r.snowflake.GenId(), monsterConfig,
						float32(monsterConfig.Pos.X),
						float32(0.0),
						float32(monsterConfig.Pos.Z))
				}
				for _, npcConfig := range groupConfig.NpcList {
					aoiManager.AddObjectToGridByPos(r.snowflake.GenId(), npcConfig,
						float32(npcConfig.Pos.X),
						float32(0.0),
						float32(npcConfig.Pos.Z))
				}
				for _, gadgetConfig := range groupConfig.GadgetList {
					aoiManager.AddObjectToGridByPos(r.snowflake.GenId(), gadgetConfig,
						float32(gadgetConfig.Pos.X),
						float32(0.0),
						float32(gadgetConfig.Pos.Z))
				}
			}
		}
		if sceneConfig.Id == 3 {
			logger.Info("init scene aoi mgr, scene: %v", sceneConfig.Id)
			aoiManager.AoiInfoLog(false)
		}
		r.sceneBlockAoiMap[uint32(sceneConfig.Id)] = aoiManager
	}
	return r
}

func (w *WorldManager) GetWorldByID(worldId uint32) *World {
	return w.worldMap[worldId]
}

func (w *WorldManager) GetAllWorld() map[uint32]*World {
	return w.worldMap
}

func (w *WorldManager) CreateWorld(owner *model.Player) *World {
	worldId := uint32(w.snowflake.GenId())
	world := &World{
		id:                  worldId,
		owner:               owner,
		playerMap:           make(map[uint32]*model.Player),
		sceneMap:            make(map[uint32]*Scene),
		entityIdCounter:     0,
		worldLevel:          0,
		multiplayer:         false,
		mpLevelEntityId:     0,
		chatMsgList:         make([]*proto.ChatInfo, 0),
		playerFirstEnterMap: make(map[uint32]int64),
		waitEnterPlayerMap:  make(map[uint32]int64),
		multiplayerTeam:     CreateMultiplayerTeam(),
		peerList:            make([]*model.Player, 0),
	}
	world.mpLevelEntityId = world.GetNextWorldEntityId(constant.EntityIdTypeConst.MPLEVEL)
	w.worldMap[worldId] = world
	return world
}

func (w *WorldManager) DestroyWorld(worldId uint32) {
	world := w.GetWorldByID(worldId)
	for _, player := range world.playerMap {
		world.RemovePlayer(player)
		player.WorldId = 0
	}
	delete(w.worldMap, worldId)
}

// GetAiWorld 获取本服务器的Ai世界
func (w *WorldManager) GetAiWorld() *World {
	return w.aiWorld
}

// InitAiWorld 初始化Ai世界
func (w *WorldManager) InitAiWorld(owner *model.Player) {
	w.aiWorld = w.GetWorldByID(owner.WorldId)
	w.aiWorld.ChangeToMultiplayer()
	go RunPlayAudio()
}

func (w *WorldManager) IsAiWorld(world *World) bool {
	return world.id == w.aiWorld.id
}

func (w *WorldManager) IsRobotWorld(world *World) bool {
	return world.owner.PlayerID < 100000000
}

func (w *WorldManager) IsBigWorld(world *World) bool {
	return (world.id == w.aiWorld.id) && (w.aiWorld.owner.PlayerID == BigWorldAiUid)
}

// 世界数据结构

type World struct {
	id                  uint32
	owner               *model.Player
	playerMap           map[uint32]*model.Player
	sceneMap            map[uint32]*Scene
	entityIdCounter     uint32 // 世界的实体id生成计数器
	worldLevel          uint8  // 世界等级
	multiplayer         bool   // 是否多人世界
	mpLevelEntityId     uint32
	chatMsgList         []*proto.ChatInfo // 世界聊天消息列表
	playerFirstEnterMap map[uint32]int64  // 玩家第一次进入世界的时间 key:uid value:进入时间
	waitEnterPlayerMap  map[uint32]int64  // 进入世界的玩家等待列表 key:uid value:开始时间
	multiplayerTeam     *MultiplayerTeam
	peerList            []*model.Player // 玩家编号列表
}

func (w *World) GetAllPlayer() map[uint32]*model.Player {
	return w.playerMap
}

func (w *World) GetAllScene() map[uint32]*Scene {
	return w.sceneMap
}

func (w *World) GetNextWorldEntityId(entityType uint16) uint32 {
	for {
		w.entityIdCounter++
		ret := (uint32(entityType) << 24) + w.entityIdCounter
		reTry := false
		for _, scene := range w.sceneMap {
			_, exist := scene.entityMap[ret]
			if exist {
				reTry = true
				break
			}
		}
		if reTry {
			continue
		} else {
			return ret
		}
	}
}

// GetPlayerPeerId 获取当前玩家世界内编号
func (w *World) GetPlayerPeerId(player *model.Player) uint32 {
	peerId := uint32(0)
	for peerIdIndex, worldPlayer := range w.peerList {
		if worldPlayer.PlayerID == player.PlayerID {
			peerId = uint32(peerIdIndex) + 1
		}
	}
	// logger.Debug("get player peer id is: %v, uid: %v", peerId, player.PlayerID)
	return peerId
}

// GetPlayerByPeerId 通过世界内编号获取玩家
func (w *World) GetPlayerByPeerId(peerId uint32) *model.Player {
	peerIdIndex := int(peerId) - 1
	if peerIdIndex >= len(w.peerList) {
		return nil
	}
	return w.peerList[peerIdIndex]
}

// GetWorldPlayerNum 获取世界中玩家的数量
func (w *World) GetWorldPlayerNum() int {
	return len(w.playerMap)
}

func (w *World) AddPlayer(player *model.Player, sceneId uint32) {
	w.peerList = append(w.peerList, player)
	w.playerMap[player.PlayerID] = player
	// 将玩家自身当前的队伍角色信息复制到世界的玩家本地队伍
	team := player.TeamConfig.GetActiveTeam()
	if player.PlayerID == w.owner.PlayerID {
		w.SetPlayerLocalTeam(player, team.GetAvatarIdList())
	} else {
		activeAvatarId := player.TeamConfig.GetActiveAvatarId()
		w.SetPlayerLocalTeam(player, []uint32{activeAvatarId})
	}
	playerNum := w.GetWorldPlayerNum()
	if playerNum > 4 {
		if !WORLD_MANAGER.IsBigWorld(w) {
			return
		}
		w.AddMultiplayerTeam(player)
	} else {
		w.UpdateMultiplayerTeam()
	}
	for _, worldPlayer := range w.playerMap {
		list := w.GetPlayerWorldAvatarList(worldPlayer)
		maxIndex := len(list) - 1
		index := int(worldPlayer.TeamConfig.CurrAvatarIndex)
		if index > maxIndex {
			w.SetPlayerAvatarIndex(worldPlayer, 0)
		} else {
			w.SetPlayerAvatarIndex(worldPlayer, index)
		}
	}
	scene := w.GetSceneById(sceneId)
	scene.AddPlayer(player)
	w.InitPlayerTeamEntityId(player)
}

func (w *World) RemovePlayer(player *model.Player) {
	peerId := w.GetPlayerPeerId(player)
	w.peerList = append(w.peerList[:peerId-1], w.peerList[peerId:]...)
	scene := w.sceneMap[player.SceneId]
	scene.RemovePlayer(player)
	delete(w.playerMap, player.PlayerID)
	delete(w.playerFirstEnterMap, player.PlayerID)
	delete(w.multiplayerTeam.localTeamMap, player.PlayerID)
	delete(w.multiplayerTeam.localAvatarIndexMap, player.PlayerID)
	delete(w.multiplayerTeam.localTeamEntityMap, player.PlayerID)
	playerNum := w.GetWorldPlayerNum()
	if playerNum > 4 {
		if !WORLD_MANAGER.IsBigWorld(w) {
			return
		}
		w.RemoveMultiplayerTeam(player)
	} else {
		if player.PlayerID != w.owner.PlayerID {
			w.UpdateMultiplayerTeam()
		}
	}
}

// WorldAvatar 世界角色
type WorldAvatar struct {
	uid            uint32
	avatarId       uint32
	avatarEntityId uint32
	weaponEntityId uint32
	abilityList    []*proto.AbilityAppliedAbility
	modifierList   []*proto.AbilityAppliedModifier
}

// GetWorldAvatarList 获取世界队伍的全部角色列表
func (w *World) GetWorldAvatarList() []*WorldAvatar {
	worldAvatarList := make([]*WorldAvatar, 0)
	for _, worldAvatar := range w.multiplayerTeam.worldTeam {
		if worldAvatar.uid == 0 {
			continue
		}
		worldAvatarList = append(worldAvatarList, worldAvatar)
	}
	return worldAvatarList
}

// GetPlayerWorldAvatar 获取某玩家在世界队伍中的某角色
func (w *World) GetPlayerWorldAvatar(player *model.Player, avatarId uint32) *WorldAvatar {
	for _, worldAvatar := range w.GetWorldAvatarList() {
		if worldAvatar.uid == player.PlayerID && worldAvatar.avatarId == avatarId {
			return worldAvatar
		}
	}
	return nil
}

// GetPlayerWorldAvatarList 获取某玩家在世界队伍中的所有角色列表
func (w *World) GetPlayerWorldAvatarList(player *model.Player) []*WorldAvatar {
	worldAvatarList := make([]*WorldAvatar, 0)
	for _, worldAvatar := range w.GetWorldAvatarList() {
		if worldAvatar.uid == player.PlayerID {
			worldAvatarList = append(worldAvatarList, worldAvatar)
		}
	}
	return worldAvatarList
}

// GetWorldAvatarByEntityId 通过场景实体id获取世界队伍中的角色
func (w *World) GetWorldAvatarByEntityId(avatarEntityId uint32) *WorldAvatar {
	for _, worldAvatar := range w.GetWorldAvatarList() {
		if worldAvatar.avatarEntityId == avatarEntityId {
			return worldAvatar
		}
	}
	return nil
}

// InitPlayerWorldAvatar 初始化某玩家在世界队伍中的所有角色
func (w *World) InitPlayerWorldAvatar(player *model.Player) {
	scene := w.GetSceneById(player.SceneId)
	for _, worldAvatar := range w.GetWorldAvatarList() {
		if worldAvatar.uid != player.PlayerID {
			continue
		}
		if !player.SceneJump && (worldAvatar.avatarEntityId != 0 || worldAvatar.weaponEntityId != 0) {
			continue
		}
		worldAvatar.avatarEntityId = scene.CreateEntityAvatar(player, worldAvatar.avatarId)
		worldAvatar.weaponEntityId = scene.CreateEntityWeapon()
	}
}

// GetPlayerTeamEntityId 获取某玩家的本地队伍实体id
func (w *World) GetPlayerTeamEntityId(player *model.Player) uint32 {
	return w.multiplayerTeam.localTeamEntityMap[player.PlayerID]
}

// InitPlayerTeamEntityId 初始化某玩家的本地队伍实体id
func (w *World) InitPlayerTeamEntityId(player *model.Player) {
	w.multiplayerTeam.localTeamEntityMap[player.PlayerID] = w.GetNextWorldEntityId(constant.EntityIdTypeConst.TEAM)
}

// GetPlayerWorldAvatarEntityId 获取某玩家在世界队伍中的某角色的实体id
func (w *World) GetPlayerWorldAvatarEntityId(player *model.Player, avatarId uint32) uint32 {
	worldAvatar := w.GetPlayerWorldAvatar(player, avatarId)
	if worldAvatar == nil {
		return 0
	}
	return worldAvatar.avatarEntityId
}

// GetPlayerWorldAvatarWeaponEntityId 获取某玩家在世界队伍中的某角色的武器的实体id
func (w *World) GetPlayerWorldAvatarWeaponEntityId(player *model.Player, avatarId uint32) uint32 {
	worldAvatar := w.GetPlayerWorldAvatar(player, avatarId)
	if worldAvatar == nil {
		return 0
	}
	return worldAvatar.weaponEntityId
}

// GetPlayerAvatarIndex 获取某玩家当前角色索引
func (w *World) GetPlayerAvatarIndex(player *model.Player) int {
	return w.multiplayerTeam.localAvatarIndexMap[player.PlayerID]
}

// SetPlayerAvatarIndex 设置某玩家当前角色索引
func (w *World) SetPlayerAvatarIndex(player *model.Player, index int) {
	if index > len(w.GetPlayerLocalTeam(player))-1 {
		return
	}
	w.multiplayerTeam.localAvatarIndexMap[player.PlayerID] = index
}

// GetPlayerActiveAvatarId 获取玩家当前活跃角色id
func (w *World) GetPlayerActiveAvatarId(player *model.Player) uint32 {
	avatarIndex := w.GetPlayerAvatarIndex(player)
	localTeam := w.GetPlayerLocalTeam(player)
	worldTeamAvatar := localTeam[avatarIndex]
	return worldTeamAvatar.avatarId
}

// GetPlayerAvatarIndexByAvatarId 获取玩家某角色的索引
func (w *World) GetPlayerAvatarIndexByAvatarId(player *model.Player, avatarId uint32) int {
	localTeam := w.GetPlayerLocalTeam(player)
	for index, worldAvatar := range localTeam {
		if worldAvatar.avatarId == avatarId {
			return index
		}
	}
	return -1
}

type MultiplayerTeam struct {
	// key:uid value:玩家的本地队伍
	localTeamMap map[uint32][]*WorldAvatar
	// key:uid value:玩家当前角色索引
	localAvatarIndexMap map[uint32]int
	localTeamEntityMap  map[uint32]uint32
	// 最终的世界队伍
	worldTeam []*WorldAvatar
}

func CreateMultiplayerTeam() (r *MultiplayerTeam) {
	r = new(MultiplayerTeam)
	r.localTeamMap = make(map[uint32][]*WorldAvatar)
	r.localAvatarIndexMap = make(map[uint32]int)
	r.localTeamEntityMap = make(map[uint32]uint32)
	r.worldTeam = make([]*WorldAvatar, 0)
	return r
}

func (w *World) GetPlayerLocalTeam(player *model.Player) []*WorldAvatar {
	return w.multiplayerTeam.localTeamMap[player.PlayerID]
}

func (w *World) SetPlayerLocalTeam(player *model.Player, avatarIdList []uint32) {
	oldLocalTeam := w.multiplayerTeam.localTeamMap[player.PlayerID]
	sameAvatarIdList := make([]uint32, 0)
	diffAvatarIdList := make([]uint32, 0)
	for _, avatarId := range avatarIdList {
		exist := false
		for _, worldAvatar := range oldLocalTeam {
			if worldAvatar.avatarId == avatarId {
				exist = true
			}
		}
		if exist {
			sameAvatarIdList = append(sameAvatarIdList, avatarId)
		} else {
			diffAvatarIdList = append(diffAvatarIdList, avatarId)
		}
	}
	newLocalTeam := make([]*WorldAvatar, len(avatarIdList))
	for _, avatarId := range sameAvatarIdList {
		for _, worldAvatar := range oldLocalTeam {
			if worldAvatar.avatarId == avatarId {
				index := 0
				for i, v := range avatarIdList {
					if avatarId == v {
						index = i
					}
				}
				newLocalTeam[index] = worldAvatar
			}
		}
	}
	for _, avatarId := range diffAvatarIdList {
		index := 0
		for i, v := range avatarIdList {
			if avatarId == v {
				index = i
			}
		}
		newLocalTeam[index] = &WorldAvatar{
			uid:            player.PlayerID,
			avatarId:       avatarId,
			avatarEntityId: 0,
			weaponEntityId: 0,
			abilityList:    make([]*proto.AbilityAppliedAbility, 0),
			modifierList:   make([]*proto.AbilityAppliedModifier, 0),
		}
	}
	w.multiplayerTeam.localTeamMap[player.PlayerID] = newLocalTeam
}

func (w *World) copyLocalTeamToWorld(start int, end int, peerId uint32) {
	player := w.GetPlayerByPeerId(peerId)
	localTeam := w.GetPlayerLocalTeam(player)
	localTeamIndex := 0
	for index := start; index <= end; index++ {
		if localTeamIndex >= len(localTeam) {
			w.multiplayerTeam.worldTeam[index] = &WorldAvatar{
				uid:            0,
				avatarId:       0,
				avatarEntityId: 0,
				weaponEntityId: 0,
				abilityList:    nil,
				modifierList:   nil,
			}
			continue
		}
		w.multiplayerTeam.worldTeam[index] = localTeam[localTeamIndex]
		localTeamIndex++
	}
}

// TODO 为了实现大世界无限人数写的
// 现在看来把世界里所有人放进队伍里发给客户端超过8个客户端会崩溃
// 看来还是不能简单的走通用逻辑 需要对大世界场景队伍做特殊处理 欺骗客户端其他玩家仅仅以场景角色实体的形式出现

func (w *World) AddMultiplayerTeam(player *model.Player) {
	if !WORLD_MANAGER.IsBigWorld(w) {
		return
	}
	localTeam := w.GetPlayerLocalTeam(player)
	w.multiplayerTeam.worldTeam = append(w.multiplayerTeam.worldTeam, localTeam...)
}

func (w *World) RemoveMultiplayerTeam(player *model.Player) {
	worldTeam := make([]*WorldAvatar, 0)
	for _, worldAvatar := range w.multiplayerTeam.worldTeam {
		if worldAvatar.uid == player.PlayerID {
			continue
		}
		worldTeam = append(worldTeam, worldAvatar)
	}
	w.multiplayerTeam.worldTeam = worldTeam
}

// UpdateMultiplayerTeam 整合所有玩家的本地队伍计算出世界队伍
func (w *World) UpdateMultiplayerTeam() {
	playerNum := w.GetWorldPlayerNum()
	if playerNum > 4 {
		return
	}
	w.multiplayerTeam.worldTeam = make([]*WorldAvatar, 4)
	switch playerNum {
	case 1:
		// 1P*4
		w.copyLocalTeamToWorld(0, 3, 1)
	case 2:
		// 1P*2 + 2P*2
		w.copyLocalTeamToWorld(0, 1, 1)
		w.copyLocalTeamToWorld(2, 3, 2)
	case 3:
		// 1P*2 + 2P*1 + 3P*1
		w.copyLocalTeamToWorld(0, 1, 1)
		w.copyLocalTeamToWorld(2, 2, 2)
		w.copyLocalTeamToWorld(3, 3, 3)
	case 4:
		// 1P*1 + 2P*1 + 3P*1 + 4P*1
		w.copyLocalTeamToWorld(0, 0, 1)
		w.copyLocalTeamToWorld(1, 1, 2)
		w.copyLocalTeamToWorld(2, 2, 3)
		w.copyLocalTeamToWorld(3, 3, 4)
	}
}

// 世界聊天

func (w *World) AddChat(chatInfo *proto.ChatInfo) {
	w.chatMsgList = append(w.chatMsgList, chatInfo)
}

func (w *World) GetChatList() []*proto.ChatInfo {
	return w.chatMsgList
}

// ChangeToMultiplayer 转换为多人世界
func (w *World) ChangeToMultiplayer() {
	w.multiplayer = true
}

// IsPlayerFirstEnter 获取玩家是否首次加入本世界
func (w *World) IsPlayerFirstEnter(player *model.Player) bool {
	_, exist := w.playerFirstEnterMap[player.PlayerID]
	if !exist {
		return true
	} else {
		return false
	}
}

func (w *World) PlayerEnter(player *model.Player) {
	w.playerFirstEnterMap[player.PlayerID] = time.Now().UnixMilli()
}

func (w *World) CreateScene(sceneId uint32) *Scene {
	scene := &Scene{
		id:                sceneId,
		world:             w,
		playerMap:         make(map[uint32]*model.Player),
		entityMap:         make(map[uint32]*Entity),
		objectIdEntityMap: make(map[int64]*Entity),
		gameTime:          18 * 60,
		createTime:        time.Now().UnixMilli(),
		meeoIndex:         0,
	}
	w.sceneMap[sceneId] = scene
	return scene
}

func (w *World) GetSceneById(sceneId uint32) *Scene {
	scene, exist := w.sceneMap[sceneId]
	if !exist {
		scene = w.CreateScene(sceneId)
	}
	return scene
}

// 场景数据结构

type Scene struct {
	id                uint32
	world             *World
	playerMap         map[uint32]*model.Player
	entityMap         map[uint32]*Entity
	objectIdEntityMap map[int64]*Entity
	gameTime          uint32 // 游戏内提瓦特大陆的时间
	createTime        int64
	meeoIndex         uint32 // 客户端风元素染色同步协议的计数器
}

func (s *Scene) GetAllPlayer() map[uint32]*model.Player {
	return s.playerMap
}

func (s *Scene) GetAllEntity() map[uint32]*Entity {
	return s.entityMap
}

type AvatarEntity struct {
	uid      uint32
	avatarId uint32
}

type MonsterEntity struct {
	monsterId uint32
}

type NpcEntity struct {
	NpcId         uint32
	RoomId        uint32
	ParentQuestId uint32
	BlockId       uint32
}

const (
	GADGET_TYPE_NORMAL = iota
	GADGET_TYPE_GATHER
	GADGET_TYPE_CLIENT
	GADGET_TYPE_VEHICLE // 载具
)

type GadgetClientEntity struct {
	configId          uint32
	campId            uint32
	campType          uint32
	ownerEntityId     uint32
	targetEntityId    uint32
	propOwnerEntityId uint32
}

type GadgetGatherEntity struct {
	gatherId uint32
}

type GadgetVehicleEntity struct {
	vehicleId  uint32
	owner      *model.Player
	maxStamina float32
	curStamina float32
	memberMap  map[uint32]*model.Player // uint32 = pos
}

type GadgetEntity struct {
	gadgetType          int
	gadgetId            uint32
	gadgetClientEntity  *GadgetClientEntity
	gadgetGatherEntity  *GadgetGatherEntity
	gadgetVehicleEntity *GadgetVehicleEntity
}

// 场景实体数据结构

type Entity struct {
	id                  uint32
	scene               *Scene
	lifeState           uint16
	pos                 *model.Vector
	rot                 *model.Vector
	moveState           uint16
	lastMoveSceneTimeMs uint32
	lastMoveReliableSeq uint32
	fightProp           map[uint32]float32
	entityType          uint32
	level               uint8
	avatarEntity        *AvatarEntity
	monsterEntity       *MonsterEntity
	npcEntity           *NpcEntity
	gadgetEntity        *GadgetEntity
	configId            uint32
	objectId            int64
}

type Attack struct {
	combatInvokeEntry *proto.CombatInvokeEntry
	uid               uint32
}

func (s *Scene) ChangeGameTime(time uint32) {
	s.gameTime = time % 1440
}

func (s *Scene) GetSceneCreateTime() int64 {
	return s.createTime
}

func (s *Scene) GetSceneTime() int64 {
	now := time.Now().UnixMilli()
	return now - s.createTime
}

func (s *Scene) AddPlayer(player *model.Player) {
	s.playerMap[player.PlayerID] = player
	s.world.InitPlayerWorldAvatar(player)
}

func (s *Scene) RemovePlayer(player *model.Player) {
	delete(s.playerMap, player.PlayerID)
	worldAvatarList := s.world.GetPlayerWorldAvatarList(player)
	for _, worldAvatar := range worldAvatarList {
		s.DestroyEntity(worldAvatar.avatarEntityId)
		s.DestroyEntity(worldAvatar.weaponEntityId)
	}
}

func (s *Scene) SetEntityLifeState(entity *Entity, lifeState uint16, dieType proto.PlayerDieType) {
	if entity.avatarEntity != nil {
		// 获取玩家对象
		player := USER_MANAGER.GetOnlineUser(entity.avatarEntity.uid)
		if player == nil {
			logger.Error("player is nil, uid: %v", entity.avatarEntity.uid)
			return
		}
		// 获取角色
		avatar, ok := player.AvatarMap[entity.avatarEntity.avatarId]
		if !ok {
			logger.Error("avatar is nil, avatarId: %v", avatar)
			return
		}
		// 设置角色存活状态
		if lifeState == constant.LifeStateConst.LIFE_REVIVE {
			avatar.LifeState = constant.LifeStateConst.LIFE_ALIVE
			// 设置血量
			entity.fightProp[uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP)] = 110
			GAME_MANAGER.EntityFightPropUpdateNotifyBroadcast(s, entity, uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP))
		}

		// PacketAvatarLifeStateChangeNotify
		avatarLifeStateChangeNotify := &proto.AvatarLifeStateChangeNotify{
			LifeState:       uint32(lifeState),
			AttackTag:       "",
			DieType:         dieType,
			ServerBuffList:  nil,
			MoveReliableSeq: entity.lastMoveReliableSeq,
			SourceEntityId:  0,
			AvatarGuid:      avatar.Guid,
		}
		for _, p := range s.playerMap {
			GAME_MANAGER.SendMsg(cmd.AvatarLifeStateChangeNotify, p.PlayerID, p.ClientSeq, avatarLifeStateChangeNotify)
		}
	} else {
		// 设置存活状态
		entity.lifeState = lifeState

		if lifeState == constant.LifeStateConst.LIFE_DEAD {
			// 设置血量
			entity.fightProp[uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP)] = 0
			GAME_MANAGER.EntityFightPropUpdateNotifyBroadcast(s, entity, uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP))
		}

		// PacketLifeStateChangeNotify
		lifeStateChangeNotify := &proto.LifeStateChangeNotify{
			EntityId:        entity.id,
			AttackTag:       "",
			MoveReliableSeq: entity.lastMoveReliableSeq,
			DieType:         dieType,
			LifeState:       uint32(lifeState),
			SourceEntityId:  0,
		}
		for _, p := range s.playerMap {
			GAME_MANAGER.SendMsg(cmd.LifeStateChangeNotify, p.PlayerID, p.ClientSeq, lifeStateChangeNotify)
		}

		// 删除实体
		s.DestroyEntity(entity.id)
		GAME_MANAGER.RemoveSceneEntityNotifyBroadcast(s, proto.VisionType_VISION_DIE, []uint32{entity.id})
	}
}

func (s *Scene) CreateEntityAvatar(player *model.Player, avatarId uint32) uint32 {
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.AVATAR)
	avatar, ok := player.AvatarMap[avatarId]
	if !ok {
		logger.Error("avatar error, avatarId: %v", avatar)
		return 0
	}
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           avatar.LifeState,
		pos:                 player.Pos,
		rot:                 player.Rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		// fightProp:           player.AvatarMap[avatarId].FightPropMap, // 使用角色结构的数据
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_AVATAR),
		// level:               0, // 使用角色结构的数据
		avatarEntity: &AvatarEntity{
			uid:      player.PlayerID,
			avatarId: avatarId,
		},
	}
	s.entityMap[entity.id] = entity
	MESSAGE_QUEUE.SendToFight(s.world.owner.FightAppId, &mq.NetMsg{
		MsgType: mq.MsgTypeFight,
		EventId: mq.FightRoutineAddEntity,
		FightMsg: &mq.FightMsg{
			FightRoutineId: s.world.id,
			EntityId:       entity.id,
			FightPropMap:   entity.fightProp,
			Uid:            entity.avatarEntity.uid,
			AvatarGuid:     player.AvatarMap[avatarId].Guid,
		},
	})
	return entity.id
}

func (s *Scene) CreateEntityWeapon() uint32 {
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.WEAPON)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 new(model.Vector),
		rot:                 new(model.Vector),
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp:           nil,
		entityType:          uint32(proto.ProtEntityType_PROT_ENTITY_WEAPON),
		level:               0,
	}
	s.entityMap[entity.id] = entity
	return entity.id
}

func (s *Scene) CreateEntityMonster(pos, rot *model.Vector, monsterId uint32, level uint8, fightProp map[uint32]float32, configId uint32, objectId int64) uint32 {
	_, exist := s.objectIdEntityMap[objectId]
	if exist {
		return 0
	}
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.MONSTER)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp:           fightProp,
		entityType:          uint32(proto.ProtEntityType_PROT_ENTITY_MONSTER),
		level:               level,
		monsterEntity: &MonsterEntity{
			monsterId: monsterId,
		},
		configId: configId,
		objectId: objectId,
	}
	s.entityMap[entity.id] = entity
	s.objectIdEntityMap[objectId] = entity
	MESSAGE_QUEUE.SendToFight(s.world.owner.FightAppId, &mq.NetMsg{
		MsgType: mq.MsgTypeFight,
		EventId: mq.FightRoutineAddEntity,
		FightMsg: &mq.FightMsg{
			FightRoutineId: s.world.id,
			EntityId:       entity.id,
			FightPropMap:   entity.fightProp,
		},
	})
	return entity.id
}

func (s *Scene) CreateEntityNpc(pos, rot *model.Vector, npcId, roomId, parentQuestId, blockId, configId uint32, objectId int64) uint32 {
	_, exist := s.objectIdEntityMap[objectId]
	if exist {
		return 0
	}
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.NPC)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp: map[uint32]float32{
			uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_MAX_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_BASE_HP): float32(1),
		},
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_NPC),
		level:      0,
		npcEntity: &NpcEntity{
			NpcId:         npcId,
			RoomId:        roomId,
			ParentQuestId: parentQuestId,
			BlockId:       blockId,
		},
		configId: configId,
		objectId: objectId,
	}
	s.entityMap[entity.id] = entity
	s.objectIdEntityMap[objectId] = entity
	return entity.id
}

func (s *Scene) CreateEntityGadgetNormal(pos, rot *model.Vector, gadgetId uint32, configId uint32, objectId int64) uint32 {
	_, exist := s.objectIdEntityMap[objectId]
	if exist {
		return 0
	}
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.GADGET)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp: map[uint32]float32{
			uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_MAX_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_BASE_HP): float32(1),
		},
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_GADGET),
		level:      0,
		gadgetEntity: &GadgetEntity{
			gadgetId:   gadgetId,
			gadgetType: GADGET_TYPE_NORMAL,
		},
		configId: configId,
		objectId: objectId,
	}
	s.entityMap[entity.id] = entity
	s.objectIdEntityMap[objectId] = entity
	return entity.id
}

func (s *Scene) CreateEntityGadgetGather(pos, rot *model.Vector, gadgetId uint32, gatherId uint32, configId uint32, objectId int64) uint32 {
	_, exist := s.objectIdEntityMap[objectId]
	if exist {
		return 0
	}
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.GADGET)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp: map[uint32]float32{
			uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_MAX_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_BASE_HP): float32(1),
		},
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_GADGET),
		level:      0,
		gadgetEntity: &GadgetEntity{
			gadgetId:   gadgetId,
			gadgetType: GADGET_TYPE_GATHER,
			gadgetGatherEntity: &GadgetGatherEntity{
				gatherId: gatherId,
			},
		},
		configId: configId,
		objectId: objectId,
	}
	s.entityMap[entity.id] = entity
	s.objectIdEntityMap[objectId] = entity
	return entity.id
}

func (s *Scene) CreateEntityGadgetClient(pos, rot *model.Vector, entityId uint32, configId, campId, campType, ownerEntityId, targetEntityId, propOwnerEntityId uint32) {
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp: map[uint32]float32{
			uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_MAX_HP):  math.MaxFloat32,
			uint32(constant.FightPropertyConst.FIGHT_PROP_BASE_HP): float32(1),
		},
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_GADGET),
		level:      0,
		gadgetEntity: &GadgetEntity{
			gadgetType: GADGET_TYPE_CLIENT,
			gadgetClientEntity: &GadgetClientEntity{
				configId:          configId,
				campId:            campId,
				campType:          campType,
				ownerEntityId:     ownerEntityId,
				targetEntityId:    targetEntityId,
				propOwnerEntityId: propOwnerEntityId,
			},
		},
	}
	s.entityMap[entity.id] = entity
}

func (s *Scene) CreateEntityGadgetVehicle(uid uint32, pos, rot *model.Vector, vehicleId uint32) uint32 {
	player := USER_MANAGER.GetOnlineUser(uid)
	if player == nil {
		logger.Error("player is nil, uid: %v", uid)
		return 0
	}
	entityId := s.world.GetNextWorldEntityId(constant.EntityIdTypeConst.GADGET)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		lifeState:           constant.LifeStateConst.LIFE_ALIVE,
		pos:                 pos,
		rot:                 rot,
		moveState:           uint16(proto.MotionState_MOTION_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp: map[uint32]float32{
			// TODO 以后使用配置表
			uint32(constant.FightPropertyConst.FIGHT_PROP_CUR_HP):  114514,
			uint32(constant.FightPropertyConst.FIGHT_PROP_MAX_HP):  114514,
			uint32(constant.FightPropertyConst.FIGHT_PROP_BASE_HP): float32(1),
		},
		entityType: uint32(proto.ProtEntityType_PROT_ENTITY_GADGET),
		level:      0,
		gadgetEntity: &GadgetEntity{
			gadgetType: GADGET_TYPE_VEHICLE,
			gadgetVehicleEntity: &GadgetVehicleEntity{
				vehicleId:  vehicleId,
				owner:      player,
				maxStamina: 240, // TODO 应该也能在配置表找到
				curStamina: 240, // TODO 与maxStamina一致
				memberMap:  make(map[uint32]*model.Player),
			},
		},
	}
	s.entityMap[entity.id] = entity
	return entity.id
}

func (s *Scene) DestroyEntity(entityId uint32) {
	entity := s.GetEntity(entityId)
	if entity == nil {
		return
	}
	delete(s.entityMap, entity.id)
	delete(s.objectIdEntityMap, entity.objectId)
	MESSAGE_QUEUE.SendToFight(s.world.owner.FightAppId, &mq.NetMsg{
		MsgType: mq.MsgTypeFight,
		EventId: mq.FightRoutineDelEntity,
		FightMsg: &mq.FightMsg{
			FightRoutineId: s.world.id,
			EntityId:       entity.id,
		},
	})
}

func (s *Scene) GetEntity(entityId uint32) *Entity {
	return s.entityMap[entityId]
}

func (s *Scene) GetEntityByObjectId(objectId int64) *Entity {
	return s.objectIdEntityMap[objectId]
}
