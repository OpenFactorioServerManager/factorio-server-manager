import regeneratorRuntime from "regenerator-runtime"
import Bus from "./notifications"
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App/App.jsx';

window.flash = (message, color="gray-light") => Bus.emit('flash', ({message, color}));

ReactDOM.render(
    <App/>
, document.getElementById('app'));
