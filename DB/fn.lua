local uuid = require('uuid')
local digest = require('digest')
local datetime = require('datetime')
local fiber = require('fiber')

local fn =  {}

local s_users = box.space.user
function fn.get_group_id(user_id)
	local t_user_id = s_user.index.primary:get{user_id}
	return t_user_id[group_id]
end

function fn.user_group(user_name)
   local t_user =  box.space.user.index.name:get{user_name}
   local t_user_group =  t_user.group_id
   return box.space.group.index.group_id:get{t_user.group_id}.group_name
end


function fn.get_name(user_id)
   return box.space.user.index.primary:get{user_id}[2]

end


function fn.test_time(group_id)
   local count = 0

   local current_date = os.date("*t")

   current_date.hour = 0
   current_date.min = 0
   current_date.sec = 0

   local start_of_day = os.time(current_date)


   local all_msg = box.space.msg.index.group_id:select(group_id)
   for _, msg in pairs(all_msg) do
      if msg[5] >= start_of_day then
         count = count + 1
         if (count > 0) then
            break
         end
      end
   end

   local current_time = os.time() 
   local current_date = os.date("*t", current_time) 
   
   local formatted_date = os.time(current_date) 


   print(count)
   if (count == 0) then
      print(os.date("%d.%m.%Y", formatted_date))
      local msg_id = uuid.str()
      local system_msg_id = uuid.str()
      print(msg_id, system_msg_id)
   end
end


function fn.new_msg(message, group_id, user)
   local count = 0

   local current_date = os.date("*t")

   current_date.hour = 0
   current_date.min = 0
   current_date.sec = 0

   local start_of_day = os.time(current_date)


   local all_msg = box.space.msg.index.group_id:select(group_id)
   for _, msg in pairs(all_msg) do
      if msg[5] >= start_of_day then
         count = count + 1
         if (count > 0) then
            break
         end
      end
   end


   local msg_id = uuid.str()
   local system_msg_id = uuid.str()
   local tm = os.time()


   if (message ~= "") then
      if (count == 0) then
         local formatted_date = os.time(current_date)
         local mew_msg_date = os.date("%d.%m.%Y", formatted_date)
         box.space.msg:insert{system_msg_id, mew_msg_date, group_id, "system", tm}
      end
      tm = os.time() + 1
      box.space.msg:insert{msg_id, message, group_id, user, tm}
   end
end


function calculateCalendarDays(pastTimestamp)
   local currentDate = os.date("*t", os.time())
   local pastDate = os.date("*t", pastTimestamp)

   currentDate.hour = 0
   currentDate.min = 0
   currentDate.sec = 0
   pastDate.hour = 0
   pastDate.min = 0
   pastDate.sec = 0

   local currentDayTimestamp = os.time(currentDate)
   local pastDayTimestamp = os.time(pastDate)

   local differenceSeconds = currentDayTimestamp - pastDayTimestamp
   local daysDifference = math.floor(differenceSeconds / (24 * 3600))

   return daysDifference
end

function fn.time_test()

   local pastTimestamp = 1714950682  -- Пример timestamp
   local daysPassed = calculateCalendarDays(pastTimestamp)
   print("Календарных дней прошло: " .. daysPassed)

end

function fn.insertAll()
   mm.new_group("EscapeWorld")
   mm.new_group("Robots")
   mm.new_user("EnderAgent_X","EscapeWorld")
   mm.new_user("Bot","Robots")
end

function fn.group_msg(group_id)
   t_msg = box.space.msg.index.group_id:select{group_id}
   local t_msg_arr = {}
   for i = 1, #t_msg do 
      table.insert( t_msg_arr, {t_msg[i][2], t_msg[i][4], t_msg[i][5]} )
   end
   return t_msg_arr
end

function fn.users()
   t_users = box.space.user:select()
   local t_users_arr = {}
   for i = 1, #t_users do 
      table.insert( t_users_arr, t_users[i][2] )
   end
   return t_users_arr
end

