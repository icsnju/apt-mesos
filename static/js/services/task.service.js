var API = 'http://192.168.33.1:3030/api';

angular.module('sher.task')
    .factory('Tasks', ['$resource', '$http', function($resource, $http) {
        var tasks = [];
        var resource = $resource(API + '/tasks', {}, {
            query: {
                method: 'get',
                timeout: 20000
            },
        })

        var getTasks = function(callback) {
            return resource.query({

            }, function(r) {
                return callback && callback(r);
            })
        };

        return {
            // 刷新任务
            refresh: function() {
                return getTasks(function(response) {
                    tasks = handleTasks(response.message);
                })
            },

            // 重置数据
            resetData: function() {
                tasks = [];
            },

            // 获取全部的任务
            getAllTasks: function() {
                return tasks;
            },

            // 搜索任务
            getTasks: function(key) {
                if(key == 'all') {
                    return tasks;
                } else {
                    var result = [];
                    var pattern = new RegExp(key,'ig');
                    for (var i = 0; i < tasks.length; i++) {
                        if(JSON.stringify(tasks[i]).match(pattern)) {
                            result.push(tasks[i]);
                        }
                    }
                    return result;
                }
            },

            // 按ID获取任务
            getById: function(id) {
                if (!!tasks) {
                    for (var i = 0; i < tasks.length; i++) {
                        if (tasks[i].id === id) {
                            return tasks[i];
                        }
                    }
                } else {
                    return null;
                }
            },

            // 提交任务
            submitTask: function(task, callback) {
                $http({
                    method: 'POST',
                    url: API + '/tasks',
                    data : task,
                    headers:{
                        'Accept': 'application/json',
                        'Content-Type': 'application/json; ; charset=UTF-8'
                    }
                }).success(function(response) {
                    return callback;
                });
            },

            // 删除任务
            deleteTask: function(id, callback) {
                $http({
                    method: 'DELETE',
                    url: API + '/tasks/' + id
                }).success(function(response) {
                    return callback && callback(response);
                })
            },

            // 杀死任务
            killTask: function(id, callback) {
                $http({
                    method: 'PUT',
                    url: API + '/tasks/' + id + '/kill'
                }).success(function(response) {
                    return callback && callback(response);
                })
            },

            // 读取任务输出
            getTaskFile: function(id, file, callback) {
                $http({
                    method: 'GET',
                    url: API + '/tasks/' + id + '/file/' + file
                }).success(function(response) {
                    return callback && callback(response);
                })
            },

            systemUsage: function(callback) {
                $http({
                    method: 'GET',
                    url: API + '/system/usage'
                }).success(function(response) {
                    return callback && callback(response);
                })
            }   

        }
    }]);

function handleTasks(tasks) {
    for(var i = 0; i < tasks.length; i++) {
        // 转换状态
        switch (tasks[i].state) {
            case "TASK_STAGING":
                tasks[i].status="STARTING";
                tasks[i].label_class="primary";
                break;
            case "TASK_RUNNING":
                tasks[i].status="RUNNING";
                tasks[i].label_class="info";
                break;
            case "TASK_FINISHED":
                tasks[i].status="FINISHED";
                tasks[i].label_class="success";
                break;
            case "TASK_FAILED":
                tasks[i].status="FAILED";
                tasks[i].label_class="danger";
                break;
            case "TASK_KILLED":
                tasks[i].status="KILLED";
                tasks[i].label_class="warning";
                break;
            case "TASK_LOST":
                tasks[i].status="LOST";
                tasks[i].label_class="default";
                break;
        }
    }

    tasks.sort(function(a, b) {
        return b.create_time - a.create_time;
    })

    return tasks
} 