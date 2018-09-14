//SPDX-License-Identifier: Apache-2.0

var stamp = require('./controller.js');

module.exports = function(app){
  app.get('/get_sp/:id', function(req, res){
    stamp.get_sp(req, res);
  });
  app.get('/add_sp/:stamp', function(req, res){
    stamp.add_sp(req, res);
  });
  app.get('/get_all_sp', function(req, res){
    stamp.get_all_sp(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    stamp.change_holder(req, res);
  });
}
