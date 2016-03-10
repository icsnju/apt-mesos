'use strict';

// Declare app level module which depends on views, and components
angular.module('sher', [
  'ui.router',
  'sher.task',
]).
config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
  $urlRouterProvider.otherwise('/task');

  $stateProvider
      .state("task", {
        url: "/task?query",
        templateUrl: "/js/templates/task.html",
        controller: 'TaskCtrl'
      });
}]);
