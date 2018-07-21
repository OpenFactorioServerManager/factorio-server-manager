import React from 'react';
import ReactDOM from 'react-dom';
import {Route} from 'react-router';
import App from './App/App.jsx';
import LoginContent from './App/components/LoginContent.jsx';
import {BrowserRouter, Switch} from "react-router-dom";

ReactDOM.render((
    <BrowserRouter>
        <Switch>
            <Route path="/login" component={LoginContent} />
            <Route component={App} />
        </Switch>
    </BrowserRouter>
), document.getElementById('app'))

