import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, browserHistory, IndexRoute} from 'react-router';
import App from './App/App.jsx';
import ModsContent from './App/components/ModsContent.jsx';
import LogsContent from './App/components/LogsContent.jsx';
import SavesContent from './App/components/SavesContent.jsx';
import ConfigContent from './App/components/ConfigContent.jsx';
import LoginContent from './App/components/LoginContent.jsx';
import UsersContent from './App/components/UsersContent.jsx';
import Index from './App/components/Index.jsx';


ReactDOM.render((
    <Router history={browserHistory}>
        <Route path="/login" component={LoginContent}/>
        <Route path="/" component={App}>
            <IndexRoute component={Index}/>
            <Route path="/server" component={Index}/> 
            <Route path="/settings" component={UsersContent}/>
            <Route path="/mods" component={ModsContent}/> 
            <Route path="/logs" component={LogsContent}/> 
            <Route path="/saves" component={SavesContent}/> 
            <Route path="/config" component={ConfigContent}/> 
        </Route>
    </Router>
), document.getElementById('app'))

