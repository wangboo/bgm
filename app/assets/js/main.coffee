app = angular.module("myapp", ['ngRoute']) 

# 路由配置
app.config ["$routeProvider", ($routeProvider)->
	$routeProvider.when("/", 
		templateUrl: 	"/tpl/main.html"
		controller: 	"MainCtrl")
	.when("/platform/:id",
		templateUrl: 	"/tpl/platform.html"
		controller: 	"PlatformCtrl")
	.when("/user/:sid/:uid",
		templateUrl: 	"/tpl/user.html"
		controller: 	"UserCtrl")
	.when("/user_charge/:sid/:uid",
		templateUrl: 	"/tpl/user_charge.html"
		controller: 	"UserChargeCtrl")
	.when("/user_setting/:sid/:uid",
		templateUrl: 	"/tpl/user_setting.html"
		controller: 	"UserSettingCtrl")
	.when("/ac",
		templateUrl:	"/tpl/ac.html"
		controller:		"AcCtrl")
	.when("/login",
		templateUrl: 	"/tpl/login.html"
		controller: 	"AuthCtrl")
]

app.controller "AuthCtrl", ["$scope", "$http", "$location", ($scope, $http, $location)->
	$scope.loginBtn = loginBtn($scope, $http, $location)
]

loginBtn = ($scope, $http, $location)-> ()->
	data = 
		username: $scope.username
		password: $scope.password
	console.log "login params = ", data
	$.post("/json/auth/login", data, (data)->
		console.log "login resp ", data
		if data.ok 
			$location.path("/")
		else 
			$scope.error = data.msg 
			# $scope.$apply()
	)
	# $http.post("/json/auth/login", data).success (data)->
	# 	console.log "login resp ", data
	# 	if data.ok 
	# 		$location.path("/")
	# 	else 
	# 		$scope.error = data.msg 
	# 		$scope.$apply()

# 主页面
app.controller "MainCtrl", ["$scope", "$http", "$location", ($scope, $http, $location)->
	$http.get("/json/auth/validate").success (data)->
		console.log "auth = ", data
		unless data.ok 
			$location.path("/login")


	$http.get("/json/platforms").success((data)-> $scope.platforms = data)
	$scope.btnClick = (index)->
		p = $scope.platforms[index]
		console.log "p = ", p
		$location.path("/platform/#{p.Id}")
]

# 平台
app.controller "PlatformCtrl", ["$scope", "$http", "$location", "$window", ($scope, $http, $location, $window)->
	path = $location.path()
	$scope.queryMode = "name"
	$scope.timeLongToString = timeLongToString
	$scope.jumpToUserInfo = jumpToUserInfo($scope, $location)
	$scope.jumpToPlatform = jumpToPlatform($scope, $http, $location)
	$scope.pid = path.replace("/platform/", "")
	$http.get("/json#{path}").success((data)-> 
		console.log "#{path} resp:", data
		console.log "data.platform.Servers = ", data.platform.Servers
		setServerState(data.platform.Servers, data.ss)
		setWorkState(data.platform.Servers)
		$scope.platform = data.platform
		$scope.ss = data.ss 
		$scope.menu = data.menu or 'findUser'
		console.log "resp", data
	)
	# 模糊查询玩家
	$scope.findUser = ()->
		unless name = $scope.findUserName then return 
		queryMode = $scope.queryMode or "name"
		if queryMode == 'id' and not /^\d+$/.test(name) then return $window.alert("id搜索模式必须输入整数")
		serverIds = for s in $scope.platform.Servers when s.check then s.Id
		if serverIds.length == 0 then return $window.alert("请勾选查询服务器")
		console.log serverIds
		data = {name: name, ids: serverIds.join(","), queryMode: queryMode}
		$.post("/json/gs/findUserByName", data, (resp)->
			$scope.findUsers = for sid in serverIds 
				sinfo = findServerById($scope.platform.Servers, sid)
				data = if resp[sid] instanceof Array then resp[sid] else [resp[sid]]
				{sinfo: sinfo, data: data}
			console.log "$scope.findUsers = ", $scope.findUsers
			$scope.$apply()
		)
]

# 跳转到用户详情
jumpToUserInfo = ($scope, $location)-> (pindex, index)->
	if arguments.length == 0 
		sid = $scope.sid
		uid = $scope.uid 
	else 
		sid = $scope.findUsers[pindex].sinfo.Id
		uid = $scope.findUsers[pindex].data[index].userId
	$location.path("/user/#{sid}/#{uid}")

# 跳转到用户充值
jumpToUserCharge = ($scope, $location)-> ()->
	sid = $scope.sid
	uid = $scope.uid 
	$location.path("/user_charge/#{sid}/#{uid}")

# 跳转到当前平台
jumpToPlatform = ($scope, $http ,$location)-> ()->
	if $scope.sid 
		$http.get("/json/pid?sid=#{$scope.sid}").success (data)->
			$location.path("/platform/#{data.pid}")
	else if $scope.pid 
		$location.path("/platform/#{$scope.pid}")

# 跳转到用户设置
jumpToUserSetting = ($scope, $location)-> ()->
	$location.path("/user_setting/#{$scope.sid}/#{$scope.uid}")
# 插件列表
directives = ['finduser', 'user_operation_menu', 'user_query_menu']
# 批量注册标签
directives.forEach (name)-> 
	tagName = name.split("_").join("")
	app.directive tagName, ()->
		templateUrl: "/tpl/directive/#{name}.html"
		restrict: "AE"
		replace: true

setServerState = (servers, ss)->
	for s in servers
		for st in ss when st.Id is s.ServerStateId then s.State = st.Name

# 0 可见&可用，1可见&不可用，2不可见&不可用
WorkStateMap = {0: "正常", 1: "维护", 2: "停服"}

setWorkState = (servers) ->
	for s in servers then s.WorkStateName = WorkStateMap[s.WorkState]

findServerById = (servers, id)-> for s in servers when s.Id == id then return s ;return {}

timeLongToString = (timeLong)-> 
	d = new Date(timeLong)
	month = ti2s(d.getMonth()+1)
	date = ti2s(d.getDate())
	hour = ti2s(d.getHours())
	minute = ti2s(d.getMinutes())
	"#{d.getFullYear()}-#{month}-#{date} #{hour}:#{minute}"

ti2s = (num) -> if num < 10 then "0#{num}" else num


