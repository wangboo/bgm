# 玩家Ctrl
app.controller "UserCtrl", ["$scope", "$http", "$location", ($scope, $http, $location)->
	initUserCtrl($scope, $location, $http)
	$scope.menu = 'info'
	sid = $scope.sid
	uid = $scope.uid
	unless $scope.user then getUserInfo($http, $scope, $scope.sid, $scope.uid)
	$scope.getQuanlityName = getQuanlityName
	$scope.qualityClass = qualityClass
	$scope.propClass = propClass
	$http.get("/json/gs/findGroup?sid=#{sid}&uid=#{uid}").success((data)->
		# console.log "resp = ", data
		$scope.heros = data.heros or []
	)
	$http.get("/json/gs/findUserProp?sid=#{sid}&uid=#{uid}").success((data)->
		console.log "findUserProp = ", data
		for i in data.list then i.propClass = "prop_#{i.q}"
		$scope.props = data.list or []
	)
]

getQuanlityName = (id)-> ["金色", "紫色", "蓝色"][id]

propClass = (i, q)-> if i % 3 != 0 then "q#{q} col-md-offset-1" else "q#{q}"

qualityClass = (q)-> "q#{q}"

getUserInfo = ($http, $scope, sid, uid)->
	$http.get("/json/gs/info?sid=#{sid}&uid=#{uid}").success((data)->
		console.log "getUserInfo = ", data
		$scope.user = data)

initUserCtrl = ($scope, $location, $http)->
	uri = $location.path()
	console.log "uri = #{uri}"
	$scope.timeLongToString = timeLongToString
	[_,_, sid, uid] = uri.split("/")
	$scope.sid = sid 
	$scope.uid = uid
	$scope.jumpToUserInfo 		= jumpToUserInfo($scope, $location)
	$scope.jumpToUserCharge 	= jumpToUserCharge($scope, $location)
	$scope.jumpToPlatform 		= jumpToPlatform($scope, $http ,$location)
	$scope.jumpToUserSetting 	= jumpToUserSetting($scope, $location)
