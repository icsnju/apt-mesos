'use strict';

var detail = angular.module('sher.detail',['ngMaterial', 'ngMessages', 'material.svgAssetsCache', 'chart.js', 'ui.router']);

detail.controller("mesCtrl", ['$scope', '$http', '$stateParams', 'Tasks','$uibModal', "$state",
	function($scope, $http, $stateParams, Tasks, $uibModal, $state){
		Tasks.refresh();
		$scope.task = Tasks.getById($stateParams.taskID);

	    // 打开日志的模态对话框
	    $scope.openLogModal = function (task) {
	        var modalInstance = $uibModal.open({
	            animation: true,
	            templateUrl: '/app/js/templates/output.modal.html',
	            controller: LogModalCtrl,
	            size: 'md',
	            windowTemplateUrl: '/app/js/components/modal/console.window.html',
	            resolve: {
	                "task": task,
	            }
	        })
	    }


	    // 杀死任务
	    $scope.kill = function (task) {
	        Tasks.killTask(task.id, function() {
	        	$state.go('navbar.tasks');
	        });
	    }

	    // 删除任务
	    $scope.delete = function (task) {
	        Tasks.deleteTask(task.id, function() {
	        	$state.go('navbar.tasks');
	        });
	    }	    
	}
]);

detail.controller("cpuCtrl", function ($scope, $http) {

	$scope.labels = ["January", "February", "March", "April", "May", "June", "July"];
	$scope.series = ['Series A', 'Series B'];
	$scope.data = [
		[28, 48, 40, 19, 86, 27, 90],
		[30, 49, 40, 9, 86, 27, 90],
	];
	setInterval(function(){
		$http.get('data/status.json').success(function(data) {
			$scope.series = ['Series A', 'Series B'];
            $scope.data = [
                [65, 59, 80, 81, 56, 55, 40],
                [28, 48, 40, 19, 86, 27, 90]
			];
		});
	},10000)
});

detail.controller("memCtrl", function ($scope, $http) {

	$scope.labels = ["January", "February", "March", "April", "May", "June", "July"];
	$scope.series = ['Series A', 'Series B'];
	$scope.data = [
		[28, 42, 41, 19, 86, 27, 90],
		[30, 56, 40, 14, 80, 23, 91],
	];
	setInterval(function(){
		$http.get('data/status.json').success(function(data) {
			$scope.series = ['Series A', 'Series B'];
            $scope.data = [
                [65, 59, 80, 81, 56, 55, 40],
                [22, 44, 49, 12, 81, 22, 94]
			];
		});
	},10000)
});

// 日志模块对话框控制器
var LogModalCtrl = function ($scope, $uibModalInstance, Tasks, task) {
    // 默认在logtab下
    $scope.logTab = 'stderr';
    consoleLog('stderr')

    $scope.console = function(file) {
        consoleLog(file)
    };

    $scope.refresh = function(file) {
        consoleLog(file)
    }

    function consoleLog(file) {
        $scope.logTab = file;
        Tasks.getTaskFile(task.id, file, function(response){
            $scope.logs = response.message.split('\n');
        });
    }
};
