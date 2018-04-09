package statik

import (
	"github.com/rakyll/statik/fs"
)

func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00\xe8HlL\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\n\x00\x00\x00index.html<html ng-app=\"cfgcfg\">\n\n<head>\n	<title>CfgCfg</title>\n	<link rel=\"stylesheet\" href=\"styles.css\">\n	<script src=\"https://code.jquery.com/jquery-3.3.1.min.js\"></script>\n	<script src=\"https://ajax.googleapis.com/ajax/libs/angularjs/1.6.9/angular.min.js\"></script>\n	<script src=\"script.js\"></script>\n</head>\n\n<body ng-controller=\"MainCtrl as mainctrl\" ng-keypress=\"mainctrl.keypress($event)\">\n	<div class=\"topbar\">\n		<a href=\"#\" ng-click=\"page = 'projects'\">Projects</a>\n		<a href=\"#\" ng-click=\"page = 'envs'\">Envs</a>\n	</div>\n	<div ng-show=\"page == 'projects'\" ng-controller=\"ProjectsCtrl as ctrl\">\n		<div class=\"leftbar\">\n			<input type=\"text\" ng-model=\"projectSearch\" class=\"project-search search\" autofocus placeholder=\"Search\" ng-keypress=\"ctrl.projectSearchPress($event)\"\n			/>\n			<ul>\n				<li ng-repeat=\"(k, v) in filteredProjects\" ng-class=\"{'active': currentProject == k}\">\n					<a href=\"#\" ng-click=\"ctrl.setCurrentProject(v)\">{{v.name}}</a>\n				</li>\n			</ul>\n		</div>\n		<div class=\"content\">\n			<h2>\n				<div class=\"var-search-container\">\n					<input type=\"text\" ng-model=\"varSearch\" class=\"var-search search\" placeholder=\"Search vars...\" ng-keydown=\"ctrl.varSearchPress($event)\"\n					 ng-blur=\"vsShow = false\">\n\n					<ul class=\"envs\">\n						<li ng-repeat=\"(k, v) in envs\" ng-class=\"{active: vsIndex == $index}\" ng-show=\"vsShow\">\n							{{k}}\n						</li>\n					</ul>\n				</div>\n				{{currentProject.name}}\n			</h2>\n\n			<div ng-repeat=\"(varName, var) in filteredVars\" class=\"var {{ctrl.value(varName, var) && 'ok' || 'nok'}}\">\n				<div class=\"value\">{{ctrl.value(varName, var)}}</div>\n				<div class=\"name\">{{varName}}</div>\n				<span ng-repeat=\"(k, v) in envs\">\n					<label>\n						<input ng-click=\"var.env = k; ctrl.updateProject(currentProject)\" ng-checked=\"var.env == k\" type=\"checkbox\" ng-disabled=\"!v[varName]\"> {{k}}\n					</label>\n				</span>\n				<label>\n					<input ng-checked=\"!var.env\" type=\"checkbox\" ng-click=\"var.env = ''; ctrl.updateProject(currentProject)\">\n					<input type=\"text\" ng-model=\"var.value\" placeholder=\"Custom\" ng-change=\"var.env = ''; ctrl.updateProject(currentProject)\"\n					/>\n				</label>\n\n				<select ng-show=\"var.value\" ng-change=\"ctrl.saveToEnv(saveToEnvModel, varName, var.value)\" ng-model=\"saveToEnvModel\">\n					<option value=\"\">Save to env...</option>\n					<option ng-repeat=\"(k, v) in envs\" value=\"{{k}}\">{{k}}</option>\n				</select>\n			</div>\n		</div>\n	</div>\n	<div ng-show=\"page == 'envs'\" ng-controller=\"EnvsCtrl as ctrl\">\n		<div class=\"leftbar\">\n			<ul>\n				<li ng-repeat=\"(k, v) in envs\" ng-click=\"ctrl.currentEnv = k\">\n					<a href=\"#\">{{k}}</a>\n				</li>\n			</ul>\n\n			<input placeholder=\"New\" ng-model=\"newEnv\">\n			<button ng-click=\"ctrl.saveNewEnv()\">Ok</button>\n		</div>\n\n		<div class=\"content\">\n			<h2>\n				{{ctrl.currentEnv}}\n				<input type=\"text\" ng-model=\"envVarSearch\" placeholder=\"Search\" />\n			</h2>\n\n			<div ng-repeat=\"(k, v) in envs[ctrl.currentEnv]\" class=\"var\">\n				<div class=\"value\">\n					<input type=\"text\" ng-model=\"envs[ctrl.currentEnv][k]\" ng-change=\"mainctrl.updateEnvs()\" />\n					<button ng-click=\"ctrl.delete(ctrl.currentEnv, k)\">X</button>\n				</div>\n				<div class=\"name\">{{k}}</div>\n			</div>\n			<div class=\"var\">\n				New:\n				<input type=\"text\" ng-model=\"newKey\" />\n				<input type=\"text\" ng-model=\"newValue\" />\n				<button ng-click=\"ctrl.saveNewKey()\">Ok</button>\n			</div>\n		</div>\n\n	</div>\n</body>\n\n</html>PK\x07\x08\xf1\xb9\xfb\x96:\x0d\x00\x00:\x0d\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00UjmL\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00	\x00\x00\x00script.jsfunction fuzzysearch(needle, haystack) {\n  needle = needle.toLowerCase();\n  haystack = haystack.toLowerCase();\n  var hlen = haystack.length;\n  var nlen = needle.length;\n  if (nlen > hlen) {\n    return false;\n  }\n  if (nlen === hlen) {\n    return needle === haystack;\n  }\n  outer: for (var i = 0, j = 0; i < nlen; i++) {\n    var nch = needle.charCodeAt(i);\n    while (j < hlen) {\n      if (haystack.charCodeAt(j++) === nch) {\n        continue outer;\n      }\n    }\n    return false;\n  }\n  return true;\n}\n\nangular.module('cfgcfg', [])\n    .controller(\n        'MainCtrl',\n        function($http, $scope) {\n          var ctrl = this;\n\n          $scope.page = 'projects';\n\n          ctrl.updateEnvs = function() {\n            $http.post('/envs', $scope.envs)\n                .then(\n                    function(response) {\n                      $scope.envs = response.data;\n                    },\n                    function errorCallback(response) {\n                      console.log(response)\n                    });\n          };\n\n          ctrl.keypress = function(ev) {\n            if ($('input:focus').length) {\n              return;\n            }\n            if (ev.keyCode == 115) {\n              $('.var-search').focus();\n              ev.preventDefault();\n            } else if (ev.keyCode == 112) {\n              $('.project-search').focus();\n              ev.preventDefault();\n            }\n          };\n        })\n    .controller(\n        'EnvsCtrl',\n        function($http, $scope) {\n          var ctrl = this;\n\n          $scope.$parent.envs = {};\n          $scope.currentEnv = 'local';\n          $http.get('/envs').then(\n              function(response) {\n                $scope.$parent.envs = response.data;\n                for (var i in $scope.$parent.envs) {\n                  ctrl.currentEnv = i;\n                  break;\n                }\n              },\n              function errorCallback(response) {\n                console.log(response)\n              });\n\n          ctrl.saveNewEnv = function() {\n            $scope.$parent.envs[$scope.newEnv] = {};\n            ctrl.currentEnv = $scope.newEnv;\n            $scope.newEnv = '';\n            $scope.$parent.mainctrl.updateEnvs();\n          };\n\n          ctrl.saveNewKey = function() {\n            if (!$scope.currentEnv) {\n              return\n            }\n            if (!$scope.newKey) {\n              return;\n            }\n\n            $scope.$parent.envs[$scope.currentEnv][$scope.newKey] =\n                $scope.newValue;\n            $scope.newKey = $scope.newValue = '';\n            $scope.$parent.mainctrl.updateEnvs();\n          };\n\n          ctrl.delete = function(env, k) {\n            delete ($scope.$parent.envs[env][k]);\n            $scope.$parent.mainctrl.updateEnvs();\n          };\n        })\n    .controller('ProjectsCtrl', function($http, $scope, $filter) {\n      var ctrl = this;\n\n      $scope.projects = [];\n      $scope.currentProject = {};\n      $http.get('/projects')\n          .then(\n              function(response) {\n                $scope.projects = response.data;\n                if ($scope.projects) {\n                  $scope.currentProject = $scope.projects[0];\n                }\n              },\n              function errorCallback(response) {\n                console.log(response)\n              });\n\n\n      ctrl.saveToEnv = function(env, k, v) {\n        $scope.$parent.envs[env][k] = v;\n        $scope.currentProject.vars[k].value = '';\n        $scope.currentProject.vars[k].env = env;\n        ctrl.updateProject($scope.currentProject);\n        $scope.$parent.mainctrl.updateEnvs();\n      };\n\n      ctrl.projectSearchPress = function(ev) {\n        if (ev.keyCode == 13) {\n          if ($scope.filteredProjects.length) {\n            ctrl.setCurrentProject($scope.filteredProjects[0]);\n            $scope.projectSearch = '';\n            $('.var-search')[0].focus();\n          }\n        }\n      };\n\n      ctrl.setCurrentProject = function(p) {\n        $scope.currentProject = p;\n        $scope.varSearch = '';\n      };\n\n      ctrl.value = function(k, v) {\n        if (!v) {\n          return '';\n        }\n        if (v.env) {\n          return ($scope.$parent.envs[v.env] || {})[k] || '';\n        }\n\n        return v.value || '';\n      };\n\n      ctrl.updateProject = function(project) {\n        $http.post('/projects', project)\n            .then(\n                function(response) {\n                  $scope.projects = response.data\n                },\n                function errorCallback(response) {\n                  console.log(response)\n                });\n      };\n\n      $scope.$watch('varSearch', buildFilteredVars);\n      $scope.$watch('projects', buildFilteredVars);\n      $scope.$watch('currentProject', buildFilteredVars, true);\n\n      function buildFilteredVars() {\n        if (!$scope.currentProject) {\n          return;\n        }\n        if (!$scope.varSearch) {\n          $scope.filteredVars = $scope.currentProject.vars;\n          return;\n        }\n        var vars = {};\n        for (var i in $scope.currentProject.vars) {\n          if (fuzzysearch($scope.varSearch, i)) {\n            vars[i] = $scope.currentProject.vars[i];\n          }\n        }\n        $scope.filteredVars = vars;\n      }\n\n      $scope.$watch('projectSearch', buildFilteredProjects);\n      $scope.$watch('projects', buildFilteredProjects);\n\n      ctrl.varSearchPress = function(ev) {\n        if (ev.keyCode == 40 && !$scope.vsShow) {\n          var envs = Object.keys($scope.$parent.envs);\n          if (!envs.length) {\n            return;\n          }\n          $scope.vsShow = true;\n          $scope.vsIndex = 0;\n        } else if (ev.keyCode == 40) {\n          $scope.vsIndex =\n              ($scope.vsIndex + 1) % Object.keys($scope.$parent.envs).length;\n        } else if (ev.keyCode == 38) {\n          ev.preventDefault();\n          if ($scope.vsIndex == 0) {\n            $scope.vsShow = false;\n            return;\n          }\n          $scope.vsIndex--;\n        } else if (ev.keyCode == 13 && $scope.vsShow) {\n          for (var i in $scope.filteredVars) {\n            $scope.filteredVars[i].env =\n                Object.keys($scope.$parent.envs)[$scope.vsIndex];\n          }\n          ctrl.updateProject($scope.currentProject);\n        }\n      };\n\n      function buildFilteredProjects() {\n        if (!$scope.projects) {\n          return;\n        }\n\n        var projects = [];\n        for (var i in $scope.projects) {\n          if (!$scope.projects[i].vars ||\n              !Object.keys($scope.projects[i].vars).length) {\n            continue;\n          }\n          if (!$scope.projectSearch ||\n              fuzzysearch($scope.projectSearch, $scope.projects[i].name)) {\n            projects.push($scope.projects[i]);\n          }\n        }\n        $scope.filteredProjects = projects;\n      }\n    });PK\x07\x08\xcb\x8b[\x0b\x8b\x1a\x00\x00\x8b\x1a\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xf0alL\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\n\x00\x00\x00styles.cssbody, html {\n	margin: 0;\n	font-family: \"Helvetica Neue\", Helvetica, Roboto, Arial, sans-serif;\n}\n\n* {\n	box-sizing: border-box;\n}\n\n.topbar {\n	height: 50px;\n	background: #a9b7ce;\n}\n\n.topbar a {\n	display: inline-block;\n	padding: 0 20px;\n	height: 50px;\n	line-height: 50px;\n	color: #242b36;\n}\n\n.leftbar {\n	float: left;\n	width: 300px;\n	background: #eee;\n	height: calc(100% - 50px);\n}\n\n.leftbar ul {\n	margin: 0;\n	padding: 20px;\n	list-style-type: none;\n	max-height: calc(100% - 70px);\n	overflow: auto;\n}\n\n.leftbar ul li {\n	padding: 5px 0;\n}\n\n.leftbar ul li a {\n	color: #475b83;\n}\n\n.leftbar ul li.active {\n	background: #acc8ff;\n}\n\n.leftbar .search {\n	display: block;\n	width: 300px;\n	height: 30px;\n	margin-top: 20px;\n	margin-bottom: 20px;\n	padding-left: 20px;\n	padding-right: 20px;\n	border-left: 0;\n	border-right: 0;\n	border-top: 1px solid #aaa;\n	border-bottom: 1px solid #aaa;\n}\n\n.content {\n	height: calc(100% - 50px);\n	overflow: auto;\n	padding: 20px;\n}\n\nh2 {\n	margin: 0;\n	padding: 0 0 15px 0;\n}\n\nh2 .search {\n	float: right;\n	height: 30px;\n}\n\np {\n	margin: 0;\n}\n\n.var {\n	padding: 10px;\n	margin: 5px 0;\n	background: #d3dcff;\n}\n\n.var.ok {\n	background: #e5ffd3;\n}\n\n.var.nok {\n	background: #ffd3d3;\n}\n\n.var .value, .var .name {\n	font-family: monospace;\n	font-size: 1.2em;\n	padding-bottom: 5px;\n}\n\n.var .name {\n	font-weight: bold;\n}\n\n.var .value {\n	float: right;\n}\n\n.var-search-container {\n	float: right;\n	position: relative;\n	width: 200px;\n}\n\n.var-search-container input {\n	width: 100%;\n}\n\n.var-search-container .envs {\n	background: white;\n	margin: 0;\n	padding: 0;\n	list-style-type: none;\n	position: absolute;\n	top: 30px;\n	width: 100%;\n	border-bottom: 1px solid rgb(238, 238, 238);\n	border-left: 1px solid rgb(238, 238, 238);\n	border-right: 1px solid rgb(238, 238, 238);\n}\n\n.var-search-container .envs li {\n	padding: 5px;\n	font-size: 15px;\n}\n\n.var-search-container .envs li.active {\n	background: #d3dcff;\n}PK\x07\x08\x06\xd5\xc7\xc0d\x07\x00\x00d\x07\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xe8HlL\xf1\xb9\xfb\x96:\x0d\x00\x00:\x0d\x00\x00\n\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00index.htmlPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00UjmL\xcb\x8b[\x0b\x8b\x1a\x00\x00\x8b\x1a\x00\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81r\x0d\x00\x00script.jsPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xf0alL\x06\xd5\xc7\xc0d\x07\x00\x00d\x07\x00\x00\n\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x814(\x00\x00styles.cssPK\x05\x06\x00\x00\x00\x00\x03\x00\x03\x00\xa7\x00\x00\x00\xd0/\x00\x00\x00\x00"
	fs.Register(data)
}
