local uuid = require('uuid')
local console = require('console')
local log = require('log')
local clock = require('clock')

console.listen('127.0.0.1:3312')

box.cfg{
	listen              = 3301,
	memtx_memory        = 100000000,
	pid_file            = 'tarantool.pid',    
	memtx_dir           = 'snap',
	wal_dir             = 'wal',
	checkpoint_interval = 3600,
	checkpoint_count    = 5,
	log         		= 'tarantool.log'
}

local session_life_time = 600


local function bootstrap()

	local s_user = box.schema.create_space('user')
	s_user:format({
		{name = 'user_id', type = 'string'},
		{name = 'user_name', type = 'string'},
		{name = 'pass', type = 'string'}
	})
	s_user:create_index('primary', {
		type = 'tree',
		parts = {'user_id'}
})
	s_user:create_index('name', {
		type = 'tree',
		parts = {'user_name'}
	})

	local s_group = box.schema.create_space('group')
	s_group:format({
		{name = "group_id", type = 'string'},
		{name = "group_name", type = 'string'}
	})
	s_group:create_index('group_id', {
			type = 'tree',
			parts = {'group_id'}
	})
	s_group:create_index('group', {
		type = 'tree',
		unique = false,
		parts = {'group_name'}
	})

	local s_usergroup = box.schema.create_space('usergroup')
	s_usergroup:format({
		{name = "usergroup_id", type = 'string'},
		{name = "user", type = 'string'},
		{name = "group_id", type = 'string'},
		{name = "role", type = 'string'}
	})
	s_usergroup:create_index('primary', {
		type = 'tree',
		parts = {'usergroup_id'}
	})
	s_usergroup:create_index('user', {
		type = 'tree',
		unique = false,
		parts = {'user'}
	})
	s_usergroup:create_index('group_id', {
		type = 'tree',
		unique = false,
		parts = {'group_id'}
	})
	s_usergroup:create_index('role', {
		type = 'tree',
		unique = false,
		parts = {'role'}
	})
	s_usergroup:create_index('user_group', {
		type = 'tree',
		unique = false,
		parts = {'user', 'group_id'}
	})




	local s_msg = box.schema.create_space('msg')
	s_msg:format({
		{name = "msg_id", type = 'string'},
		{name = "message", type = 'string'},
		{name = "group_id", type = 'string'},
		{name = "user", type = 'string'},
		{name = "msg_time", type = 'unsigned'}
	})
	s_msg:create_index('primary', {
		type = 'tree',
		parts = {'msg_id'}
	})

	s_msg:create_index('user', {
		type = 'tree',
		unique = false,
		parts = {'user'}
	})

	s_msg:create_index('time', {
		type = 'tree',
		unique = false,
		parts = {'msg_time'}
	})

	s_msg:create_index('group_id', {
		type = 'tree',
		unique = false,
		parts = {'group_id'}
	})

	s_msg:create_index('user_group', {
		type = 'tree',
		unique = false,
		parts = {'user', 'group_id'}
	})


	box.schema.user.create('ex', { password = 'secret' })
	box.schema.user.grant('ex','read,write,execute,create,drop','universe')

	box.schema.user.create('repl', { password = 'replication' })
	box.schema.user.grant('repl', 'replication')
end

box.once('replica', bootstrap)

local function is_str_empty(str)
	return str == nil or str == ''
end


local function get_group_id(user_id)
	local t_user_id = s_user.index.primary:get{user_id}
	return t_user_id[group_id]
end

fn = require("fn")

console.start()