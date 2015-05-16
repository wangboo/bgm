app.controller "AcCtrl", ["$scope", "$http", "$location", "$window", ($scope, $http, $location, $window)->
	$http.get("/json/platforms").success (data)-> $scope.platforms = data
	$http.get("/json/ac/all_type").success (data)-> $scope.all_type = data
	$http.get("/json/ac/all_reward").success (data)-> $scope.all_reward = data
	$scope.beginAt 	= "2015-01-01"
	$scope.endAt 		= "2016-01-01"
	$scope.acCommit = acCommit($scope, $http, $window)
	$scope.queryServer = queryServer($scope)
	$scope.all = true 
	$scope.size = 1000
	$scope.allServer = true 
]

queryServer = ($scope)-> ()->
	pids = for p in $scope.platforms when p.select then p.Id 
	console.log "queryServer ", pids
	$.get("/json/platformServers/#{pids.join(",")}", (data)-> $scope.selectServers = data)


acCommit = ($scope, $http, $window)->()->
	all = $scope.all
	allServer = $scope.allServer
	isMuti = if $scope.isMuti then true else false 
	mutiTimes = if isMuti then $scope.mutiTimes else 0 
	size = parseInt($scope.size)
	if size.toString() == 'NaN' or size <= 0
		$window.alert("数量必须为正整数")
		return 
	if not /^\d\d\d\d-\d\d-\d\d$/.test($scope.beginAt) 
		$window.alert("开始时间格式必须为2015-01-02样式")
		return 
	if not /^\d\d\d\d-\d\d-\d\d$/.test($scope.endAt)
		$window.alert("结束时间格式必须为2015-01-02样式")
		return 
	data = 
		all: 				all
		allServer: 	allServer
		name: 			$scope.name
		size: 			size 
		rType: 			$scope.type
		reward: 		$scope.reward
		beginAt: 		$scope.beginAt
		endAt:			$scope.endAt
		isMuti: 		isMuti
		prefix: 		$scope.prefix
		mutiTimes: 	mutiTimes
		desc: 			$scope.desc
	for k, v of data when v == undefined 
		$window.alert("#{k} 没有填写")
		return 
	data.pids = if all 
		[]
	else
		for p in $scope.platforms when p.select then p.Id 
	data.pids = data.pids.join(",")
	data.sids = if allServer 
		[]
	else
		for s in $scope.selectServers when s.select then s.Id 
	console.log "data.sids.length = ", data.sids.length
	if not allServer and data.sids.length == 0 
		$window.alert("请至少选择一个平台")
		return 
	data.sids = data.sids.join(",")
	console.log "create = ", data
	$http.get("/json/ac/create", {params: data}).success (data)-> 
		console.log "resp=", data
		if data.ok 
			$window.alert("创建成功")
			$scope.all = true 
			$scope.allServer = true 
			$scope.prefix = ""
			$scope.name = ""
			$scope.type = ""
			$scope.reward = ""
			$scope.desc = ""
			$scope.isMuti = false
			$scope.size = 1000
			data.pids = ""
			data.sids = ""
			for s in $scope.selectServers then s.select = false 
			for p in $scope.platforms then p.select = false 
		else
			$window.alert(data.msg)
