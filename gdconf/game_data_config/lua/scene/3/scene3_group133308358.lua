-- 基础信息
local base_info = {
	group_id = 133308358
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
	{ config_id = 358001, gadget_id = 70300086, pos = { x = -1721.654, y = 184.781, z = 4691.986 }, rot = { x = 1.613, y = 26.564, z = 10.231 }, level = 32, area_id = 26 },
	{ config_id = 358002, gadget_id = 70300086, pos = { x = -1717.096, y = 175.066, z = 4724.911 }, rot = { x = 12.253, y = 0.477, z = 330.564 }, level = 32, area_id = 26 },
	{ config_id = 358003, gadget_id = 70300086, pos = { x = -1745.162, y = 186.730, z = 4706.634 }, rot = { x = 358.590, y = 6.578, z = 330.803 }, level = 32, area_id = 26 }
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
		gadgets = { 358001, 358002, 358003 },
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