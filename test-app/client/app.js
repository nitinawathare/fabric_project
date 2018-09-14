// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory)	{

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllStampPaper = function()	{

		appFactory.queryAllStampPaper(function(data)	{
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_sp = array;
		});
	}

	$scope.queryStampPaper = function(){

		var id = $scope.stamp_id;

		appFactory.queryStampPaper(id, function(data){
			$scope.query_sp = data;

			if ($scope.queryStampPaper == "Stamp Paper not found"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordStamp = function()	{

		appFactory.recordStamp($scope.stamp, function(data){
			$scope.create_stamp = data;
			$("#success_create").show();
		});
	}

	$scope.changeHolder = function(){

		appFactory.changeHolder($scope.holder, function(data){
			$scope.change_holder = data;
			if ($scope.change_holder == "Error: no tuna catch found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllStampPaper = function(callback)	{

    	$http.get('/get_all_sp/').success(function(output){
			callback(output)
		});
	}

	factory.queryStampPaper = function(id, callback){
    	$http.get('/get_sp/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordStamp = function(data, callback){

		var stamp = data.stampid + "-" + data.uid + "-" + data.stamp_holder + "-" + data.location + "-" + data.doc_type+ "-" + data.doc_content;

    	$http.get('/add_sp/'+stamp).success(function(output){
			callback(output)
		});
	}

	factory.changeHolder = function(data, callback){

		var holder = data.id + "-" + data.name;

    	$http.get('/change_holder/'+holder).success(function(output){
			callback(output)
		});
	}

	return factory;
});


