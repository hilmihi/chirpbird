import React, { useReducer, createContext } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginComp from './user/login';
import { Chat } from './chat/Chat';

function App() {
  return (
      <BrowserRouter>
        <Routes>
            <Route exact path="/" element={<LoginComp/>}/>
            <Route exact path="/chat" element={<Chat/>}/>
        </Routes>
      </BrowserRouter>
  );
}

export default App;