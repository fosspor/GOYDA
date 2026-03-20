import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navigation from './components/Navigation';
import Home from './pages/Home';
import Discover from './pages/Discover';
import SuggestedRoutes from './pages/SuggestedRoutes';
import RouteDetail from './pages/RouteDetail';
import Profile from './pages/Profile';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Navigation />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/discover" element={<Discover />} />
        <Route path="/suggested-routes" element={<SuggestedRoutes />} />
        <Route path="/route/:id" element={<RouteDetail />} />
        <Route path="/profile" element={<Profile />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
