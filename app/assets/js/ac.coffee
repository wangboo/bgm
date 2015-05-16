app.controller "AcCtrl", ["$scope", "$http", "$location", "$window", ($scope, $http, $location, $window)->
	$http.get("/json/platforms").success (data)-> $scope.platforms = data
	$http.get("/json/ac/all_type").success (data)-> $scope.all_type = data
	$http.get("/json/ac/all_reward").success (data)-> $scope.all_reward = data
	$scope.beginAt 	= "2015-01-01"
	$scope.EndAt 		= "2016-01-01"
	$scope.acCommit = acCommit($scope, $http, $window)
]

acCommit = ($scope, $http, $window)->()->
	all = if $scope.all then true else false
	isMuti = if $scope.isMuti then true else false 
	mutiTimes = if isMuti then $scope.mutiTimes else 0 
	size = parseInt($scope.size)
	if size.toString() == 'NaN' or size <= 0
		$window.alert("数量必须为正整数")
		return 
	if not /^\d\d\d\d-\d\d-\d\d$/.test($scope.beginAt) 
		$window.alert("开始时间格式必须为2015-01-02样式")
		return 
	if not /^\d\d\d\d-\d\d-\d\d$/.test($scope.EndAt) 
		$window.alert("结束时间格式必须为2015-01-02样式")
		return 
	data = 
		all: 				all
		name: 			$scope.name
		size: 			size 
		rType: 			$scope.type
		reward: 		$scope.reward
		beginAt: 		$scope.beginAt
		endAt:			$scope.endAt
		isMuti: 		isMuti
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
	console.log "create = ", data
	$http.get("/json/ac/create", {params: data}).success (data)-> console.log "resp=", data
