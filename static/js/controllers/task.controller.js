/**
 * Created by jeff on 16/3/9.
 */

var FILTER_NAME = ['All', 'Running', 'Finished', 'Failed', 'Killed', 'Lost', 'Staging', 'Error'];

angular.module('sher.task', ['ngResource', 'ui.bootstrap'])

.controller('TaskCtrl', [
    '$scope',
    '$http',
    '$timeout',
    '$state',
    '$stateParams',
    '$uibModal',
    'Tasks',
function($scope, $http, $timeout, $state, $stateParams, $uibModal, Tasks) {
    $scope.query = $stateParams.query || "all";
    $scope.filter = $scope.query

    // 加载数据
    var reload = function (query) {
        Tasks.refresh().$promise.then(function(response) {
            //TODO 错误处理
            $scope.tasks = Tasks.getTasks(query);
            console.log($scope.tasks);
        });
    }

    // 提交任务
    $scope.submitTask = function (task) {
        Tasks.submitTask(task, reload($scope.query))
    }

    // 杀死任务
    $scope.kill = function (task) {
        Tasks.killTask(task.id, reload($scope.query));
    }

    // 删除任务
    $scope.delete = function (task) {
        Tasks.deleteTask(task.id, reload($scope.query));
    }

    // 搜索任务
    $scope.search = function () {
        $state.go('task', {query: $scope.search_key})
    }

    // 打开提交任务的模态框
    $scope.openTaskModal = function () {
        var modalInstance = $uibModal.open({
            animation: true,
            templateUrl: '/js/templates/task.modal.html',
            controller: TaskModalCtrl,
            size: 'md',
            windowTemplateUrl: '/js/components/modal/modal.window.html',
            resolve: {
            }
        });
    }

    // 打开日志的模态对话框
    $scope.openLogModal = function (task) {
        var modalInstance = $uibModal.open({
            animation: true,
            templateUrl: '/js/templates/output.modal.html',
            controller: LogModalCtrl,
            size: 'md',
            windowTemplateUrl: '/js/components/modal/console.window.html',
            resolve: {
                "task": task,
            }
        })
    }

    // 加载任务, 定时监控
    reload($scope.query);
    setInterval(function(){
        Tasks.monitor(reload($scope.query))
    },1000)
}]);


// 模块对话框控制器
var TaskModalCtrl = function ($scope, $uibModalInstance, Tasks) {
    // 数据初始化
    $scope.task = {
        cpus:'0.5',
        mem:'64',
        disk:'0',
        docker_image:'busybox',
        cmd:'ls',
        volumes: [
            {
                container_path: "/data",
                host_path: "/vagrant",
                mode: "RW"
            }
        ],
        port_mappings: [
            {
                container_port: "8080",
                host_port: "31020",
                protocol: "tcp"
            }
        ]
    }

    $scope.addPortMapping = function() {
        $scope.task.port_mappings.push({
            container_port: "8080",
                host_port: "31020",
            protocol: "tcp"
        })
    }

    $scope.deletePortMapping = function(index) {
        $scope.task.port_mappings.splice(index, 1);
    }

    $scope.addVolume = function() {
        $scope.task.volumes.push({
            container_path: "/data",
            host_path: "/vagrant",
            mode: "RW"
        })
    }

    $scope.deleteVolume = function(index) {
        $scope.task.volumes.splice(index, 1);
    }

    $scope.submit = function () {
        Tasks.submitTask($scope.task, function(){
            // TODO 消息通知
        });
        $uibModalInstance.close();
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
};

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
