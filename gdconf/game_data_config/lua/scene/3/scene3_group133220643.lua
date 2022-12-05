-- 基础信息
local base_info = {
	group_id = 133220643
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
}

-- NPC
npcs = {
	{ config_id = 643001, npc_id = 30269, pos = { x = -2786.207, y = 224.123, z = -4389.994 }, rot = { x = 0.000, y = 38.116, z = 0.000 }, area_id = 11 },
	{ config_id = 643002, npc_id = 30270, pos = { x = -2806.990, y = 227.148, z = -4470.637 }, rot = { x = 0.000, y = 283.120, z = 0.000 }, area_id = 11 },
	{ config_id = 643003, npc_id = 30271, pos = { x = -2808.441, y = 227.323, z = -4471.201 }, rot = { x = 0.000, y = 76.602, z = 0.000 }, area_id = 11 },
	{ config_id = 643004, npc_id = 30272, pos = { x = -2808.516, y = 227.251, z = -4470.196 }, rot = { x = 0.000, y = 105.812, z = 0.000 }, area_id = 11 },
	{ config_id = 643005, npc_id = 30273, pos = { x = -2807.985, y = 227.203, z = -4469.466 }, rot = { x = 0.000, y = 146.417, z = 0.000 }, area_id = 11 }
}

-- 装置
gadgets = {
}

-- 区域
regions = {
}

-- 触发器
triggers = {
}

-- 变量
variables = {
}

--================================================================
-- 
-- 初始化配置
-- 
--================================================================

-- 初始化时创建
init_config = {
	suite = 1,
	end_suite = 0,
	rand_suite = false
}

--================================================================
-- 
-- 小组配置
-- 
--================================================================

suites = {
	{
		-- suite_id = 1,
		-- description = ,
		monsters = { },
		gadgets = { },
		regions = { },
		triggers = { },
		npcs = { 643001, 643002, 643003, 643004, 643005 },
		rand_weight = 100
	}
}

--================================================================
-- 
-- 触发器
-- 
--================================================================