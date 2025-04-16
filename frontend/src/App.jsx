import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import LoginSuccessPage from './pages/LoginSuccessPage';
import HomePage from './pages/HomePage';
import SubscriptionsPage from './pages/SubscriptionsPage';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/login-success" element={<LoginSuccessPage />} />
        
        {/* Protected routes */}
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/subscriptions" element={<SubscriptionsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
