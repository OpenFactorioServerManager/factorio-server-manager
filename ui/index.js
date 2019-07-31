global.$ = global.jQuery = require('jquery');
require('bootstrap');
require('admin-lte/build/js/AdminLTE.js');
require('bootstrap-fileinput');
require('bootstrap-fileinput/themes/fas/theme');


/**
 * Change the plus and minus when opening or closing bootstraps collapse object
 */
// $('body').on("show.bs.collapse hide.bs.collapse", (e) => {
//     let $target = $(e.target);
//     let $box = $target.parent(".box");
//     let $fontAwesome = $box.children(".box-header").children("i");
//
//     if(e.type == "show") {
//         $fontAwesome.removeClass("fa-plus").addClass("fa-minus");
//     } else if(e.type == "hide") {
//         $fontAwesome.removeClass("fa-minus").addClass("fa-plus");
//     }
// });


/**
 * Import Stuff for React
 */
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App/App.jsx';
import LoginContent from './App/components/LoginContent.jsx';
import {BrowserRouter, Switch, Route} from "react-router-dom";

/**
 * Start React Render
 */
ReactDOM.render((
    <BrowserRouter>
        <Switch>
            <Route path="/login" component={LoginContent} />
            <Route component={App} />
        </Switch>
    </BrowserRouter>
), document.getElementById('app'));