function fn.time_group_msg(group_id, datetime) --загружает сообщения больше определенного времени
   local seven_days_ago = os.time() - (7 * 24 * 60 * 60)
   t_msg1 = box.space.msg.index.time:select({datetime}, {iterator = 'GT'})
   t_msg2 = box.space.msg.index.group_id:select{group_id}
   local combined_result = {}
   for _, v in ipairs(t_msg1) do
      for _, w in ipairs(t_msg2) do
         if v == w then
            table.insert(combined_result, v)
         end
      end
   end 
   local t_msg_arr = {}
   local cnt = 0
   if #combined_result > 99 then
      for i = #combined_result - 99, #combined_result do
         if combined_result[i][5] >= seven_days_ago then
            table.insert( t_msg_arr, {combined_result[i][2], combined_result[i][4], combined_result[i][5], combined_result[i][1]} ) --id тоже возвращяем и пихаем в name listbox
            cnt = cnt + 1
         end
      end
   else 
      for i = 1, #combined_result do 
         table.insert( t_msg_arr, {combined_result[i][2], combined_result[i][4], combined_result[i][5], combined_result[i][1]} ) --id тоже возвращяем и пихаем в name listbox
         cnt = cnt + 1
      end
   end
   return cnt, t_msg_arr
end

function fn.time_group_msg_test(group_id, datetime)
   local thirty_days_ago = os.time() - (30 * 24 * 60 * 60)
   t_msg1 = box.space.msg.index.time:select({datetime}, {iterator = 'GT'})
   t_msg2 = box.space.msg.index.group_id:select{group_id}
   local combined_result = {}
   for _, v in ipairs(t_msg1) do
      for _, w in ipairs(t_msg2) do
         if v == w then
            table.insert(combined_result, v)
         end
      end
   end 
   local t_msg_arr = {}
   local cnt = 0
   for i = 1, #combined_result do

      table.insert( t_msg_arr, {combined_result[i][2], combined_result[i][4], combined_result[i][5], combined_result[i][1]} ) --id тоже возвращяем и пихаем в name listbox
      cnt = cnt + 1
   end
   return cnt, t_msg_arr
end


function fn.new_user(name, pass) 
   local user_id = uuid.str()
   local hashed_password = digest.sha256(pass)
   if box.space.user.index.name:get(name) ~= nil then
      return false
   end
   if name == "system" then
      return false
   end
   box.space.user:insert{user_id, name, hashed_password}
   return true
end

function fn.login(name, pass)
   if box.space.user.index.name:get(name) == nil then
      return false
   end
   local user_password = box.space.user.index.name:get(name).pass 
   local hashed_password = digest.sha256(pass)
   if user_password == hashed_password then
      return true
   else 
      return false
   end
end


function fn.new_group(name, group_id, group)
   local usergroup_id = uuid.str() 
   if isStringOnlySpaces(group_id) then
      return "0"
   end
   if box.space.group.index.group_id:get(group_id) == nil then
      box.space.usergroup:insert{usergroup_id, name, group_id, "admin"}
      box.space.group:insert{group_id, group}
      return "true"
   else
      return "false"
   end
end

function fn.join_group(name, group_id)
   local usergroup_id = uuid.str()
   if box.space.group.index.group_id:get(group_id) ~= nil and box.space.usergroup.index.user_group:select({name, group_id})[1] == nil then
      box.space.usergroup:insert{usergroup_id, name, group_id, "user"}
      return box.space.group.index.group_id:get(group_id).group_name
   else
      return "false"
   end
end


function fn.hash_password(password)
   local hashed_password = digest.sha256(password)
   return hashed_password
end

function fn.get_user_groups(name)
   local t_user_groups = box.space.usergroup.index.user:select(name)
   local t_user_groups_names = {}
   for i = 1, #t_user_groups do
      local group_name = box.space.group.index.group_id:get(t_user_groups[i][3]).group_name
      table.insert( t_user_groups_names, {t_user_groups[i][3], group_name} )
   end
   return t_user_groups_names
end

function fn.del_group(user, group_id)
   local del_id = box.space.usergroup.index.user_group:select({user, group_id})[1][1]
   box.space.usergroup.index.primary:delete(del_id)
   if #box.space.usergroup.index.group_id:select(group_id) == 0 then 
      fn.del_group_msg(group_id)
      box.space.group.index.group_id:delete(group_id)
   end
end

