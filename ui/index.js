import regeneratorRuntime from "regenerator-runtime"

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
