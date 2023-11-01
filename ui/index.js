import regeneratorRuntime from "regenerator-runtime"
import Bus from "./notifications"
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App/App.jsx';

window.flash = (message, color="gray-light") => Bus.emit('flash', ({message, color}));

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<App/>);
