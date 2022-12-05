-- 基础信息
local base_info = {
	group_id = 133302566
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 566001, monster_id = 21010101, pos = { x = -737.764, y = 164.888, z = 2634.635 }, rot = { x = 0.000, y = 72.769, z = 0.000 }, level = 27, drop_tag = "丘丘人", pose_id = 9003, area_id = 24 },
	{ config_id = 566002, monster_id = 21010101, pos = { x = -736.844, y = 164.351, z = 2637.601 }, rot = { x = 0.000, y = 140.030, z = 0.000 }, level = 27, drop_tag = "丘丘人", pose_id = 9012, area_id = 24 }
}

-- NPC
npcs = {
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
		monsters = { 566001, 566002 },
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