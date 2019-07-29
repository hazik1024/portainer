import { jsonObjectsToArrayHandler, genericHandler } from './response/handlers';

angular.module('portainer.docker')
.factory('Build', ['$resource', 'API_ENDPOINT_ENDPOINTS', 'EndpointProvider', function BuildFactory($resource, API_ENDPOINT_ENDPOINTS, EndpointProvider) {
  'use strict';
  return $resource(API_ENDPOINT_ENDPOINTS + '/:endpointId/docker/build', {
    endpointId: EndpointProvider.endpointID
  },
  {
    buildImage: {
      method: 'POST', ignoreLoadingBar: true,
      transformResponse: jsonObjectsToArrayHandler, isArray: true,
      headers: { 'Content-Type': 'application/x-tar' }
    },
    buildImageOverride: {
      method: 'POST', ignoreLoadingBar: true,
      transformResponse: jsonObjectsToArrayHandler, isArray: true
    },
  });
}]);


angular.module('portainer.docker')
.factory('TestBuild', ['$resource', 'API_CUSTOM', 'EndpointProvider', function BuildFactory($resource, API_CUSTOM) {
  'use strict';
  return $resource(API_CUSTOM + '/build', {
  },
  {
    testBuild: {
      method: 'POST', ignoreLoadingBar: true,
      transformResponse: genericHandler, isArray: false
    }
  });
}]);

angular.module('portainer.docker')
.factory('TestStack', ['$resource', 'API_CUSTOM', 'EndpointProvider', function BuildFactory($resource, API_CUSTOM) {
  'use strict';
  return $resource(API_CUSTOM + '/stack', {
  },
  {
    testStack: {
      method: 'POST', ignoreLoadingBar: true,
      transformResponse: genericHandler, isArray: false
    }
  });
}]);