function fn.del_group_msg(group_id)
   local messages = box.space.msg.index.group_id:select(group_id)
   for _, msg in pairs(messages) do
      local id = msg[1]
      box.space.msg.index.primary:delete(id)
   end
end

function fn.edit_group(user, group_id, new_group)

   box.space.group.index.group_id:update(group_id,{{'=', 2, new_group}})
end

function fn.edit_msg(msg_id, new_msg)
   box.space.msg.index.primary:update(msg_id,{{'=', 2, new_msg}})
   local msg = box.space.msg.index.primary:get(msg_id).message
   local user = box.space.msg.index.primary:get(msg_id).user
   local msg_time = box.space.msg.index.primary:get(msg_id).msg_time
   local group_id = box.space.msg.index.primary:get(msg_id).group_id
   local group_name = box.space.group.index.group_id:get(group_id).group_name
   return user, group_name, msg, msg_time
end

function fn.get_selected_msg(msg_id)
   local msg = box.space.msg.index.primary:get(msg_id).message
   local user = box.space.msg.index.primary:get(msg_id).user
   local group_id = box.space.msg.index.primary:get(msg_id).group_id
   local group_name = box.space.group.index.group_id:get(group_id).group_name
   return user, group_name, msg
end

function fn.del_msg(msg_id)
   box.space.msg.index.primary:delete(msg_id)
end

function fn.group_users_cnt(group_id)
   local count = #box.space.usergroup.index.group_id:select(group_id)
   return tostring(count)
end

function fn.group_users(group_id)
   local group_users = box.space.usergroup.index.group_id:select(group_id)
   local group_users_role = {}
   for i = 1, #group_users do
      local user_name = group_users[i][2]
      local user_role = group_users[i][4]
      local daysPassed = "no"
      local user_last_time = 0
      if #box.space.msg.index.user_group:select({user_name, group_id}) ~= 0 then
         local userMsgArr = box.space.msg.index.user_group:select({user_name, group_id})
         for _, message in ipairs(userMsgArr) do
            if message[5] > user_last_time then
               user_last_time = message[5]
            end
        end
        print(user_last_time)
         daysPassed = tostring(calculateCalendarDays(user_last_time))
      end

      table.insert( group_users_role, {user_name, user_role, daysPassed} )
   end
   return group_users_role
end

function fn.check_group_id(group_id)
   if isStringOnlySpaces(group_id) then
      return "0"
   end
   if box.space.group.index.group_id:get(group_id) == nil then
      return "true"
   else
      return "false"
   end
end

function isStringOnlySpaces(str)
   return str:gsub("%s+", "") == ""
end

function fn.promote_user(user, group_id, rank)
   if user ~= nil and group_id ~= nil then
      local promoteId = box.space.usergroup.index.user_group:select({user, group_id})[1][1]
      box.space.usergroup.index.primary:update(promoteId,{{'=', 4, rank}})
   end
end

function fn.get_user_role(user, group_id)
   return box.space.usergroup.index.user_group:select({user, group_id})[1][4]
end

function fn.group_exists(user, group_id)
   if box.space.usergroup.index.user_group:select({user, group_id})[1] == nil then
      return false
   else
      return true
   end
end

function fn.get_msg_cnt(group_id, days)
   local days_ago = os.time() - (days * 24 * 60 * 60)
   local count = 0
   local all_msg = box.space.msg.index.group_id:select(group_id)
   for _, msg in pairs(all_msg) do
      if msg[4] ~= 'system' and msg[5] >= days_ago then
         count = count + 1
      end
   end
   return tostring(count)
end

function fn.get_max_user_sg(group_id, days)
   local days_ago = os.time() - (days * 24 * 60 * 60)
   local user_messages = {}

   local all_msg = box.space.msg.index.group_id:select(group_id)

   for _, msg in pairs(all_msg) do
      local username = msg[4]
      if user_messages[username] and msg[4] ~= 'system' and msg[5] >= days_ago then
         user_messages[username] = user_messages[username] + 1
      else
         if msg[4] ~= 'system' and msg[5] >= days_ago then
            user_messages[username] = 1
         end
      end
   end

   local max_user = nil
   local max_count = 0
   for username, count in pairs(user_messages) do
      if count > max_count then
         max_user = username
         max_count = count
      end
   end
   return max_user, tostring(max_count)
end


return fn	
