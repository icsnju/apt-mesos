'use strict';

var overview = angular.module("sher.overview", ["chart.js"]);

overview.controller("tableCtrl", ['$scope', 'Tasks', function ($scope, Tasks) {
	Tasks.refresh();
	$scope.tasks = Tasks.getTasks().slice(0, 5);
}
]);

overview.controller("pieCtrl", ['$scope', 'Tasks', function ($scope, Tasks) {
	//TODO 增加service
	Tasks.refresh();
	Tasks.systemUsage(function(response) {
		var metrics = response.message;
	  	$scope.labels = ["free", "used"];
	  	$scope.cpus = [metrics.free_cpus, metrics.used_cpus];
	  	$scope.mem = [metrics.free_mem, metrics.used_mem];
	  	$scope.disk = [metrics.free_disk, metrics.used_disk];	
	  	console.log($scope.mem)	
	});
}
]);

overview.controller("pmemCtrl", function ($scope) {
  $scope.labels = ["Download Sales", "In-Store Sales", "Mail-Order Sales"];
  $scope.data = [200, 500, 100];
});

overview.controller("pnetCtrl", function ($scope) {
  $scope.labels = ["Download Sales", "In-Store Sales", "Mail-Order Sales"];
  $scope.data = [300, 500, 400];
});
