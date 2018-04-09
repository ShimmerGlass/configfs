function fuzzysearch(needle, haystack) {
  needle = needle.toLowerCase();
  haystack = haystack.toLowerCase();
  var hlen = haystack.length;
  var nlen = needle.length;
  if (nlen > hlen) {
    return false;
  }
  if (nlen === hlen) {
    return needle === haystack;
  }
  outer: for (var i = 0, j = 0; i < nlen; i++) {
    var nch = needle.charCodeAt(i);
    while (j < hlen) {
      if (haystack.charCodeAt(j++) === nch) {
        continue outer;
      }
    }
    return false;
  }
  return true;
}

angular.module('cfgcfg', [])
    .controller(
        'MainCtrl',
        function($http, $scope) {
          var ctrl = this;

          $scope.page = 'projects';

          ctrl.updateEnvs = function() {
            $http.post('/envs', $scope.envs)
                .then(
                    function(response) {
                      $scope.envs = response.data;
                    },
                    function errorCallback(response) {
                      console.log(response)
                    });
          };

          ctrl.keypress = function(ev) {
            if ($('input:focus').length) {
              return;
            }
            if (ev.keyCode == 115) {
              $('.var-search').focus();
              ev.preventDefault();
            } else if (ev.keyCode == 112) {
              $('.project-search').focus();
              ev.preventDefault();
            }
          };
        })
    .controller(
        'EnvsCtrl',
        function($http, $scope) {
          var ctrl = this;

          $scope.$parent.envs = {};
          $scope.currentEnv = 'local';
          $http.get('/envs').then(
              function(response) {
                $scope.$parent.envs = response.data;
                for (var i in $scope.$parent.envs) {
                  ctrl.currentEnv = i;
                  break;
                }
              },
              function errorCallback(response) {
                console.log(response)
              });

          ctrl.saveNewEnv = function() {
            $scope.$parent.envs[$scope.newEnv] = {};
            ctrl.currentEnv = $scope.newEnv;
            $scope.newEnv = '';
            $scope.$parent.mainctrl.updateEnvs();
          };

          ctrl.saveNewKey = function() {
            if (!$scope.currentEnv) {
              return
            }
            if (!$scope.newKey) {
              return;
            }

            $scope.$parent.envs[$scope.currentEnv][$scope.newKey] =
                $scope.newValue;
            $scope.newKey = $scope.newValue = '';
            $scope.$parent.mainctrl.updateEnvs();
          };

          ctrl.delete = function(env, k) {
            delete ($scope.$parent.envs[env][k]);
            $scope.$parent.mainctrl.updateEnvs();
          };
        })
    .controller('ProjectsCtrl', function($http, $scope, $filter) {
      var ctrl = this;

      $scope.projects = [];
      $scope.currentProject = {};
      $http.get('/projects')
          .then(
              function(response) {
                $scope.projects = response.data;
                if ($scope.projects) {
                  $scope.currentProject = $scope.projects[0];
                }
              },
              function errorCallback(response) {
                console.log(response)
              });


      ctrl.saveToEnv = function(env, k, v) {
        $scope.$parent.envs[env][k] = v;
        $scope.currentProject.vars[k].value = '';
        $scope.currentProject.vars[k].env = env;
        ctrl.updateProject($scope.currentProject);
        $scope.$parent.mainctrl.updateEnvs();
      };

      ctrl.projectSearchPress = function(ev) {
        if (ev.keyCode == 13) {
          if ($scope.filteredProjects.length) {
            ctrl.setCurrentProject($scope.filteredProjects[0]);
            $scope.projectSearch = '';
            $('.var-search')[0].focus();
          }
        }
      };

      ctrl.setCurrentProject = function(p) {
        $scope.currentProject = p;
        $scope.varSearch = '';
      };

      ctrl.value = function(k, v) {
        if (!v) {
          return '';
        }
        if (v.env) {
          return ($scope.$parent.envs[v.env] || {})[k] || '';
        }

        return v.value || '';
      };

      ctrl.updateProject = function(project) {
        $http.post('/projects', project)
            .then(
                function(response) {
                  $scope.projects = response.data
                },
                function errorCallback(response) {
                  console.log(response)
                });
      };

      $scope.$watch('varSearch', buildFilteredVars);
      $scope.$watch('projects', buildFilteredVars);
      $scope.$watch('currentProject', buildFilteredVars, true);

      function buildFilteredVars() {
        if (!$scope.currentProject) {
          return;
        }
        if (!$scope.varSearch) {
          $scope.filteredVars = $scope.currentProject.vars;
          return;
        }
        var vars = {};
        for (var i in $scope.currentProject.vars) {
          if (fuzzysearch($scope.varSearch, i)) {
            vars[i] = $scope.currentProject.vars[i];
          }
        }
        $scope.filteredVars = vars;
      }

      $scope.$watch('projectSearch', buildFilteredProjects);
      $scope.$watch('projects', buildFilteredProjects);

      ctrl.varSearchPress = function(ev) {
        if (ev.keyCode == 40 && !$scope.vsShow) {
          var envs = Object.keys($scope.$parent.envs);
          if (!envs.length) {
            return;
          }
          $scope.vsShow = true;
          $scope.vsIndex = 0;
        } else if (ev.keyCode == 40) {
          $scope.vsIndex =
              ($scope.vsIndex + 1) % Object.keys($scope.$parent.envs).length;
        } else if (ev.keyCode == 38) {
          ev.preventDefault();
          if ($scope.vsIndex == 0) {
            $scope.vsShow = false;
            return;
          }
          $scope.vsIndex--;
        } else if (ev.keyCode == 13 && $scope.vsShow) {
          for (var i in $scope.filteredVars) {
            $scope.filteredVars[i].env =
                Object.keys($scope.$parent.envs)[$scope.vsIndex];
          }
          ctrl.updateProject($scope.currentProject);
        }
      };

      function buildFilteredProjects() {
        if (!$scope.projects) {
          return;
        }

        var projects = [];
        for (var i in $scope.projects) {
          if (!$scope.projects[i].vars ||
              !Object.keys($scope.projects[i].vars).length) {
            continue;
          }
          if (!$scope.projectSearch ||
              fuzzysearch($scope.projectSearch, $scope.projects[i].name)) {
            projects.push($scope.projects[i]);
          }
        }
        $scope.filteredProjects = projects;
      }
    });