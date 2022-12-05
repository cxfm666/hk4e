-- 基础信息
local base_info = {
	group_id = 133108170
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
}

-- 装置
gadgets = {
	{ config_id = 170001, gadget_id = 70220042, pos = { x = -267.116, y = 200.967, z = -861.900 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 2, area_id = 7 }
}

-- 区域
regions = {
}

-- 触发器
triggers = {
}

-- 点位
points = {
	{ config_id = 170004, pos = { x = -273.982, y = 200.638, z = -863.593 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 7, tag = 4 },
	{ config_id = 170005, pos = { x = -269.126, y = 200.531, z = -868.311 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 7, tag = 4 },
	{ config_id = 170006, pos = { x = -263.131, y = 200.967, z = -858.357 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 7, tag = 4 },
	{ config_id = 170007, pos = { x = -267.801, y = 201.000, z = -855.381 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 7, tag = 4 },
	{ config_id = 170008, pos = { x = -272.276, y = 200.647, z = -857.428 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 7, tag = 4 }
}

-- 变量
variables = {
}

-- 怪物随机池
monster_pools = {
	{ pool_id = 1004, rand_weight = 100 },
	{ pool_id = 1005, rand_weight = 100 },
	{ pool_id = 1006, rand_weight = 100 },
	{ pool_id = 1007, rand_weight = 100 }
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
		rand_weight = 100
	}
}

--================================================================
-- 
-- 触发器
-- 
--================================================================

require "TreasureMapEvent"