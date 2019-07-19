angular.module('portainer.app')
  .controller('BuildController', ['$q', '$scope',
    function($q, $scope) {
        function initView() {
            console.log("test build")
        }

        initView();
    }]);
